package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if domain == "" {
		return nil, errors.New("domain received empty")
	}
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	for scanner.Scan() {
		email := fastjson.GetString(scanner.Bytes(), "Email")
		if strings.Contains(email, domain) {
			result[strings.ToLower(strings.SplitN(email, "@", 2)[1])]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
