package scramble

import (
	"encoding/json"
	"os"
	"strings"
)

var (
	words       []string
	wordsLoaded bool
	suffixes    = []string{"", "ed", "s", "es", "ing", "er", "ers", "y", "ally"}
)

func Unscramble(scrambledWord string) (possibleMatches []string) {
	if !wordsLoaded {
		loadWords()
		wordsLoaded = true
	}

	scrambledWord = strings.Trim(strings.ToLower(scrambledWord), " ")

wordLoop:
	for _, word := range words {
		word = strings.ToLower(word)
		wordSplit := strings.Split(word, "")
		var matchingLetters int

		for _, wordLetter := range wordSplit {
			if strings.Count(scrambledWord, wordLetter) == strings.Count(word, wordLetter) {
				matchingLetters++
			} else {
				continue wordLoop
			}
		}

		if matchingLetters == len(scrambledWord) {
			possibleMatches = append(possibleMatches, word)
		}
	}

	return possibleMatches
}

func AddWord(word string) {
	words = append(words, word)
	saveWords()
}

func loadWords() {
	f, err := os.Open("./scramble.json")
	defer f.Close()
	if err != nil {
		return
	}

	err = json.NewDecoder(f).Decode(&words)
	if err != nil {
		panic(err)
	}
}

func saveWords() {
	f, err := os.OpenFile("./scramble.json", os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	c, err := json.MarshalIndent(words, "", " ")
	if err != nil {
		panic(err)
	}

	_, err = f.Write(c)
	if err != nil {
		panic(err)
	}
}
