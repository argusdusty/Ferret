package ferret

import "sort"

type InvertedSuffix struct {
	WordIndex   []int // WordIndex and SuffixIndex are sorted by Words[WordIndex[i]][SuffixIndex[i]:]
	SuffixIndex []int
	Words       [][]byte
	Results     []string
	Values      []interface{} // Some data mapped to the words
	Converter   func(string) []byte
}

func (IS *InvertedSuffix) Swap(i, j int) {
	IS.WordIndex[i], IS.WordIndex[j] = IS.WordIndex[j], IS.WordIndex[i]
	IS.SuffixIndex[i], IS.SuffixIndex[j] = IS.SuffixIndex[j], IS.SuffixIndex[i]
}

func (IS *InvertedSuffix) Len() int {
	return len(IS.WordIndex)
}

// Equivalent to:
// bytes.Compare(S.Words[S.WordIndex[i]][S.SuffixIndex[i]:], S.Words[S.WordIndex[j]][S.SuffixIndex[j]:]) <= 0
// but faster
func (IS *InvertedSuffix) Less(i, j int) bool {
	x := IS.WordIndex[i]
	y := IS.WordIndex[j]
	a := IS.Words[x]
	b := IS.Words[y]
	pa := IS.SuffixIndex[i]
	pb := IS.SuffixIndex[j]
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
	if pa == na {
		return true
	}
	return false
}

// Creates an inverted suffix from a dictionary of byte arrays
func MakeInvertedSuffix(Words, Results []string, Data []interface{}, Converter func(string) []byte) *InvertedSuffix {
	WordIndex := make([]int, 0)
	SuffixIndex := make([]int, 0)
	NewWords := make([][]byte, 0)
	for i, Word := range Words {
		word := Converter(Word)
		for j := 0; j < len(word); j++ {
			WordIndex = append(WordIndex, i)
			SuffixIndex = append(SuffixIndex, j)
		}
		NewWords = append(NewWords, word)
	}
	Suffixes := &InvertedSuffix{WordIndex, SuffixIndex, NewWords, Results, Data, Converter}
	sort.Sort(Suffixes)
	return Suffixes
}

// Adds a word to the dictionary that IS was built on.
// This is pretty slow, so stick to MakeInvertedSuffix when you can
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
	high := IS.Len()
	n := len(Query)
	for a := 0; a < n; a++ {
		c := Query[a]
		oldlow := low
		oldhigh := high
		i := low
		j := high
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

// Returns the strings which contain the query, and their stored values
// Unsorted. Might be partially sorted alphabetically?
// Input:
//     Word: The substring to search for.
//     ResultsLimit: Limit the results to some number of values. Set to -1 for no limit
func (IS *InvertedSuffix) Query(Word string, ResultsLimit int) ([]string, []interface{}) {
	Query := IS.Converter(Word)
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

// Returns the strings which contain the query
// Sorted. The function sorter produces a value to sort by (largest first)
// Input:
//     Word: The substring to search for.
//     ResultsLimit: Limit the results to some number of values. Set to -1 for no limit
//     Sorter: Takes (Result, Value, Length, Index (where Query begins in Result)) (string, []byte, int, int)
//         and produces a value (float64) to sort by (largest first).
func (IS *InvertedSuffix) SortedQuery(Word string, ResultsLimit int, Sorter func(string, interface{}, int, int) float64) ([]string, []interface{}, []float64) {
	Query := IS.Converter(Word)
	Results := make([]string, 0, ResultsLimit)
	Values := make([]interface{}, 0, ResultsLimit)
	Scores := make([]float64, 0, ResultsLimit)
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
	Results := make([]string, 0, ResultsLimit)
	Values := make([]interface{}, 0, ResultsLimit)
	Scores := make([]float64, 0, ResultsLimit)
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
