package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type WordContainer struct {
	Word  string
	Count int
}

type Words struct {
	container []WordContainer
}

func New(size int) *Words {
	return &Words{
		container: make([]WordContainer, size),
	}
}

func (w Words) FindWord(s string) int {
	for i, v := range w.container {
		if v.Word == s {
			return i
		}
	}
	return -1
}

func handleString(in string) []string {
	tmp := strings.ReplaceAll(in, "\n", " ")
	return strings.Split(strings.ReplaceAll(tmp, "\t", " "), " ")
}

func sortSlice(list []string) []string {
	sort.Strings(list)
	return list
}

func uniqueSlice(list []string) []string {
	keys := make(map[string]bool)
	var slice []string
	for _, v := range list {
		if _, ok := keys[v]; !ok {
			keys[v] = true
			slice = append(slice, v)
		}
	}
	return sortSlice(slice)
}

func unique(m map[int][]string) map[int][]string {
	for i := range m {
		m[i] = uniqueSlice(m[i])
	}
	return m
}

func Top10(in string) (out []string) {
	if len(in) == 0 {
		return []string{}
	}
	rg := regexp.MustCompile(`-|\p{Cyrillic}+`)
	s := handleString(in)
	words := New(len(s))
	for i, v := range s {
		if rg.Match([]byte(v)) {
			if position := words.FindWord(v); position != -1 {
				words.container[position].Count++
			} else {
				words.container[i].Word = v
				words.container[i].Count++
			}
		}
	}
	rate := make(map[int][]string)
	sort.Slice(words.container, func(i, j int) bool {
		if words.container[i].Count > words.container[j].Count {
			rate[words.container[i].Count] = append(rate[words.container[i].Count], words.container[i].Word)
			return true
		}
		return false
	})
	uRate := unique(rate)
	for i := len(words.container); ; i-- {
		out = append(out, uRate[i]...)
		if len(out) == 10 {
			break
		}
	}
	return
}
