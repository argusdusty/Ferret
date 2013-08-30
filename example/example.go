// Copyright 2013 Mark Canning
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// Author: Mark Canning
// Developed at: Tamber, Inc. (http://www.tamber.com/).
//
// Tamber also has this really cool recommendation engine for music
// (also development by me) which prioritizes up-and-coming artists, so
// it doesn't succomb to the popularity biases that plague modern
// recommendation engines, and still produces excellent personalized
// recommendations! Make sure to check us out at http://www.tamber.com
// or https://itunes.apple.com/us/app/tamber-concerts/id658240483

package main

import (
	"fmt"
	"github.com/argusdusty/Ferret"
)

// Some words to search for and return
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

// Length of data
var ExampleData = []interface{}{
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

var ExampleCorrection = func(b []byte) [][]byte { return ferret.ErrorCorrect(b, ferret.LowercaseLetters) }
var ExampleSorter = func(s string, v interface{}, l int, i int) float64 { return -float64(l + i) }
var ExampleConverter = func(s string) []byte { return []byte(s) }

func main() {
	ExampleSearchEngine := ferret.New(ExampleWords, ExampleWords, ExampleData, ExampleConverter)
	fmt.Println(ExampleSearchEngine.Query("ar", 5))
	fmt.Println(ExampleSearchEngine.Query("test", 5))
	fmt.Println(ExampleSearchEngine.ErrorCorrectingQuery("tsst", 5, ExampleCorrection))
	fmt.Println(ExampleSearchEngine.SortedErrorCorrectingQuery("tsst", 5, ExampleCorrection, ExampleSorter))
	fmt.Println(ExampleSearchEngine.SortedQuery("a", 5, ExampleSorter))
	fmt.Println(ExampleSearchEngine.Query("a", 5))
	ExampleSearchEngine.Insert("asdfghjklqwertyuiopzxcvbnm", "asdfghjklqwertyuiopzxcvbnm", []uint64{26})
	fmt.Println(ExampleSearchEngine.Query("sdfghjklqwert", 5))
	fmt.Println(ExampleSearchEngine.Query("ferret", 5))
	ExampleSearchEngine.Insert("ferret", "ferret", []uint64{7})
	fmt.Println(ExampleSearchEngine.Query("ferret", 5))
}
