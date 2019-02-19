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

package ferret

// AllASCII is all ASCII bytes (0-127)
var AllASCII = []byte{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 17, 18,
	19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35,
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52,
	53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69,
	70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86,
	87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102,
	103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116,
	117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
}

// PrintableASCII is all printable ASCII bytes
var PrintableASCII = []byte{
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48,
	49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65,
	66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82,
	83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
	100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113,
	114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126,
}

// LowercaseASCII is all printable ASCII bytes excluding capitalized letters (A-Z)
var LowercaseASCII = []byte{
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48,
	49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65,
	91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105,
	106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118,
	119, 120, 121, 122, 123, 124, 125, 126,
}

// LowercaseLetters is All lowercase ASCII bytes (a-z/97-122)
var LowercaseLetters = []byte{
	97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
	111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
}

// ErrorCorrect returns all byte-arrays which are Levenshtein distance of 1 away from Word
// within an allowed array of byte characters.
func ErrorCorrect(Word []byte, AllowedBytes []byte) [][]byte {
	results := make([][]byte, 0)
	N := len(Word)
	for i := 0; i < N; i++ {
		t := Word[i]
		// Remove Character
		temp := make([]byte, N)
		copy(temp, Word)
		temp = append(temp[:i], temp[i+1:]...)
		results = append(results, temp)
		if i != 0 {
			// Add Character
			for _, c := range AllowedBytes {
				temp := make([]byte, N)
				copy(temp, Word)
				temp = append(temp[:i], append([]byte{c}, temp[i:]...)...)
				results = append(results, temp)
			}
			// Transpose Character
			temp := make([]byte, N)
			copy(temp, Word)
			temp[i], temp[i-1] = temp[i-1], temp[i]
			results = append(results, temp)
		}
		// Insert Character
		for _, c := range AllowedBytes {
			if c == t {
				continue
			}
			temp := make([]byte, N)
			copy(temp, Word)
			temp[i] = c
			results = append(results, temp)
		}
	}
	return results
}
