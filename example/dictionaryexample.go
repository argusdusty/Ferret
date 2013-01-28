package main

import (
	"bytes"
	"fmt"
	"github.com/argusdusty/ferret"
	"io/ioutil"
	"time"
)

var Correction = func(b []byte) [][]byte { return ferret.ErrorCorrect(b, ferret.LowercaseASCII) }
var LengthSorter = func(s []byte) float64 { return float64(-len(s)) }

func main() {
	t := time.Now()
	Data, err := ioutil.ReadFile("dictionary.dat")
	if err != nil {
		panic(err)
	}
	Words := make([][]byte, 0)
	for _, Word := range bytes.Split(Data, []byte("\r\n")) {
		Words = append(Words, Word)
	}
	fmt.Println("Loaded dictionary in:", time.Now().Sub(t))
	t = time.Now()

	InvertedSuffix := ferret.MakeInvertedSuffix(Words)
	fmt.Println("Created index in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(InvertedSuffix.Query([]byte("ar"), 5))
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(InvertedSuffix.Query([]byte("test"), 5))
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(InvertedSuffix.ErrorCorrectingQuery([]byte("tsst"), 5, Correction))
	fmt.Println("Performed error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(InvertedSuffix.SortedErrorCorrectingQuery([]byte("tsst"), 5, Correction, LengthSorter))
	fmt.Println("Performed sorted error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(InvertedSuffix.SortedQuery([]byte("a"), 5, LengthSorter))
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(InvertedSuffix.Query([]byte("a"), 5))
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(InvertedSuffix.Query([]byte("the"), 25))
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()
	InvertedSuffix.Insert([]byte("asdfghjklqwertyuiopzxcvbnm"))
	fmt.Println("Performed insert in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(InvertedSuffix.Query([]byte("sdfghjklqwert"), 5))
	fmt.Println("Performed search in:", time.Now().Sub(t))
	t = time.Now()

	SIS := ferret.MakeSortedInvertedSuffix(Words, LengthSorter, []float64{-20.0, -15.0, -10.0, -5.0})
	fmt.Println("Created SIS in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(SIS.Query([]byte("ar"), 5))
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(SIS.Query([]byte("test"), 5))
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(SIS.ErrorCorrectingQuery([]byte("tsst"), 5, Correction))
	fmt.Println("Performed sorted error correcting search in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(SIS.Query([]byte("a"), 5))
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(SIS.Query([]byte("the"), 25))
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
	SIS.Insert([]byte("asdfghjklqwertyuiopzxcvbnm"))
	fmt.Println("Performed insert in:", time.Now().Sub(t))
	t = time.Now()
	PrintArray(SIS.Query([]byte("sdfghjklqwert"), 5))
	fmt.Println("Performed sorted search in:", time.Now().Sub(t))
	t = time.Now()
}

func PrintArray(b [][]byte) {
	Result := make([]string, 0)
	for _, Word := range b {
		Result = append(Result, string(Word))
	}
	fmt.Println(Result)
}
