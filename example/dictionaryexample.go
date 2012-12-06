package main

import (
	"io/ioutil"
	"bytes"
	"fmt"
	"github.com/argusdusty/ferret"
)

var Conversion = func(s string) []byte { return []byte(s) }

func main() {
	Data, err := ioutil.ReadFile("dictionary.dat")
	if err != nil {
		panic(err)
	}
	Dictionary := make([]string, 0)
	for _, Word := range(bytes.Split(Data, []byte("\n"))) {
		Dictionary = append(Dictionary, string(Word))
	}
	InvertedSuffix := ferret.MakeInvertedSuffix(Dictionary, Conversion)
	fmt.Println(InvertedSuffix.Query("ar", 5, false, false, make([]byte, 0)))
	fmt.Println(InvertedSuffix.Query("test", 5, false, false, make([]byte, 0)))
	fmt.Println(InvertedSuffix.Query("tsst", 5, false, true, ferret.LowercaseLetters))
	fmt.Println(InvertedSuffix.Query("a", 5, true, false, make([]byte, 0)))
	fmt.Println(InvertedSuffix.Query("a", 5, false, false, make([]byte, 0)))
}