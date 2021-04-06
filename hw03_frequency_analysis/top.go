package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type word struct {
	word string
	freq int
}

var (
	rg          = regexp.MustCompile(`[a-zA-ZА-Яа-я]+`)
	punctuation = regexp.MustCompile(`[.,!?]`)
)

func getWords(s []word) []string {
	arr := make([]string, 0, len(s))
	for _, v := range s {
		arr = append(arr, v.word)
	}
	return arr
}

func setWords(res map[string]int) (words []word) {
	for w, f := range res {
		words = append(words, word{w, f})
	}

	sort.Slice(words, func(i, j int) bool {
		if words[i].freq > words[j].freq {
			return true
		} else if words[i].freq < words[j].freq {
			return false
		}
		return words[i].word < words[j].word
	})
	return
}

func Top10(in string) []string {
	if len(in) == 0 {
		return nil
	}
	res := make(map[string]int)
	for _, v := range strings.Fields(in) {
		v = strings.ToLower(v)
		if punctuation.MatchString(v) {
			v = v[:len(v)-1]
		}
		if rg.MatchString(v) {
			res[v]++
		}
	}
	words := setWords(res)
	if len(words) >= 10 {
		return getWords(words[:10])
	}
	return getWords(words)
}
