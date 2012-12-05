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
	ExampleInvertedSuffix := ferret.MakeInvertedSuffix(ExampleDictionary, ExampleConversion, 5)
	fmt.Println(ExampleInvertedSuffix.Query([]byte("test"), false, make([]byte, 0)))
	fmt.Println(ExampleInvertedSuffix.Query([]byte("a"), false, make([]byte, 0)))
	fmt.Println(ExampleInvertedSuffix.Query([]byte("ar"), false, make([]byte, 0)))
	fmt.Println(ExampleInvertedSuffix.Query([]byte("tsst"), true, ferret.LowercaseASCII))
}