package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"twitch-trivia-unscrambler/config"
	"twitch-trivia-unscrambler/twitch"
)

var regex = regexp.MustCompile(`"(.*?)"`)
var commonEmotesRegex = regexp.MustCompile(`KEKWait\s?|FeelsDankMan\s?|GuitarTime\s?|TeaTime\s?|FeelsGoodMan\s?|FeelsBadMan\s?|DankgG\s?`)

func isTrivia(message twitch.Message) bool {
	if message.Author == "amazefulbot" && strings.Contains(message.Message, "[Trivia]") {
		if strings.Contains(message.Message, "Hint:") || strings.Contains(message.Message, "is close.") {
			return false
		}
		return true
	}
	return false
}

func isScramble(message twitch.Message) bool {
	if message.Author == "amazefulbot" && strings.Contains(message.Message, "[Scramble]") {
		if strings.Contains(message.Message, "Hint:") || strings.Contains(message.Message, "is close.") {
			return false
		}
		return true
	}
	return false
}

func isType(sentence string) string {
	if strings.Contains(sentence, "The answer was") || strings.Contains(sentence, "The word was") || strings.Contains(sentence, "@") {
		return "answer"
	}
	return "question"
}

func RandomNumber(min, max int) int {
	var result int
	switch {
	case min > max:
		// Fail with error
		return result
	case max == min:
		result = max
	case max > min:
		maxRand := max - min
		b, err := rand.Int(rand.Reader, big.NewInt(int64(maxRand)))
		if err != nil {
			return result
		}
		result = min + int(b.Int64())
	}
	return result
}

func extractAnswer(message string) string {
	return strings.Trim(regex.FindStringSubmatch(message)[1], " ")
}

func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
}

func removeCommonEmotes(message string) string {
	return commonEmotesRegex.ReplaceAllString(message, "")
}

func isRandomlyRejected() bool {
	num := RandomNumber(0, 100)

	if num > config.Config.RandomRejectionPercentage {
		return false
	}

	return true
}

func isPartialAnswerFirst(answer string) bool {
	if len(strings.Split(answer, " ")) == 1 {
		return false
	}

	num := RandomNumber(0, 100)

	if num > config.Config.PartiallyAnswerPercentage {
		return false
	}

	return true
}
