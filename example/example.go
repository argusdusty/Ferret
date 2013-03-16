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
	fmt.Println(ExampleInvertedSuffix.Query("ar", 5))
	fmt.Println(ExampleInvertedSuffix.Query("test", 5))
	fmt.Println(ExampleInvertedSuffix.ErrorCorrectingQuery("tsst", 5, ExampleCorrection))
	fmt.Println(ExampleInvertedSuffix.SortedErrorCorrectingQuery("tsst", 5, ExampleCorrection, ExampleSorter))
	fmt.Println(ExampleInvertedSuffix.SortedQuery("a", 5, ExampleSorter))
	fmt.Println(ExampleInvertedSuffix.Query("a", 5))
	ExampleInvertedSuffix.Insert("asdfghjklqwertyuiopzxcvbnm", []uint64{26})
	fmt.Println(ExampleInvertedSuffix.Query("sdfghjklqwert", 5))
	fmt.Println(ExampleInvertedSuffix.Query("ferret", 5))
	ExampleInvertedSuffix.Insert("ferret", []uint64{7})
	fmt.Println(ExampleInvertedSuffix.Query("ferret", 5))
}
