package ferret

import (
	"sort"
)

type InvertedSuffix struct {
	WordIndex   []int // WordIndex and SuffixIndex are sorted by Words[LengthIndex[i]][WordIndex[i]][SuffixIndex[i]:]
	SuffixIndex []int
	Words       [][]byte // Words[i] = Conversion(Dictionary[i])
	Lengths     []int    // Lengths[i] == len(Words[i])
	Dictionary  []string
	Conversion  func(string) []byte
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
	na := IS.Lengths[x]
	nb := IS.Lengths[y]
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

// Creates an inverted suffix from a dictionary of strings
func MakeInvertedSuffix(Dictionary []string, Conversion func(string) []byte) *InvertedSuffix {
	WordIndex := make([]int, 0)
	SuffixIndex := make([]int, 0)
	Words := make([][]byte, 0)
	Lengths := make([]int, 0)
	for i, Word := range Dictionary {
		ByteWord := Conversion(Word)
		N := len(ByteWord)
		for j := 0; j < N; j++ {
			WordIndex = append(WordIndex, i)
			SuffixIndex = append(SuffixIndex, j)
		}
		Words = append(Words, ByteWord)
		Lengths = append(Lengths, N)
	}
	Suffixes := &InvertedSuffix{WordIndex, SuffixIndex, Words, Lengths, Dictionary, Conversion}
	sort.Sort(Suffixes)
	return Suffixes
}

// Adds a word to the dictionary that IS was built on.
// This is pretty slow, so stick to MakeInvertedSuffix when you can
func (IS *InvertedSuffix) Insert(Word string) {
	i := len(IS.Words)
	ByteWord := IS.Conversion(Word)
	IS.Words = append(IS.Words, ByteWord)
	Length := len(ByteWord)
	IS.Lengths = append(IS.Lengths, Length)
	IS.Dictionary = append(IS.Dictionary, Word)
	for j := 0; j < Length; j++ {
		k := sort.Search(IS.Len(), func(h int) bool {
			y := IS.WordIndex[h]
			a := ByteWord
			b := IS.Words[y]
			pa := j
			pb := IS.SuffixIndex[h]
			na := Length
			nb := IS.Lengths[y]
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
		})
		IS.WordIndex = append(IS.WordIndex[:k], append([]int{i}, IS.WordIndex[k:]...)...)
		IS.SuffixIndex = append(IS.SuffixIndex[:k], append([]int{j}, IS.SuffixIndex[k:]...)...)
	}
}

// Performs an exact substring search for the query in the word dictionary
// Returns the boundaries (low/high) of sorted suffixes which have the query as a prefix
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
			Length := IS.Lengths[Index]
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
			return low, high
		}
		j = high
		for i < j {
			h := (i + j) >> 1
			Index := IS.WordIndex[h]
			Word := IS.Words[Index]
			Length := IS.Lengths[Index]
			d := IS.SuffixIndex[h] + a
			if d >= Length {
				i = h + 1
			} else {
				e := Word[d]
				if e <= c {
					i = h + 1
				} else {
					j = h
				}
			}
		}
		high = j
		if low == high {
			return low, high
		}
	}
	return low, high
}

// Returns the strings which contain the query
// Unsorted, I think it's partially sorted alphabetically
// Input:
//     Query: The substring to search for. Will be converted to []byte with SE.Conversion
//     ResultsLimit: Limit the results to some number of values. Set to -1 for no limit
func (IS *InvertedSuffix) Query(Query string, ResultsLimit int) []string {
	Data := IS.Conversion(Query)
	Results := make([]string, 0)
	low, high := IS.Search(Data)
	a := 0
	used := make(map[int]bool, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		if _, ok := used[x]; ok {
			continue
		}
		used[x] = true
		Results = append(Results, IS.Dictionary[x])
		a++
		if a == ResultsLimit {
			return Results
		}
	}
	return Results
}

// Returns the strings which contain the query
// Sorted. The function sorter produces a value to sort by (largest first)
// Input:
//     Query: The substring to search for. Will be converted to []byte with SE.Conversion
//     ResultsLimit: Limit the results to some number of values. Set to -1 for no limit
//     Sorter: Takes a Word (string) an produces a value (float64) to sort by (largest first).
func (IS *InvertedSuffix) SortedQuery(Query string, ResultsLimit int, Sorter func(string) float64) []string {
	Data := IS.Conversion(Query)
	Values := make([]float64, 0)
	Results := make([]string, 0)
	low, high := IS.Search(Data)
	a := 0
	used := make(map[int]bool, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		if _, ok := used[x]; ok {
			continue
		}
		used[x] = true
		w := IS.Dictionary[x]
		v := Sorter(w)
		if a == ResultsLimit {
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
			if i < a {
				Values = append(Values[:i], append([]float64{v}, Values[i:a-1]...)...)
				Results = append(Results[:i], append([]string{w}, Results[i:a-1]...)...)
			}
		} else {
			a++
			Values = append(Values, v)
			Results = append(Results, w)
		}
	}
	return Results
}

