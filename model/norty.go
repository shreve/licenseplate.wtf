package model

// Catch norty words

import (
	"os"
	"regexp"
	"strings"
)

var nortyWords []string

var isAlphaNum = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`).MatchString

func getNortyWords() []string {
	// Read list of norty words from environment variable and split by newlines
	nortyWords := strings.ToUpper(os.Getenv("NORTY_WORDS"))
	return strings.Split(nortyWords, "\n")
}

func isNorty(word string) bool {
	upper := strings.ToUpper(word)

	// Check if word is in any of the words in the list of norty words
	for _, nortyWord := range nortyWords {
		if strings.Contains(nortyWord, upper) {
			return true
		}
	}
	return false
}

func init() {
	nortyWords = getNortyWords()
}
