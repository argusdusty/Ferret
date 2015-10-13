// Copyright 2013 Mark Canning
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// Author: Mark Canning
// Developed at: Tamber, Inc. (http://www.tamber.com/).
//
// Tamber also has this really cool recommendation engine for music
// (also development by me) which prioritizes up-and-coming artists, so
// it doesn't succomb to the popularity biases that plague modern
// recommendation engines, and still produces excellent personalized
// recommendations! Make sure to check us out at http://www.tamber.com
// or https://itunes.apple.com/us/app/tamber-concerts/id658240483

package ferret

import "sort"

type InvertedSuffix struct {
	WordIndex   []int               // WordIndex and SuffixIndex are sorted by Words[WordIndex[i]][SuffixIndex[i]:]
	SuffixIndex []int               // WordIndex and SuffixIndex are sorted by Words[WordIndex[i]][SuffixIndex[i]:]
	Words       [][]byte            // The words to perform substring searches over
	Results     []string            // The 'true' value of the words. Used as a return value
	Values      []interface{}       // Some data mapped to the words. Used for sorting, and as a return value
	Converter   func(string) []byte // A converter for an inserted word/query to a byte array to search for/with
}

// A wrapper type used to sort the three arrays according to sort.sort
type sortWrapper struct {
	WordIndex   []int
	SuffixIndex []int
	Words       [][]byte
}

func (SW *sortWrapper) Swap(i, j int) {
	SW.WordIndex[i], SW.WordIndex[j] = SW.WordIndex[j], SW.WordIndex[i]
	SW.SuffixIndex[i], SW.SuffixIndex[j] = SW.SuffixIndex[j], SW.SuffixIndex[i]
}

func (SW *sortWrapper) Len() int {
	return len(SW.WordIndex)
}

// Equivalent to:
// bytes.Compare(S.Words[S.WordIndex[i]][S.SuffixIndex[i]:], S.Words[S.WordIndex[j]][S.SuffixIndex[j]:]) <= 0
// but faster
func (SW *sortWrapper) Less(i, j int) bool {
	x := SW.WordIndex[i]
	y := SW.WordIndex[j]
	a := SW.Words[x]
	b := SW.Words[y]
	pa := SW.SuffixIndex[i]
	pb := SW.SuffixIndex[j]
	na := len(a)
	nb := len(b)
	for pa < na && pb < nb {
		wa := a[pa]
		wb := b[pb]
		if wa < wb {
			return true
		}
		if wb < wa {
			return false
		}
		pa++
		pb++
	}
	return pa == na
}

// Creates an inverted suffix from a dictionary of byte arrays, mapping data, and a string->[]byte converter
func New(Words, Results []string, Data []interface{}, Converter func(string) []byte) *InvertedSuffix {
	CharCount := 0
	NewWords := make([][]byte, len(Words))
	for i, Word := range Words {
		NewWord := Converter(Word)
		NewWords[i] = NewWord
		CharCount += len(NewWord)
	}
	WordIndex := make([]int, 0, CharCount)
	SuffixIndex := make([]int, 0, CharCount)
	for i, NewWord := range NewWords {
		for j := 0; j < len(NewWord); j++ {
			WordIndex = append(WordIndex, i)
			SuffixIndex = append(SuffixIndex, j)
		}
	}
	sort.Sort(&sortWrapper{WordIndex, SuffixIndex, NewWords})
	Suffixes := &InvertedSuffix{WordIndex, SuffixIndex, NewWords, Results, Data, Converter}
	return Suffixes
}

// Adds a word to the dictionary that IS was built on.
// This is pretty slow, because of linear-time insertion into an array,
// so stick to New when you can
func (IS *InvertedSuffix) Insert(Word, Result string, Data interface{}) {
	Query := IS.Converter(Word)
	low, high := IS.Search(Query)
	for k := low; k < high; k++ {
		if IS.Results[IS.WordIndex[k]] == Word {
			IS.Values[IS.WordIndex[k]] = Data
			return
		}
	}
	i := len(IS.Words)
	IS.Words = append(IS.Words, Query)
	Length := len(Query)
	IS.Results = append(IS.Results, Result)
	IS.Values = append(IS.Values, Data)
	for j := 0; j < Length; j++ {
		k, _ := IS.Search(Query[j:])
		IS.WordIndex = append(IS.WordIndex, 0)
		copy(IS.WordIndex[k+1:], IS.WordIndex[k:])
		IS.WordIndex[k] = i
		IS.SuffixIndex = append(IS.SuffixIndex, 0)
		copy(IS.SuffixIndex[k+1:], IS.SuffixIndex[k:])
		IS.SuffixIndex[k] = j
	}
}

