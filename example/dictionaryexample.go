package main

import (
	"bytes"
	"fmt"
	"github.com/argusdusty/ferret"
	"io/ioutil"
	"strconv"
	"time"
)

var Correction = func(b []byte) [][]byte { return ferret.ErrorCorrect(b, ferret.LowercaseLetters) }
var LengthSorter = func(s string, v interface{}, l int, i int) float64 { return -float64(l + i) }
var FreqSorter = func(s string, v interface{}, l int, i int) float64 { return float64(v.(uint64)) }
var Converter = ferret.UnicodeToLowerASCII

func main() {
	t := time.Now()
	Data, err := ioutil.ReadFile("dictionary.dat")
	if err != nil {
		panic(err)
	}
	Words := make([]string, 0)
	Values := make([]interface{}, 0)
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
		Values = append(Values, Freq)
	}
	fmt.Println("Loaded dictionary in:", time.Now().Sub(t))
	t = time.Now()

	InvertedSuffix := ferret.MakeInvertedSuffix(Words, Words, Values, Converter)
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
	InvertedSuffix.Insert("asdfghjklqwertyuiopzxcvbnm", "asdfghjklqwertyuiopzxcvbnm", uint64(0))
	fmt.Println("Performed insert in:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println(InvertedSuffix.Query("sdfghjklqwert", 5))
	fmt.Println("Performed search in:", time.Now().Sub(t))
	fmt.Println("Running benchmarks...")
	t = time.Now()
	n := 0
	for _, Query := range InvertedSuffix.Words {
		InvertedSuffix.Query(string(Query), 5)
		n++
	}
	fmt.Println("Performed", n, "limit-5 searches in:", time.Now().Sub(t))
	t = time.Now()
	n = 0
	for _, Query := range InvertedSuffix.Words {
		InvertedSuffix.Query(string(Query), 25)
		n++
	}
	fmt.Println("Performed", n, "limit-25 searches in:", time.Now().Sub(t))
	t = time.Now()
	n = 0
	for _, Query := range InvertedSuffix.Words {
		InvertedSuffix.SortedQuery(string(Query), 5, LengthSorter)
		n++
	}
	fmt.Println("Performed", n, "limit-5 length sorted searches in:", time.Now().Sub(t))
	t = time.Now()
	n = 0
	for _, Query := range InvertedSuffix.Words {
		InvertedSuffix.SortedQuery(string(Query), 25, LengthSorter)
		n++
	}
	fmt.Println("Performed", n, "limit-25 length sorted searches in:", time.Now().Sub(t))
	t = time.Now()
	n = 0
	for _, Query := range InvertedSuffix.Words {
		InvertedSuffix.SortedQuery(string(Query), 5, FreqSorter)
		n++
	}
	fmt.Println("Performed", n, "limit-5 frequency sorted searches in:", time.Now().Sub(t))
	t = time.Now()
	n = 0
	for _, Query := range InvertedSuffix.Words {
		InvertedSuffix.SortedQuery(string(Query), 25, FreqSorter)
		n++
	}
	fmt.Println("Performed", n, "limit-25 frequency sorted searches in:", time.Now().Sub(t))
	t = time.Now()
	n = 0
	for _, Query := range InvertedSuffix.Words[:2048] {
		InvertedSuffix.ErrorCorrectingQuery(string(Query)+"0", 5, Correction)
		n++
	}
	fmt.Println("Performed", n, "limit-5 error correcting searches in:", time.Now().Sub(t))
	t = time.Now()
	n = 0
	for _, Query := range InvertedSuffix.Words[:2048] {
		InvertedSuffix.ErrorCorrectingQuery(string(Query)+"0", 25, Correction)
		n++
	}
	fmt.Println("Performed", n, "limit-25 error correcting searches in:", time.Now().Sub(t))
	t = time.Now()
	n = 0
	for _, Query := range InvertedSuffix.Words[:2048] {
		InvertedSuffix.SortedErrorCorrectingQuery(string(Query)+"0", 5, Correction, LengthSorter)
		n++
	}
	fmt.Println("Performed", n, "limit-5 length sorted error correcting searches in:", time.Now().Sub(t))
	t = time.Now()
	n = 0
	for _, Query := range InvertedSuffix.Words[:2048] {
		InvertedSuffix.SortedErrorCorrectingQuery(string(Query)+"0", 25, Correction, LengthSorter)
		n++
	}
	fmt.Println("Performed", n, "limit-25 length sorted error correcting searches in:", time.Now().Sub(t))
	t = time.Now()
	n = 0
	for _, Query := range InvertedSuffix.Words[:2048] {
		InvertedSuffix.SortedErrorCorrectingQuery(string(Query)+"0", 5, Correction, FreqSorter)
		n++
	}
	fmt.Println("Performed", n, "limit-5 frequency sorted error correcting searches in:", time.Now().Sub(t))
	t = time.Now()
	n = 0
	for _, Query := range InvertedSuffix.Words[:2048] {
		InvertedSuffix.SortedErrorCorrectingQuery(string(Query)+"0", 25, Correction, FreqSorter)
		n++
	}
	fmt.Println("Performed", n, "limit-25 frequency sorted error correcting searches in:", time.Now().Sub(t))
	/*
		t = time.Now()
		n = 0
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