// Returns the strings which contain the query
// Unsorted, I think it's partially sorted alphabetically
// Will search for all substrings defined by ErrorCorrection
// if no results are found on the initial query
// Input:
//     Query: The substring to search for. Will be converted to []byte with SE.Conversion
//     ResultsLimit: Limit the results so you don't return your whol dictionary by accident. Set to -1 for no limit
//     Unsorted: Whether or not the results should be sorted by length.
//     ErrorCorrection: Returns a list of alternate queries
func (IS *InvertedSuffix) ErrorCorrectingQuery(Query string, ResultsLimit int, ErrorCorrection func([]byte) [][]byte) []string {
	Data := IS.Conversion(Query)
	Results := make([]string, 0)
	low, high := IS.Search(Data)
	a := 0
	used := make(map[int]bool, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		if _, ok := used[x]; ok {
			continue
		}
		used[x] = true
		Results = append(Results, IS.Dictionary[x])
		a++
		if a == ResultsLimit {
			return Results
		}
	}
	if a == 0 {
		for _, q := range ErrorCorrection(Data) {
			low, high := IS.Search(q)
			for k := low; k < high; k++ {
				x := IS.WordIndex[k]
				if _, ok := used[x]; ok {
					continue
				}
				used[x] = true
				Results = append(Results, IS.Dictionary[x])
				a++
				if a == ResultsLimit {
					return Results
				}
			}
		}
	}
	return Results
}

// Returns the strings which contain the query
// Sorted. The function sorter produces a value to sort by (largest first)
// Will search for all substrings defined by ErrorCorrection
// if no results are found on the initial query
// Input:
//     Query: The substring to search for. Will be converted to []byte with SE.Conversion
//     ResultsLimit: Limit the results so you don't return your whol dictionary by accident. Set to -1 for no limit
//     ErrorCorrection: Returns a list of alternate queries
//     Sorter: Takes a Word (string) an produces a value (float64) to sort by (largest first).
func (IS *InvertedSuffix) SortedErrorCorrectingQuery(Query string, ResultsLimit int, ErrorCorrection func([]byte) [][]byte, Sorter func(string) float64) []string {
	Data := IS.Conversion(Query)
	Values := make([]float64, 0)
	Results := make([]string, 0)
	low, high := IS.Search(Data)
	a := 0
	used := make(map[int]bool, 0)
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		if _, ok := used[x]; ok {
			continue
		}
		used[x] = true
		w := IS.Dictionary[x]
		v := Sorter(w)
		if a == ResultsLimit {
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
			if i < a {
				Values = append(Values[:i], append([]float64{v}, Values[i:a-1]...)...)
				Results = append(Results[:i], append([]string{w}, Results[i:a-1]...)...)
			}
		} else {
			a++
			Values = append(Values, v)
			Results = append(Results, w)
		}
	}
	if a == 0 {
		for _, q := range ErrorCorrection(Data) {
			low, high := IS.Search(q)
			for k := low; k < high; k++ {
				x := IS.WordIndex[k]
				if _, ok := used[x]; ok {
					continue
				}
				used[x] = true
				w := IS.Dictionary[x]
				v := Sorter(w)
				if a == ResultsLimit {
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
					if i < a {
						Values = append(Values[:i], append([]float64{v}, Values[i:a-1]...)...)
						Results = append(Results[:i], append([]string{w}, Results[i:a-1]...)...)
					}
				} else {
					a++
					Values = append(Values, v)
					Results = append(Results, w)
				}
			}
		}
	}
	return Results
}

