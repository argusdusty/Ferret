package main

import (
	"github.com/argusdusty/ferret"
	"fmt"
)

var ExampleDictionary = []string{
	"abdeblah"
	"foobar"
	"barfoo"
	"qwerty"
	"testing"
	"example"
	"dictionary"
	"dvorak"
	"ferret"
}

func main() {
	ExampleConversion := func(s string) []byte { return []byte(s) }
	ExampleInvertedSuffix := ferret.MakeInvertedSuffix(ExampleDictionary, ExampleConversion, 5)
	fmt.Println(ferret.Query([]byte{"test"}))
	fmt.Println(ferret.Query([]byte{"a"}))
	fmt.Println(ferret.PrefixQuery([]byte{"a"}))
	fmt.Println(ferret.Query([]byte{"ar"}))
}