package main

import (
	"bytes"
	"fmt"
	"github.com/argusdusty/ferret"
	"io/ioutil"
	"strconv"
	"time"
)

var Correction = func(b []byte) [][]byte { return ferret.ErrorCorrect(b, ferret.LowercaseASCII) }
var LengthSorter = func(s string, v []uint64, l int, i int) float64 { return -float64(l + i) }
var FreqSorter = func(s string, v []uint64, l int, i int) float64 { return float64(v[0]) }
var Converter = ferret.UnicodeToLowerASCII

func main() {
	t := time.Now()
	Data, err := ioutil.ReadFile("dictionary.dat")
	if err != nil {
		panic(err)
	}
	Words := make([]string, 0)
	Values := make([][]uint64, 0)
	for _, Vals := range bytes.Split(Data, []byte("\r\n")) {
		WordFreq := bytes.Split(Vals, []byte(" "))
		if len(WordFreq) != 2 {
			continue
		}
		Freq, err := strconv.ParseUint(string(WordFreq[1]), 10, 64)
		if err != nil {
			continue
		}
		Words = append(Words, string(WordFreq[0]))
		Values = append(Values, []uint64{Freq})
	}
	fmt.Println("Loaded dictionary in:", time.Now().Sub(t))
	t = time.Now()

	InvertedSuffix := ferret.MakeInvertedSuffix(Words, Values, Converter)
	fmt.Println("Created index in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.Query("ar", 5))
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.Query("test", 5))
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.ErrorCorrectingQuery("tsst", 5, Correction))
	fmt.Println("Performed error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.SortedErrorCorrectingQuery("tssst", 5, Correction, LengthSorter))
	fmt.Println("Performed sorted error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.SortedErrorCorrectingQuery("tssst", 5, Correction, FreqSorter))
	fmt.Println("Performed sorted error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.SortedQuery("a", 5, LengthSorter))
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.SortedQuery("a", 5, FreqSorter))
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.Query("a", 5))
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.Query("the", 25))
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.SortedQuery("the", 25, FreqSorter))
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	InvertedSuffix.Insert("asdfghjklqwertyuiopzxcvbnm", []uint64{0})
	fmt.Println("Performed insert in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.Query("sdfghjklqwert", 5))
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println("Running benchmarks...")
	for _, Query := range InvertedSuffix.OrigWords {
		InvertedSuffix.Query(Query, 5)
	}
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()
	for _, Query := range InvertedSuffix.OrigWords {
		InvertedSuffix.Query(Query, 25)
	}
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()
	for _, Query := range InvertedSuffix.OrigWords {
		InvertedSuffix.SortedQuery(Query, 5, LengthSorter)
	}
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	for _, Query := range InvertedSuffix.OrigWords {
		InvertedSuffix.SortedQuery(Query, 25, LengthSorter)
	}
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	for _, Query := range InvertedSuffix.OrigWords {
		InvertedSuffix.SortedQuery(Query, 5, FreqSorter)
	}
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	for _, Query := range InvertedSuffix.OrigWords {
		InvertedSuffix.SortedQuery(Query, 25, FreqSorter)
	}
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	for _, Query := range InvertedSuffix.OrigWords[:1000] {
		InvertedSuffix.ErrorCorrectingQuery(Query+"0", 5, Correction)
	}
	fmt.Println("Performed error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	for _, Query := range InvertedSuffix.OrigWords[:1000] {
		InvertedSuffix.ErrorCorrectingQuery(Query+"0", 25, Correction)
	}
	fmt.Println("Performed error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	for _, Query := range InvertedSuffix.OrigWords[:1000] {
		InvertedSuffix.SortedErrorCorrectingQuery(Query+"0", 5, Correction, LengthSorter)
	}
	fmt.Println("Performed sorted error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	for _, Query := range InvertedSuffix.OrigWords[:1000] {
		InvertedSuffix.SortedErrorCorrectingQuery(Query+"0", 25, Correction, LengthSorter)
	}
	fmt.Println("Performed sorted error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	for _, Query := range InvertedSuffix.OrigWords[:1000] {
		InvertedSuffix.SortedErrorCorrectingQuery(Query+"0", 5, Correction, FreqSorter)
	}
	fmt.Println("Performed sorted error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	for _, Query := range InvertedSuffix.OrigWords[:1000] {
		InvertedSuffix.SortedErrorCorrectingQuery(Query+"0", 25, Correction, FreqSorter)
	}
	fmt.Println("Performed sorted error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	/*
		SIS := ferret.MakeSortedInvertedSuffix(Words, LengthSorter, []float64{-20.0, -15.0, -10.0, -5.0})
		fmt.Println("Created SIS in:", time.Now().Sub(t))
		t = time.Now()
		fmt.Println(SIS.Query("ar", 5))
		fmt.Println("Performed sorted search in:", time.Now().Sub(t))
		t = time.Now()
		fmt.Println(SIS.Query("test", 5))
		fmt.Println("Performed sorted search in:", time.Now().Sub(t))
		t = time.Now()
		fmt.Println(SIS.ErrorCorrectingQuery("tsst", 5, Correction))
		fmt.Println("Performed sorted error correcting search in:", time.Now().Sub(t))
		t = time.Now()
		fmt.Println(SIS.Query("a", 5))
		fmt.Println("Performed sorted search in:", time.Now().Sub(t))
		t = time.Now()
		fmt.Println(SIS.Query("the", 25))
		fmt.Println("Performed sorted search in:", time.Now().Sub(t))
		t = time.Now()
		SIS.Insert("asdfghjklqwertyuiopzxcvbnm")
		fmt.Println("Performed insert in:", time.Now().Sub(t))
		t = time.Now()
		fmt.Println(SIS.Query("sdfghjklqwert", 5))
		fmt.Println("Performed sorted search in:", time.Now().Sub(t))
		t = time.Now()
	*/
}
