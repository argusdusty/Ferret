package ferret

import (
	"strings"
	"unicode"
)

// Only handles Lowercase Latin-1 Supplement
var UnicodeToASCII = map[rune]rune {
	'À':'A', 'Á':'A', 'Â':'A', 'Ã':'A', 'Ä':'A', 'Å':'A',
	//'Æ':'AE',
	'Ç':'C',
	'È':'E', 'É':'E', 'Ê':'E', 'Ë':'E',
	'Ì':'I', 'Í':'I', 'Î':'I', 'Ï':'I',
	//'Ð':'D',
	'Ñ':'N',
	'Ò':'O', 'Ó':'O', 'Ô':'O', 'Õ':'O', 'Ö':'O', 'Ø':'O',
	'Ù':'U', 'Ú':'U', 'Û':'U', 'Ü':'U',
	'Ý':'Y',
	'ß':'B',
	'à':'a', 'á':'a', 'â':'a', 'ã':'a', 'ä':'a', 'å':'a',
	//'æ':'ae',
	'ç':'c',
	'è':'e', 'é':'e', 'ê':'e', 'ë':'e',
	'ì':'i', 'í':'i', 'î':'i', 'ï':'i',
	'ð':'o',
	'ñ':'n',
	'ò':'o', 'ó':'o', 'ô':'o', 'õ':'o', 'ö':'o', 'ø':'o',
	'ù':'u', 'ú':'u', 'û':'u', 'ü':'u',
	'ý':'y', 'ÿ':'y',
}

func ToASCII(r rune) rune {
	a, ok := UnicodeToASCII[r]
	if ok { return a }
	return r
}

func UnicodeToLowerASCII(s string) []byte {
	return []byte(strings.Map(func(r rune) rune { return ToASCII(unicode.ToLower(r)) }, s))
}
