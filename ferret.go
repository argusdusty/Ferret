package ferret

import "sort"

type InvertedSuffix struct {
	WordIndex   []int // WordIndex and SuffixIndex are sorted by Words[WordIndex[i]][SuffixIndex[i]:]
	SuffixIndex []int
	Words       [][]byte
	OrigWords   []string
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
func MakeInvertedSuffix(Words []string, Data []interface{}, Converter func(string) []byte) *InvertedSuffix {
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
	Suffixes := &InvertedSuffix{WordIndex, SuffixIndex, NewWords, Words, Data, Converter}
	sort.Sort(Suffixes)
	return Suffixes
}

// Adds a word to the dictionary that IS was built on.
// This is pretty slow, so stick to MakeInvertedSuffix when you can
func (IS *InvertedSuffix) Insert(Word string, Data interface{}) {
	Query := IS.Converter(Word)
	low, high := IS.Search(Query)
	for k := low; k < high; k++ {
		if IS.OrigWords[IS.WordIndex[k]] == Word {
			IS.Values[IS.WordIndex[k]] = Data
			return
		}
	}
	i := len(IS.Words)
	IS.Words = append(IS.Words, Query)
	Length := len(Query)
	IS.OrigWords = append(IS.OrigWords, Word)
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
	Results := make([]string, 0)
	Values := make([]interface{}, 0)
	low, high := IS.Search(Query)
	a := 0
	used := make(map[int]bool, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		if _, ok := used[x]; ok {
			continue
		}
		used[x] = true
		Results = append(Results, IS.OrigWords[x])
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
	Results := make([]string, 0)
	Values := make([]interface{}, 0)
	Scores := make([]float64, 0)
	low, high := IS.Search(Query)
	a := 0
	used := make(map[int]float64, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		w := IS.OrigWords[x]
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
	Results := make([]string, 0)
	Values := make([]interface{}, 0)
	low, high := IS.Search(Query)
	a := 0
	used := make(map[int]bool, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		if _, ok := used[x]; ok {
			continue
		}
		used[x] = true
		Results = append(Results, IS.OrigWords[x])
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
				Results = append(Results, IS.OrigWords[x])
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
	Results := make([]string, 0)
	Values := make([]interface{}, 0)
	Scores := make([]float64, 0)
	low, high := IS.Search(Query)
	a := 0
	used := make(map[int]float64, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		w := IS.OrigWords[x]
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
				w := IS.OrigWords[x]
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

// Disabled for now while I work out how to adapt the below algorithm to the new data structure
/*

// A variant of the InvertedSuffix, which splits the InvertedSuffixes by value
// This allows for faster sorted searches on most dictionaries
type SortedInvertedSuffix struct {
	Sorter    func([]byte) float64
	Divisions []float64
	Data      map[int]*InvertedSuffix
}

// Creates an sorted inverted suffix from a dictionary of byte arrays,
// a sorting function which takes a byte array, and returns a float64 to sort by,
// and a list of divisions
func MakeSortedInvertedSuffix(Words [][]byte, Sorter func([]byte) float64, Divisions []float64) *SortedInvertedSuffix {
	sort.Float64s(Divisions)
	Data := make(map[int]*InvertedSuffix, 0)
	Wordss := make(map[int][][]byte, 0)
	for _, Word := range Words {
		x := Sorter(Word)
		N := 0
		j := len(Divisions)
		for N < j {
			h := (N + j) >> 1
			if Divisions[h] < x {
				N = h + 1
			} else {
				j = h
			}
		}
		if N != 0 {
			N -= 1
		}
		if _, ok := Wordss[N]; !ok {
			Wordss[N] = make([][]byte, 0)
		}
		Wordss[N] = append(Wordss[N], Word)
	}
	for N, Words := range Wordss {
		WordIndex := make([]int, 0)
		SuffixIndex := make([]int, 0)
		ISLengths := make([]int, 0)
		for i, Word := range Words {
			Length := len(Word)
			for j := 0; j < Length; j++ {
				WordIndex = append(WordIndex, i)
				SuffixIndex = append(SuffixIndex, j)
			}
			ISLengths = append(ISLengths, Length)
		}
		Suffixes := &InvertedSuffix{WordIndex, SuffixIndex, Words, ISLengths}
		sort.Sort(Suffixes)
		Data[N] = Suffixes
	}
	return &SortedInvertedSuffix{Sorter, Divisions, Data}
}

// Just like InvertedSuffix.Insert, but not always as slow
func (SIS *SortedInvertedSuffix) Insert(Word []byte) {
	v := SIS.Sorter(Word)
	i := 0
	j := len(SIS.Divisions)
	for i < j {
		h := (i + j) >> 1
		if SIS.Divisions[h] < v {
			i = h + 1
		} else {
			j = h
		}
	}
	if i != 0 {
		i -= 1
	}
	IS, ok := SIS.Data[i]
	if !ok {
		Words := [][]byte{Word}
		WordIndex := make([]int, 0)
		SuffixIndex := make([]int, 0)
		Length := len(Word)
		for j := 0; j < Length; j++ {
			WordIndex = append(WordIndex, i)
			SuffixIndex = append(SuffixIndex, j)
		}
		ISLengths := []int{Length}
		ISLengths = append(ISLengths, i)
		Suffixes := &InvertedSuffix{WordIndex, SuffixIndex, Words, ISLengths}
		sort.Sort(Suffixes)
		SIS.Data[i] = Suffixes
	} else {
		IS.Insert(Word)
	}
}

// Same as InvertedSuffix.SortedQuery
func (SIS *SortedInvertedSuffix) Query(Query []byte, ResultsLimit int) [][]byte {
	Limit := len(SIS.Divisions)
	if Limit == 0 {
		return make([][]byte, 0)
	}
	Values := make([]float64, 0)
	Results := make([][]byte, 0)
	a := 0
	for i := Limit - 1; i >= 0; i-- {
		IS, ok := SIS.Data[i]
		if !ok {
			continue
		}
		low, high := IS.Search(Query)
		used := make(map[int]bool, 0)
		for k := low; k < high; k++ {
			x := IS.WordIndex[k]
			if _, ok := used[x]; ok {
				continue
			}
			used[x] = true
			w := IS.Words[x]
			v := SIS.Sorter(w)
			i := 0
			j := a
			for i < j {
				h := (i + j) >> 1
				if Values[h] > v {
					i = h + 1
				} else {
					j = h
				}
			}
			if a == ResultsLimit {
				if i < a {
					Values = append(Values[:i], append([]float64{v}, Values[i:a-1]...)...)
					Results = append(Results[:i], append([][]byte{w}, Results[i:a-1]...)...)
				}
			} else {
				a++
				if i < a {
					Values = append(Values[:i], append([]float64{v}, Values[i:]...)...)
					Results = append(Results[:i], append([][]byte{w}, Results[i:]...)...)
				} else {
					Values = append(Values, v)
					Results = append(Results, w)
				}
			}
		}
		if a == ResultsLimit {
			return Results
		}
	}
	return Results
}

// Same as InvertedSuffix.SortedErrorCorrectingQuery
func (SIS *SortedInvertedSuffix) ErrorCorrectingQuery(Query []byte, ResultsLimit int, ErrorCorrection func([]byte) [][]byte) [][]byte {
	Limit := len(SIS.Divisions)
	if Limit == 0 {
		return make([][]byte, 0)
	}
	Values := make([]float64, 0)
	Results := make([][]byte, 0)
	a := 0
	for i := Limit - 1; i >= 0; i-- {
		IS, ok := SIS.Data[i]
		if !ok {
			continue
		}
		low, high := IS.Search(Query)
		used := make(map[int]bool, 0)
		for k := low; k < high; k++ {
			x := IS.WordIndex[k]
			if _, ok := used[x]; ok {
				continue
			}
			used[x] = true
			w := IS.Words[x]
			v := SIS.Sorter(w)
			i := 0
			j := a
			for i < j {
				h := (i + j) >> 1
				if Values[h] > v {
					i = h + 1
				} else {
					j = h
				}
			}
			if a == ResultsLimit {
				if i < a {
					Values = append(Values[:i], append([]float64{v}, Values[i:a-1]...)...)
					Results = append(Results[:i], append([][]byte{w}, Results[i:a-1]...)...)
				}
			} else {
				a++
				if i < a {
					Values = append(Values[:i], append([]float64{v}, Values[i:]...)...)
					Results = append(Results[:i], append([][]byte{w}, Results[i:]...)...)
				} else {
					Values = append(Values, v)
					Results = append(Results, w)
				}
			}
		}
		if a == ResultsLimit {
			return Results
		}
	}
	if a == 0 {
		Errors := make([][]byte, 0)
		for _, q := range ErrorCorrection(Query) {
			Errors = append(Errors, q)
		}
		for i := Limit - 1; i >= 0; i-- {
			IS, ok := SIS.Data[i]
			if !ok {
				continue
			}
			for _, q := range Errors {
				low, high := IS.Search(q)
				used := make(map[int]bool, 0)
				for k := low; k < high; k++ {
					x := IS.WordIndex[k]
					if _, ok := used[x]; ok {
						continue
					}
					used[x] = true
					w := IS.Words[x]
					v := SIS.Sorter(w)
					i := 0
					j := a
					for i < j {
						h := (i + j) >> 1
						if Values[h] > v {
							i = h + 1
						} else {
							j = h
						}
					}
					if a == ResultsLimit {
						if i < a {
							Values = append(Values[:i], append([]float64{v}, Values[i:a-1]...)...)
							Results = append(Results[:i], append([][]byte{w}, Results[i:a-1]...)...)
						}
					} else {
						a++
						if i < a {
							Values = append(Values[:i], append([]float64{v}, Values[i:]...)...)
							Results = append(Results[:i], append([][]byte{w}, Results[i:]...)...)
						} else {
							Values = append(Values, v)
							Results = append(Results, w)
						}
					}
				}
			}
			if a == ResultsLimit {
				return Results
			}
		}
	}
	return Results
}
*/