// Disabled while I experiment with data structures
/*
// A variant of the InvertedSuffix, which splits the InvertedSuffixes by value
// This allows for faster length-sorted searches on most dictionaries
type SortedInvertedSuffix struct {
	Values     []int // Sorted index for Data. Data[Scores[0]] represents the dictionary with the largest score
	Conversion func(string) []byte
	Sorter     func(string) float64
	Data       map[int]*InvertedSuffix
}

func MakeLengthSortedInvertedSuffix(Dictionary []string, Conversion func(string) []byte) *LengthSortedInvertedSuffix {
	Data := make(map[int]*InvertedSuffix, 0)
	Values := make([]float64, 0)
	Wordss := make(map[int][][]byte, 0)
	Dictionaries := make(map[int][]string, 0)
	for _, Word := range Dictionary {
		ByteWord := Conversion(Word)
		N := len(ByteWord)
		if _, ok := Wordss[N]; !ok {
			Lengths = append(Lengths, N)
			Wordss[N] = make([][]byte, 0)
			Dictionaries[N] = make([]string, 0)
		}
		Wordss[N] = append(Wordss[N], ByteWord)
		Dictionaries[N] = append(Dictionaries[N], Word)
	}
	sort.Ints(Lengths)
	for n, Words := range Wordss {
		WordIndex := make([]int, 0)
		SuffixIndex := make([]int, 0)
		ISLengths := make([]int, 0)
		for i, ByteWord := range Words {
			for j := 0; j < n; j++ {
				WordIndex = append(WordIndex, i)
				SuffixIndex = append(SuffixIndex, j)
			}
			Words = append(Words, ByteWord)
			ISLengths = append(ISLengths, n)
		}
		Suffixes := &InvertedSuffix{WordIndex, SuffixIndex, Words, ISLengths, Dictionaries[n], Conversion}
		sort.Sort(Suffixes)
		Data[n] = Suffixes
	}
	return &LengthSortedInvertedSuffix{Lengths, Conversion, Data}
}

func (LSIS *LengthSortedInvertedSuffix) Insert(Word string) {
	ByteWord := LSIS.Conversion(Word)
	N := len(ByteWord)
	if _, ok := LSIS.Data[N]; !ok {
		Words := [][]byte{ByteWord}
		Dictionary := []string{Word}
		i := 0
		j := len(LSIS.Lengths)
		for i < j {
			h := i + (j-i)/2
			if LSIS.Lengths[h] < N {
				i = h + 1
			} else {
				j = h
			}
		}
		LSIS.Lengths = append(LSIS.Lengths[:i], append([]int{N}, LSIS.Lengths[i:]...)...)
		WordIndex := make([]int, 0)
		SuffixIndex := make([]int, 0)
		for j := 0; j < N; j++ {
			WordIndex = append(WordIndex, 0)
			SuffixIndex = append(SuffixIndex, j)
		}
		Suffixes := &InvertedSuffix{WordIndex, SuffixIndex, Words, []int{N}, Dictionary, LSIS.Conversion}
		sort.Sort(Suffixes)
		LSIS.Data[N] = Suffixes
	} else {
		LSIS.Data[N].Insert(Word)
	}
}

func (LSIS *LengthSortedInvertedSuffix) Query(Query string, ResultsLimit int) []string {
	if len(LSIS.Lengths) == 0 {
		return make([]string, 0)
	}
	Data := LSIS.Conversion(Query)
	Results := make([]string, 0)
	Limit := len(LSIS.Lengths)
	a := 0
	n := len(Data)
	i := 0
	j := len(LSIS.Lengths)
	for i < j {
		h := i + (j-i)/2
		if LSIS.Lengths[h] < n {
			i = h + 1
		} else {
			j = h
		}
	}
	for ; i < Limit; i++ {
		IS, ok := LSIS.Data[LSIS.Lengths[i]]
		if !ok {
			continue
		}
		low, high := IS.Search(Data)
		used := make(map[int]bool, 0)
		for k := low; k < high; k++ {
			x := IS.WordIndex[k]
			if _, ok := used[x]; ok {
				continue
			}
			used[x] = true
			Results = append(Results, IS.Dictionary[x])
			a++
			if a == ResultsLimit {
				return Results
			}
		}
	}
	return Results
}

func (LSIS *LengthSortedInvertedSuffix) ErrorCorrectingQuery(Query string, ResultsLimit int, ErrorCorrection func([]byte) [][]byte) []string {
	if len(LSIS.Lengths) == 0 {
		return make([]string, 0)
	}
	Data := LSIS.Conversion(Query)
	Results := make([]string, 0)
	Limit := len(LSIS.Lengths)
	a := 0
	n := len(Data)
	i := 0
	j := len(LSIS.Lengths)
	for i < j {
		h := i + (j-i)/2
		if LSIS.Lengths[h] < n {
			i = h + 1
		} else {
			j = h
		}
	}
	for ; i < Limit; i++ {
		IS, ok := LSIS.Data[LSIS.Lengths[i]]
		if !ok {
			continue
		}
		low, high := IS.Search(Data)
		used := make(map[int]bool, 0)
		for k := low; k < high; k++ {
			x := IS.WordIndex[k]
			if _, ok := used[x]; ok {
				continue
			}
			used[x] = true
			Results = append(Results, IS.Dictionary[x])
			a++
			if a == ResultsLimit {
				return Results
			}
		}
	}
	if a == 0 {
		n--
		i := 0
		j := len(LSIS.Lengths)
		for i < j {
			h := i + (j-i)/2
			if LSIS.Lengths[h] < n {
				i = h + 1
			} else {
				j = h
			}
		}
		for ; i < Limit; i++ {
			for _, q := range ErrorCorrection(Data) {
				if LSIS.Lengths[i] < len(q) {
					continue
				}
				IS, ok := LSIS.Data[LSIS.Lengths[i]]
				if !ok {
					continue
				}
				low, high := IS.Search(q)
				used := make(map[int]bool, 0)
				for k := low; k < high; k++ {
					x := IS.WordIndex[k]
					if _, ok := used[x]; ok {
						continue
					}
					used[x] = true
					Results = append(Results, IS.Dictionary[x])
					a++
					if a == ResultsLimit {
						return Results
					}
				}
			}
		}
	}
	return Results
}
*/
