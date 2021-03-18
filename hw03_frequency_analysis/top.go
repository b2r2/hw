package hw03_frequency_analysis

import (
	"regexp"
	"sort"
	"strings"
)

var rg = map[string]*regexp.Regexp{
	"cyrillic":    regexp.MustCompile(`\p{Cyrillic}+`),
	"punctuation": regexp.MustCompile(`[.,!?@#$%^&*_]`),
}

type wordContainer struct {
	word  string
	count int
}

type words struct {
	container  []wordContainer
	countWords int
}

func newWords(size int) *words {
	return &words{
		container: make([]wordContainer, size),
	}
}

func (w words) findWord(s string) int {
	if rg["punctuation"].Match([]byte(s)) {
		s = s[:len(s)-1]
	}
	for i, v := range w.container {
		if v.word == s {
			return i
		}
	}
	return -1
}

func handleString(in string) []string {
	tmp := strings.ReplaceAll(in, "\n", " ")
	return strings.Split(strings.ReplaceAll(tmp, "\t", " "), " ")
}

func Top10(in string) (out []string) {
	if len(in) == 0 || !rg["cyrillic"].Match([]byte(in)) {
		return nil
	}
	s := handleString(in)
	words := newWords(len(s))
	for i, v := range s {
		v = strings.ToLower(v)
		if rg["cyrillic"].Match([]byte(v)) {
			if position := words.findWord(v); position != -1 {
				words.container[position].count++
			} else {
				if rg["punctuation"].Match([]byte(v)) {
					v = v[:len(v)-1]
				}
				words.container[i].word = v
				words.container[i].count++
			}
		}
		words.countWords++
	}
	if words.countWords < 10 {
		return nil
	}
	sort.Slice(words.container, func(i, j int) bool {
		if words.container[i].count > words.container[j].count {
			return true
		} else if words.container[i].count < words.container[j].count {
			return false
		}
		return words.container[i].word < words.container[j].word
	})

	for i := 0; i < 10; i++ {
		out = append(out, words.container[i].word)
	}
	return
}
