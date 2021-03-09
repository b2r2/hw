package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString return simple error by invalid string.
var (
	ErrInvalidString = errors.New("invalid string")
	spec             = '\\'
	isSpec           bool
	next             rune
)

// Unpack ...
func Unpack(in string) (string, error) {
	if len(in) == 0 {
		return "", nil
	}
	if len(in) == 1 {
		return string(in[0]), nil
	}
	if unicode.IsDigit(rune(in[0])) ||
		(spec == rune(in[0]) && len(in) == 2 && unicode.IsDigit(rune(in[1]))) {
		return "", ErrInvalidString
	}
	var out strings.Builder
	for inx, char := range in[:len(in)-1] {
		next = rune(in[inx+1])
		switch {
		case !isSpec && unicode.IsDigit(char):
			if unicode.IsDigit(next) {
				return "", ErrInvalidString
			}
		case !isSpec && char == spec:
			if !unicode.IsDigit(next) && next != spec {
				return "", ErrInvalidString
			}
			isSpec = true
			continue
		case unicode.IsDigit(next):
			count, _ := strconv.Atoi(string(next))
			out.WriteString(strings.Repeat(string(char), count))
		default:
			out.WriteRune(char)
		}
		isSpec = false
	}
	if isSpec || !unicode.IsDigit(next) {
		out.WriteRune(next)
	}
	return out.String(), nil
}
