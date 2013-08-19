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

package ferret

import (
	"strings"
	"unicode"
)

// Only handles Latin-1 Supplement
var UnicodeToASCII = map[rune]rune{
	'À': 'A', 'Á': 'A', 'Â': 'A', 'Ã': 'A', 'Ä': 'A', 'Å': 'A',
	'Ç': 'C',
	'È': 'E', 'É': 'E', 'Ê': 'E', 'Ë': 'E',
	'Ì': 'I', 'Í': 'I', 'Î': 'I', 'Ï': 'I',
	'Ñ': 'N',
	'Ò': 'O', 'Ó': 'O', 'Ô': 'O', 'Õ': 'O', 'Ö': 'O', 'Ø': 'O',
	'Ù': 'U', 'Ú': 'U', 'Û': 'U', 'Ü': 'U',
	'Ý': 'Y',
	'à': 'a', 'á': 'a', 'â': 'a', 'ã': 'a', 'ä': 'a', 'å': 'a',
	'ç': 'c',
	'è': 'e', 'é': 'e', 'ê': 'e', 'ë': 'e',
	'ì': 'i', 'í': 'i', 'î': 'i', 'ï': 'i',
	'ð': 'o',
	'ñ': 'n',
	'ò': 'o', 'ó': 'o', 'ô': 'o', 'õ': 'o', 'ö': 'o', 'ø': 'o',
	'ù': 'u', 'ú': 'u', 'û': 'u', 'ü': 'u',
	'ý': 'y', 'ÿ': 'y',
}

func ToASCII(r rune) rune {
	a, ok := UnicodeToASCII[r]
	if ok {
		return a
	}
	return r
}

func UnicodeToLowerASCII(s string) []byte {
	return []byte(strings.Map(func(r rune) rune { return unicode.ToLower(ToASCII(r)) }, s))
}
