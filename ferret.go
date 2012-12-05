package ferret

import (
	"sort"
	"strings"
)

type InvertedSuffix struct {
	WordIndex []int // WordIndex and SuffixIndex are sorted by S.Words[S.WordIndex[i]][S.SuffixIndex[i]:]
	SuffixIndex []int
	Words [][]byte // Unsorted list of []byte-converted dictionary words
	Lengths []int // Caching for performance. IS.Lengths[i] == len(IS.Words[i])
	OrigWords []string // Original value of []byte-converted words
	SortedWords []string // Sorted list of original dictionary words
	ResultsLimit int // Set to -1 for no results limit
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
	a := IS.Words[IS.WordIndex[i]]; b := IS.Words[IS.WordIndex[j]]
	pa := IS.SuffixIndex[i]; pb := IS.SuffixIndex[j]
	na := IS.Lengths[i]; nb := IS.Lengths[j]
	for pa < na && pb < nb {
		wa := a[pa]; wb := b[pb]
		if wa < wb { return true }
		if wb < wa { return false }
		pa++; pb++
	}
	if pa == na { return true }
	return false
}

// Creates an inverted suffix from a dictionary of strings
func MakeInvertedSuffix(Words []string, Conversion func(string) []byte, ResultsLimit int) *InvertedSuffix {
	WordIndex := make([]int, 0)
	SuffixIndex := make([]int, 0)
	NewWords := make([][]byte, 0)
	Lengths := make([]int, 0)
	OrigWords := make([]string, 0)
	SortedWords := make([]string, 0)
	for i, Word := range(Words) {
		ByteWord := Conversion(Word); N := len(ByteWord)
		for j := 0; j < N; j++ {
			WordIndex = append(WordIndex, i)
			SuffixIndex = append(SuffixIndex, j)
		}
		NewWords = append(NewWords, ByteWord)
		Lengths = append(Lengths, N)
		OrigWords = append(OrigWords, Word)
		SortedWords = append(SortedWords, Word)
	}
	sort.Strings(SortedWords)
	Suffixes := &InvertedSuffix{WordIndex, SuffixIndex, NewWords, Lengths, OrigWords, Words, ResultsLimit}
	sort.Sort(Suffixes)
	return Suffixes
}

// Adds a word to the dictionary that IS was built on.
// This is pretty slow, so stick to MakeInvertedSuffix when you can
func (IS *InvertedSuffix) Insert(Word []byte) {
	Length := len(Word)
	i := len(IS.Words)
	IS.Words = append(IS.Words, Word)
	IS.Lengths = append(IS.Lengths, Length)
	for j := 0; j < Length; j++ {
		k := sort.Search(IS.Len(), func(h int) bool {
			a := Word; b := IS.Words[IS.WordIndex[h]]
			pa := j; pb := IS.SuffixIndex[j]
			na := Length; nb := IS.Lengths[j]
			for pa < na && pb < nb {
				wa := a[pa]; wb := b[pb]
				if wa < wb { return true }
				if wb < wa { return false }
				pa++; pb++
			}
			if pa == na { return true }
			return false
		})
		IS.WordIndex = append(IS.WordIndex[:k], append([]int{i}, IS.WordIndex[k:]...)...)
		IS.SuffixIndex = append(IS.SuffixIndex[:k], append([]int{j}, IS.SuffixIndex[k:]...)...)
	}
}

// Performs an exact substring search for the query in the word dictionary
// Returns the boundaries of sorted suffixes which match have the query as a prefix
func (IS *InvertedSuffix) Search(Query []byte) (int, int) {
	low := 0; high := IS.Len(); n := len(Query)
	for a := 0; a < n; a++ {
		c := Query[a]
		i := low; j := high
		for i < j {
			h := (i + j) >> 1
			Index := IS.WordIndex[h]
			Word := IS.Words[Index]
			Length := IS.Lengths[Index]
			d := IS.SuffixIndex[h]+a
			if d >= Length {
				i = h + 1
			} else {
				e := Word[d]
				if e < c {
					i = h + 1
				} else {
					j = h
					if e > c { high = h }
				}
			}
		}
		low = i
		if low == high { return low, high }
		j = high
		for i < j {
			h := (i + j) >> 1
			Index := IS.WordIndex[h]
			Word := IS.Words[Index]
			Length := IS.Lengths[Index]
			d := IS.SuffixIndex[h]+a
			if d >= Length {
				i = h + 1
			} else {
				e := Word[d]
				if e < c {
					i = h + 1
				} else {
					j = h
				}
			}
		}
		high = j
		if low == high { return low, high }
	}
	return low, high
}

// Returns the strings which contain the query
func (IS *InvertedSuffix) Query(Query []byte, FaultTolerance bool, AllowedBytes []byte) []string {
	Values := make([]int, 0); Results := make([]string, 0); n := len(Query)
	low, high := IS.Search(Query)
	a := 0; used := make(map[int]bool, 0); Limit := IS.ResultsLimit
	for k := low; k < high; k++ {
		x := IS.WordIndex[k]
		if _, ok := used[x]; ok { continue }
		used[x] = true
		z := IS.Lengths[x]
		w := IS.OrigWords[x]
		i := sort.Search(a, func(h int) bool { return Values[h] > z })
		Values = append(Values[:i], append([]int{z}, Values[i:]...)...)
		Results = append(Results[:i], append([]string{w}, Results[i:]...)...)
		if a == Limit {
			Values = Values[:a]; Results = Results[:a]
		} else {
			a++
		}
	}
	if FaultTolerance && a == 0 && n < 20 {
		for _, q := range(ErrorCorrect(Query, LowercaseASCII)) {
			low, high := IS.Search(q)
			for k := low; k < high; k++ {
				x := IS.WordIndex[k]
				if _, ok := used[x]; ok { continue }
				used[x] = true
				z := IS.Lengths[x]
				w := IS.OrigWords[x]
				i := sort.Search(a, func(h int) bool { return Values[h] > z })
				Values = append(Values[:i], append([]int{z}, Values[i:]...)...)
				Results = append(Results[:i], append([]string{w}, Results[i:]...)...)
				if a == Limit {
					Values = Values[:a]; Results = Results[:a]
				} else {
					a++
				}
			}
		}
	}
	return Results
}

// Returns the strings which have the query as a prefix
// Unoptomized, but should still run in the optimal running time (O(ln(IS.Len())*len(Query)))
func (IS *InvertedSuffix) PrefixQuery(Query string) []string {
	Results := make([]string, 0); a := 0
	for i := sort.SearchStrings(IS.SortedWords, Query); ; i++ {
		Word := IS.SortedWords[i]
		if !strings.HasPrefix(Word, Query) { break }
		Results = append(Results, Word); a++
		if a == IS.ResultsLimit { return Results }
	}
	return Results
}