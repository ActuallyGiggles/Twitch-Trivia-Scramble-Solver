package trivia

import (
	"encoding/json"
	"os"
	"regexp"
)

var (
	t            []trivia
	triviaLoaded bool
	regex        = regexp.MustCompile(`\<(.*?)\>`)
)

type trivia struct {
	Source   string `json:"source"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// Example trivia dialogue:
//
// [Trivia] After playing the nightclub circuit, he broke into movies with "What's New, Pussycat?"
// [Trivia] Hint: Wood_ _____
// [Trivia] @speedking1994 warren is close. [Similarity: 63%]
// [Trivia] No one answered correctly. FeelsBadMan The answer was: " Woody Allen "
// [Trivia] @porohappy You answered the question correctly and got 10 points. FeelsGoodMan The answer was " News Anchor ". [Similarity: 100%]

func SearchTrivia(question string) (answer string, found bool) {
	if !triviaLoaded {
		loadTrivia()
		triviaLoaded = true
	}

	for _, trivia := range t {
		if question == trivia.Question {
			return trivia.Answer, true
		}
	}
	return "", false
}

func AddTrivia(question string, answer string) {
	t = append(t, trivia{Source: "new", Question: question, Answer: answer})
	saveTrivia()
}

func saveTrivia() {
	f, err := os.OpenFile("./trivia.json", os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	c, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		panic(err)
	}

	_, err = f.Write(c)
	if err != nil {
		panic(err)
	}
}

func loadTrivia() {
	triviaFile, err := os.Open("./trivia.json")
	if err != nil {
		panic(err)
	}
	defer triviaFile.Close()

	err = json.NewDecoder(triviaFile).Decode(&t)
	if err != nil {
		panic(err)
	}
}
