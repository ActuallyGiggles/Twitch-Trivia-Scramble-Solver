package print

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/pterm/pterm"
)

func Print(p Instructions) {
	var Q string
	var QType string
	var AType string
	var A []string

	if p.Service == "trivia" {
		Q = p.Trivia.Question
		A = []string{p.Trivia.Answer}
		QType = "Question:    "
		AType = "Answer:      "
	} else if p.Service == "scramble" {
		Q = p.Scramble.Question
		A = p.Scramble.Matches
		QType = "Scrambled:   "
		AType = "Unscrambled: "
	}

	if p.NoteOnly && p.Error {
		pterm.Error.Printf("%s\n%s\n%s\n", strings.ToUpper(p.Channel), pterm.Gray(p.Service), pterm.Gray("Note: "+p.Note))
		pterm.Println()
	} else if p.NoteOnly && !p.Error {
		pterm.Info.Printf("%s\n%s\n%s\n", strings.ToUpper(p.Channel), pterm.Gray(p.Service), pterm.Gray("Note: "+p.Note))
		pterm.Println()
	} else {
		chunks := chunks(Q, 7)

		message := fmt.Sprintf("%-20s  |  %s", strings.ToUpper(p.Channel), QType)

		if utf8.RuneCountInString(Q) < 50 {
			message += fmt.Sprintf("%s\n", Q)
		} else {
			message += fmt.Sprintf("%s\n", chunks[0])
			for _, chunk := range chunks[1:] {
				message += fmt.Sprintf("%-20s  |               %s\n", "", chunk)
			}
		}

		if A[0] == "" || A == nil {
			message += fmt.Sprintf("%-20s           |  %s%s\n", pterm.Gray(p.Service), AType, "[N/A]")
		} else {
			message += fmt.Sprintf("%-20s           |  %s%s\n", pterm.Gray(p.Service), AType, A)
		}

		message += fmt.Sprintf("%-20s  |  %s\n\n", "", pterm.Gray("Note:        "+p.Note))

		pterm.Success.Println(message)
		pterm.Println()
	}
}

func Page(title string, content func()) {
	print("\033[H\033[2J")
	if title == "Exited" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightRed)).WithFullWidth().Println("Twitch Trivia/Scramble Autosolver by ActuallyGiggles")
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightRed)).WithFullWidth().Println(title)
	} else if title == "Started" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).WithFullWidth().Println("Twitch Trivia/Scramble Autosolver by ActuallyGiggles")
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).WithFullWidth().Println(title)
	} else if title == "Set up" {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithFullWidth().Println("Twitch Trivia/Scramble Autosolver by ActuallyGiggles")
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithFullWidth().Println(title)
	}
	pterm.Println()
	content()
}

func chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	split := strings.Split(s, " ")
	chunks := make([]string, 0, (len(split)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range split {
		if currentLen == chunkSize {
			chunks = append(chunks, strings.Join(split[currentStart:i], " "))
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, strings.Join(split[currentStart:], " "))
	return chunks
}

func Clear() {
	print("\033[H\033[2J")
}

type Instructions struct {
	Channel string

	Service  string
	Trivia   TriviaMode
	Scramble ScrambleMode

	Note     string
	NoteOnly bool

	Error bool
}

type TriviaMode struct {
	Question string
	Answer   string
}

type ScrambleMode struct {
	Question string
	Matches  []string
}
