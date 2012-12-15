package main

import (
	"github.com/argusdusty/ferret"
	"fmt"
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

func main() {
	ExampleConversion := func(s string) []byte { return []byte(s) }
	ExampleInvertedSuffix := ferret.MakeInvertedSuffix(ExampleDictionary, ExampleConversion)
	fmt.Println(ExampleInvertedSuffix.Query("ar", 5, false))
	fmt.Println(ExampleInvertedSuffix.Query("test", 5, false))
	fmt.Println(ExampleInvertedSuffix.ErrorCorrectingQuery("tsst", 5, false, ExampleCorrection))
	fmt.Println(ExampleInvertedSuffix.Query("a", 5, true))
	fmt.Println(ExampleInvertedSuffix.Query("a", 5, false))
	ExampleInvertedSuffix.Insert("asdfghjklqwertyuiopzxcvbnm")
	fmt.Println(ExampleInvertedSuffix.Query("sdfghjklqwert", 5, false))
}