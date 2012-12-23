package main

import (
	"fmt"
	"github.com/argusdusty/ferret"
)

var ExampleDictionary = []string{
	"abdeblah",
	"foobar",
	"barfoo",
	"qwerty",
	"testing",
	"example",
	"dictionary",
	"dvorak",
	"ferret",
}
var ExampleCorrection = func(b []byte) [][]byte { return ferret.ErrorCorrect(b, ferret.LowercaseASCII) }
var ExampleSorter = func(s string) float64 { return float64(-len(s)) }
var ExampleConversion = func(s string) []byte { return []byte(s) }

func main() {
	ExampleInvertedSuffix := ferret.MakeInvertedSuffix(ExampleDictionary, ExampleConversion)
	fmt.Println(ExampleInvertedSuffix.Query("ar", 5, false))
	fmt.Println(ExampleInvertedSuffix.Query("test", 5, false))
	fmt.Println(ExampleInvertedSuffix.ErrorCorrectingQuery("tsst", 5, false, ExampleCorrection))
	fmt.Println(ExampleInvertedSuffix.SortedErrorCorrectingQuery("tsst", 5, false, ExampleCorrection, ExampleSorter))
	fmt.Println(ExampleInvertedSuffix.SortedQuery("a", 5, ExampleSorter))
	fmt.Println(ExampleInvertedSuffix.Query("a", 5, false))
	ExampleInvertedSuffix.Insert("asdfghjklqwertyuiopzxcvbnm")
	fmt.Println(ExampleInvertedSuffix.Query("sdfghjklqwert", 5, false))
}
