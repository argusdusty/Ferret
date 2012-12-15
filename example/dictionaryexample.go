package main

import (
	"io/ioutil"
	"bytes"
	"fmt"
	"github.com/argusdusty/ferret"
	"time"
	"sort"
)

var Conversion = func(s string) []byte { return []byte(s) }
var Correction = func(b []byte) [][]byte { return ferret.ErrorCorrect(b, ferret.LowercaseASCII) }

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
	fmt.Println(InvertedSuffix.Query("ar", 5, false))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.Query("test", 5, false))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.ErrorCorrectingQuery("tsst", 5, false, Correction))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.Query("a", 5, true))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.Query("a", 5, false))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.Query("the", 25, false))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	InvertedSuffix.Insert("asdfghjklqwertyuiopzxcvbnm")
	fmt.Println("Performed insert in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(InvertedSuffix.Query("sdfghjklqwert", 5, false))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()

	LSIS := ferret.MakeLengthSortedInvertedSuffix(Dictionary, Conversion)
	fmt.Println("Created LSIS in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.Query("ar", 5))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.Query("test", 5))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.ErrorCorrectingQuery("tsst", 5, Correction))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.Query("a", 5))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.Query("a", 5))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.Query("the", 25))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
	LSIS.Insert("asdfghjklqwertyuiopzxcvbnm")
	fmt.Println("Performed insert in:", time.Now().Sub(t)); t = time.Now()
	fmt.Println(LSIS.Query("sdfghjklqwert", 5))
	fmt.Println("Performed search in:", time.Now().Sub(t)); t = time.Now()
}