// Performs an exact substring search for the query in the word dictionary
// Returns the boundaries (low/high) of sorted suffixes which have the query as a prefix
// This is a low-level interface. I wouldn't recommend using this yourself
func (IS *InvertedSuffix) Search(Query []byte) (int, int) {
	low := 0
	high := len(IS.WordIndex)
	n := len(Query)
	for a := 0; a < n; a++ {
		c := Query[a]
		oldlow := low
		oldhigh := high
		i := low
		j := high
		// Raise the lower-bound
		for i < j {
			h := (i + j) >> 1
			Index := IS.WordIndex[h]
			Word := IS.Words[Index]
			Length := len(Word)
			d := IS.SuffixIndex[h] + a
			if d >= Length {
				i = h + 1
			} else {
				e := Word[d]
				if e < c {
					i = h + 1
				} else {
					j = h
					if e > c {
						high = h
					}
				}
			}
		}
		low = i
		if low == high {
			break
		}
		j = high
		if low == oldlow && high == oldhigh {
			// nothing to do here: Word[IS.SuffixIndex[(i + j) >> 1] + a] == c
			continue
		}
		// Lower the upper-bound
		for i < j {
			h := (i + j) >> 1
			Index := IS.WordIndex[h]
			Word := IS.Words[Index]
			Length := len(Word)
			d := IS.SuffixIndex[h] + a
			if d >= Length {
				i = h + 1
			} else {
				e := Word[d]
				if e <= c {
					i = h + 1
					if e < c {
						low = i
					}
				} else {
					j = h
				}
			}
		}
		high = j
		if low == high {
			break
		}
	}
	return low, high
}

// Returns the strings which contain the query, and their stored values unsorted
// Input:
//     Word: The substring to search for.
//     ResultsLimit: Limit the results to some number of values. Set to -1 for no limit
func (IS *InvertedSuffix) Query(Word string, ResultsLimit int) ([]string, []interface{}) {
	Query := IS.Converter(Word)
	if ResultsLimit == 0 {
		return []string{}, []interface{}{}
	}
	if ResultsLimit < 0 {
		ResultsLimit = 0
	}
	Results := make([]string, 0, ResultsLimit)
	Values := make([]interface{}, 0, ResultsLimit)
	low, high := IS.Search(Query)
	a := 0
	used := make(map[int]bool, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		if _, ok := used[x]; ok {
			continue
		}
		used[x] = true
		Results = append(Results, IS.Results[x])
		Values = append(Values, IS.Values[x])
		a++
		if a == ResultsLimit {
			return Results, Values
		}
	}
	return Results, Values
}

// Returns the strings which contain the query sorted
// The function sorter produces a value to sort by (largest first)
// Input:
//     Word: The substring to search for.
//     ResultsLimit: Limit the results to some number of values. Set to -1 for no limit
//     Sorter: Takes (Result, Value, Length, Index (where Query begins in Result)) (string, []byte, int, int)
//         and produces a value (float64) to sort by (largest first).
func (IS *InvertedSuffix) SortedQuery(Word string, ResultsLimit int, Sorter func(string, interface{}, int, int) float64) ([]string, []interface{}, []float64) {
	Query := IS.Converter(Word)
	if ResultsLimit == 0 {
		return []string{}, []interface{}{}, []float64{}
	}
	if ResultsLimit < 0 {
		ResultsLimit = 0
	}
	Results := make([]string, 0, ResultsLimit)
	Values := make([]interface{}, 0, ResultsLimit)
	Scores := make([]float64, 0, ResultsLimit)
	if ResultsLimit == 0 {
		ResultsLimit = -1
	}
	low, high := IS.Search(Query)
	a := 0
	used := make(map[int]float64, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		w := IS.Results[x]
		v := IS.Values[x]
		s := Sorter(w, v, len(IS.Words[x]), IS.SuffixIndex[k])
		if ps, ok := used[x]; ok && ps >= s {
			continue
		}
		used[x] = s
		i := 0
		j := a
		for i < j {
			h := (i + j) >> 1
			if Scores[h] > s {
				i = h + 1
			} else {
				j = h
			}
		}
		if a == ResultsLimit {
			if i < a {
				copy(Results[i+1:], Results[i:a-1])
				Results[i] = w
				copy(Values[i+1:], Values[i:a-1])
				Values[i] = v
				copy(Scores[i+1:], Scores[i:a-1])
				Scores[i] = s
			}
		} else {
			a++
			if i < a {
				Results = append(Results, "")
				copy(Results[i+1:], Results[i:])
				Results[i] = w
				Values = append(Values, nil)
				copy(Values[i+1:], Values[i:a-1])
				Values[i] = v
				Scores = append(Scores, 0.0)
				copy(Scores[i+1:], Scores[i:a-1])
				Scores[i] = s
			} else {
				Results = append(Results, w)
				Values = append(Values, v)
				Scores = append(Scores, s)
			}
		}
	}
	return Results, Values, Scores
}

