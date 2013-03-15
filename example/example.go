package main

import (
	"fmt"
	"github.com/argusdusty/ferret"
)

var ExampleWords = []string{
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

var ExampleData = [][]uint64{
	[]uint64{8},
	[]uint64{6},
	[]uint64{6},
	[]uint64{6},
	[]uint64{7},
	[]uint64{7},
	[]uint64{10},
	[]uint64{6},
	[]uint64{6},
}

var ExampleCorrection = func(b []byte) [][]byte { return ferret.ErrorCorrect(b, ferret.LowercaseASCII) }
var ExampleSorter = func(s string, v []uint64, l int, i int) float64 { return -float64(l + i) }
var ExampleConverter = func(s string) []byte { return []byte(s) }

func main() {
	ExampleInvertedSuffix := ferret.MakeInvertedSuffix(ExampleWords, ExampleData, ExampleConverter)
	PrintArrays(ExampleInvertedSuffix.Query("ar", 5))
	PrintArrays(ExampleInvertedSuffix.Query("test", 5))
	PrintArrays(ExampleInvertedSuffix.ErrorCorrectingQuery("tsst", 5, ExampleCorrection))
	PrintSortedArrays(ExampleInvertedSuffix.SortedErrorCorrectingQuery("tsst", 5, ExampleCorrection, ExampleSorter))
	PrintSortedArrays(ExampleInvertedSuffix.SortedQuery("a", 5, ExampleSorter))
	PrintArrays(ExampleInvertedSuffix.Query("a", 5))
	ExampleInvertedSuffix.Insert("asdfghjklqwertyuiopzxcvbnm", []uint64{26})
	PrintArrays(ExampleInvertedSuffix.Query("sdfghjklqwert", 5))
	PrintArrays(ExampleInvertedSuffix.Query("ferret", 5))
	ExampleInvertedSuffix.Insert("ferret", []uint64{7})
	PrintArrays(ExampleInvertedSuffix.Query("ferret", 5))
}

func PrintArrays(r []string, v [][]uint64) {
	fmt.Println(r, v)
}

func PrintSortedArrays(r []string, v [][]uint64, s []float64) {
	fmt.Println(r, v, s)
}
