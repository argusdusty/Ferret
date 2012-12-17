package ferret

// Not ready yet
/*
// All ASCII
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

// All printable ASCII
var PrintableASCII = []byte{
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48,
	49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65,
	66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82,
	83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
	100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113,
	114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126,
}
*/

// All printable ASCII excluding capitalized letters (A-Z)
var LowercaseASCII = []byte{
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48,
	49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65,
	91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105,
	106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118,
	119, 120, 121, 122, 123, 124, 125, 126,
}

// All lowercase ASCII (a-z)
var LowercaseLetters = []byte{
	97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
	111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
}

// Returns all byte-arrays which are Levenshtein distance of 1 away from Word
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
