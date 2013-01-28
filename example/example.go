package main

import (
	"fmt"
	"github.com/argusdusty/ferret"
)

var ExampleWords = [][]byte{
	[]byte("abdeblah"),
	[]byte("foobar"),
	[]byte("barfoo"),
	[]byte("qwerty"),
	[]byte("testing"),
	[]byte("example"),
	[]byte("dictionary"),
	[]byte("dvorak"),
	[]byte("ferret"),
}

var ExampleCorrection = func(b []byte) [][]byte { return ferret.ErrorCorrect(b, ferret.LowercaseASCII) }
var ExampleSorter = func(s []byte) float64 { return float64(-len(s)) }

func main() {
	ExampleInvertedSuffix := ferret.MakeInvertedSuffix(ExampleWords)
	PrintArray(ExampleInvertedSuffix.Query([]byte("ar"), 5))
	PrintArray(ExampleInvertedSuffix.Query([]byte("test"), 5))
	PrintArray(ExampleInvertedSuffix.ErrorCorrectingQuery([]byte("tsst"), 5, ExampleCorrection))
	PrintArray(ExampleInvertedSuffix.SortedErrorCorrectingQuery([]byte("tsst"), 5, ExampleCorrection, ExampleSorter))
	PrintArray(ExampleInvertedSuffix.SortedQuery([]byte("a"), 5, ExampleSorter))
	PrintArray(ExampleInvertedSuffix.Query([]byte("a"), 5))
	ExampleInvertedSuffix.Insert([]byte("asdfghjklqwertyuiopzxcvbnm"))
	PrintArray(ExampleInvertedSuffix.Query([]byte("sdfghjklqwert"), 5))
}

func PrintArray(b [][]byte) {
	Result := make([]string, 0)
	for _, Word := range b {
		Result = append(Result, string(Word))
	}
	fmt.Println(Result)
}
