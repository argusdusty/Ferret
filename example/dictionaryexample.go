package main

import (
	"io/ioutil"
	"bytes"
	"fmt"
	"github.com/argusdusty/ferret"
	"time"
)

var Conversion = func(s string) []byte { return []byte(s) }

func main() {
	t := time.Now()
	Data, err := ioutil.ReadFile("dictionary.dat")
	if err != nil {
		panic(err)
	}
	Dictionary := make([]string, 0)
	for _, Word := range(bytes.Split(Data, []byte("\r\n"))) {
		Dictionary = append(Dictionary, string(Word))
	}
	fmt.Println("Loaded dictionary in:", time.Now().Sub(t)); t = time.Now()
	InvertedSuffix := ferret.MakeInvertedSuffix(Dictionary, Conversion)
	fmt.Println("Created index in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.Query("ar", 5, false, false, make([]byte, 0)))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.Query("test", 5, false, false, make([]byte, 0)))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.Query("tsst", 5, false, true, ferret.LowercaseLetters))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.Query("a", 5, true, false, make([]byte, 0)))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.Query("a", 5, false, false, make([]byte, 0)))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.Query("the", 25, false, false, make([]byte, 0)))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	LSIS := ferret.MakeLengthSortedInvertedSuffix(Dictionary, Conversion)
	fmt.Println("Created LSIS in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.Query("ar", 5, 100))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.Query("test", 5, 100))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.Query("a", 5, 100))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.Query("a", 5, 100))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.Query("the", 25, 100))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
}