// Returns the strings which contain the query
// Unsorted, I think it's partially sorted alphabetically
// Will search for all substrings defined by ErrorCorrection
// if no results are found on the initial query
// Input:
//     Query: The substring to search for.
//     ResultsLimit: Limit the results so you don't return your whole dictionary by accident. Set to -1 for no limit
//     ErrorCorrection: Returns a list of alternate queries
func (IS *InvertedSuffix) ErrorCorrectingQuery(Word string, ResultsLimit int, ErrorCorrection func([]byte) [][]byte) ([]string, []interface{}) {
	Query := IS.Converter(Word)
	if ResultsLimit == 0 {
		return []string{}, []interface{}{}
	}
	if ResultsLimit < 0 {
		ResultsLimit = 0
	}
	Results := make([]string, 0, ResultsLimit)
	Values := make([]interface{}, 0, ResultsLimit)
	low, high := IS.Search(Query)
	a := 0
	used := make(map[int]bool, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		if _, ok := used[x]; ok {
			continue
		}
		used[x] = true
		Results = append(Results, IS.Results[x])
		Values = append(Values, IS.Values[x])
		a++
		if a == ResultsLimit {
			return Results, Values
		}
	}
	if a != ResultsLimit {
		for _, q := range ErrorCorrection(Query) {
			low, high := IS.Search(q)
			for k := low; k < high; k++ {
				x := IS.WordIndex[k]
				if _, ok := used[x]; ok {
					continue
				}
				used[x] = true
				Results = append(Results, IS.Results[x])
				Values = append(Values, IS.Values[x])
				a++
				if a == ResultsLimit {
					return Results, Values
				}
			}
		}
	}
	return Results, Values
}

// Returns the strings which contain the query
// Sorted. The function sorter produces a value to sort by (largest first)
// Will search for all substrings defined by ErrorCorrection
// if no results are found on the initial query
// Input:
//     Query: The substring to search for.
//     ResultsLimit: Limit the results so you don't return your whole dictionary by accident. Set to -1 for no limit
//     ErrorCorrection: Returns a list of alternate queries
//     Sorter: Takes (Result, Value, Length, Index (where Query begins in Result))
//         (string, []byte, int, int), and produces a value (float64) to sort by (largest first).
func (IS *InvertedSuffix) SortedErrorCorrectingQuery(Word string, ResultsLimit int, ErrorCorrection func([]byte) [][]byte, Sorter func(string, interface{}, int, int) float64) ([]string, []interface{}, []float64) {
	Query := IS.Converter(Word)
	if ResultsLimit == 0 {
		return []string{}, []interface{}{}, []float64{}
	}
	if ResultsLimit < 0 {
		ResultsLimit = 0
	}
	Results := make([]string, 0, ResultsLimit)
	Values := make([]interface{}, 0, ResultsLimit)
	Scores := make([]float64, 0, ResultsLimit)
	if ResultsLimit == 0 {
		ResultsLimit = -1
	}
	low, high := IS.Search(Query)
	a := 0
	used := make(map[int]float64, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		w := IS.Results[x]
		v := IS.Values[x]
		s := Sorter(w, v, len(IS.Words[x]), IS.SuffixIndex[k])
		if ps, ok := used[x]; ok && ps >= s {
			continue
		}
		used[x] = s
		i := 0
		j := a
		for i < j {
			h := (i + j) >> 1
			if Scores[h] > s {
				i = h + 1
			} else {
				j = h
			}
		}
		if a == ResultsLimit {
			if i < a {
				copy(Results[i+1:], Results[i:a-1])
				Results[i] = w
				copy(Values[i+1:], Values[i:a-1])
				Values[i] = v
				copy(Scores[i+1:], Scores[i:a-1])
				Scores[i] = s
			}
		} else {
			a++
			if i < a {
				Results = append(Results, "")
				copy(Results[i+1:], Results[i:])
				Results[i] = w
				Values = append(Values, nil)
				copy(Values[i+1:], Values[i:a-1])
				Values[i] = v
				Scores = append(Scores, 0.0)
				copy(Scores[i+1:], Scores[i:a-1])
				Scores[i] = s
			} else {
				Results = append(Results, w)
				Values = append(Values, v)
				Scores = append(Scores, s)
			}
		}
	}
	if a == 0 {
		for _, q := range ErrorCorrection(Query) {
			low, high := IS.Search(q)
			for k := low; k < high; k++ {
				x := IS.WordIndex[k]
				w := IS.Results[x]
				v := IS.Values[x]
				s := Sorter(w, v, len(IS.Words[x]), IS.SuffixIndex[k])
				if ps, ok := used[x]; ok && ps >= s {
					continue
				}
				used[x] = s
				i := 0
				j := a
				for i < j {
					h := (i + j) >> 1
					if Scores[h] > s {
						i = h + 1
					} else {
						j = h
					}
				}
				if a == ResultsLimit {
					if i < a {
						copy(Results[i+1:], Results[i:a-1])
						Results[i] = w
						copy(Values[i+1:], Values[i:a-1])
						Values[i] = v
						copy(Scores[i+1:], Scores[i:a-1])
						Scores[i] = s
					}
				} else {
					a++
					if i < a {
						Results = append(Results, "")
						copy(Results[i+1:], Results[i:])
						Results[i] = w
						Values = append(Values, nil)
						copy(Values[i+1:], Values[i:a-1])
						Values[i] = v
						Scores = append(Scores, 0.0)
						copy(Scores[i+1:], Scores[i:a-1])
						Scores[i] = s
					} else {
						Results = append(Results, w)
						Values = append(Values, v)
						Scores = append(Scores, s)
					}
				}
			}
		}
	}
	return Results, Values, Scores
}
