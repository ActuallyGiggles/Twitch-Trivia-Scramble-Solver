package scramble

import (
	"encoding/json"
	"os"
	"sort"
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
	scrambledSlice := strings.Split(scrambledWord, "")
	sort.Strings(scrambledSlice)

wordLoop:
	for _, word := range words {
		if len(scrambledWord) != len(word) {
			continue wordLoop
		}

		word = strings.Trim(strings.ToLower(word), " ")
		wordSlice := strings.Split(word, "")
		sort.Strings(wordSlice)

		for i := 0; i < len(wordSlice); i++ {
			if wordSlice[i] != scrambledSlice[i] {
				continue wordLoop
			}
		}

		possibleMatches = append(possibleMatches, word)
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
