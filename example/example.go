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

func main() {
	ExampleConversion := func(s string) []byte { return []byte(s) }
	ExampleInvertedSuffix := ferret.MakeInvertedSuffix(ExampleDictionary, ExampleConversion)
	fmt.Println(ExampleInvertedSuffix.Query([]byte("ar"), 5, false, false, make([]byte, 0)))
	fmt.Println(ExampleInvertedSuffix.Query([]byte("test"), 5, false, false, make([]byte, 0)))
	fmt.Println(ExampleInvertedSuffix.Query([]byte("tsst"), 5, false, true, ferret.LowercaseASCII))
	fmt.Println(ExampleInvertedSuffix.Query([]byte("a"), 5, true, false, make([]byte, 0)))
	fmt.Println(ExampleInvertedSuffix.Query([]byte("a"), 5, false, false, make([]byte, 0)))
}