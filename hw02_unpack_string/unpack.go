package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString return simple error by invalid string.
var (
	ErrInvalidString        = errors.New("invalid string")
	slash            string = `\`
)

// Unpack ...
func Unpack(in string) (string, error) {
	if len(in) == 0 {
		return "", nil
	}
	var out strings.Builder
	for inx, char := range in {
		if unicode.IsDigit(char) && inx == 0 {
			return "", ErrInvalidString
		}
		if inx > 0 {
			prev := in[inx-1]
			if unicode.IsDigit(char) && unicode.IsDigit(rune(prev)) {
				return "", ErrInvalidString
			}
			if unicode.IsDigit(char) || string(char) == slash {
				count, _ := strconv.Atoi(string(char))
				tmp := strings.Repeat(string(prev), count)
				s := out.String()
				out.Reset()
				out.WriteString(s[:len(s)-1] + tmp)
			}
		}
		if !unicode.IsDigit(char) {
			out.WriteRune(char)
		}
	}
	return out.String(), nil
}
