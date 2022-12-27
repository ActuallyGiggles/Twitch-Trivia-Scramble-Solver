package main

import (
	"strconv"
	"strings"
	"time"
	"twitch-trivia-unscrambler/config"
	"twitch-trivia-unscrambler/scramble"
	"twitch-trivia-unscrambler/trivia"
	"twitch-trivia-unscrambler/twitch"
)

var (
	workers = make(map[string]*worker)
)

type worker struct {
	Channel string

	TriviaUnknown  bool
	TriviaQuestion string
	TriviaCancel   chan bool

	ScrambleUnknown  bool
	ScrambleQuestion string
	ScrambleCancel   chan bool
}

func StartWorkers() {
	for _, channel := range config.Config.ChannelsToJoin {
		w := worker{
			Channel:        channel,
			TriviaCancel:   make(chan bool),
			ScrambleCancel: make(chan bool),
		}

		workers[channel] = &w
	}
}

func (w *worker) playScramble(message twitch.Message) {
	var split []string
	split = strings.Split(removeCommonEmotes(message.Message), " ")[1:]
	sentence := strings.Join(split, " ")
	var question string
	var answer string

	switch isType(sentence) {
	case "question":
		if isRandomlyRejected() {
			Print(printInstructions{
				Channel: message.Channel,
				Service: "scramble",
				Scramble: scrambleMode{
					Question: "",
					Matches:  []string{},
				},
				Note:     "randomly rejected",
				NoteOnly: true,
				Error:    false,
			})

			config.UpdateStats(w.Channel, "scramble", "rejected")
			return
		}

		question = split[len(split)-1]
		matches := scramble.Unscramble(question)

		if len(matches) > 0 {
			go w.answer(printInstructions{
				Channel: message.Channel,
				Service: "scramble",
				Scramble: scrambleMode{
					Question: question,
					Matches:  matches,
				},
				NoteOnly: false,
			})
			return
		}

		if len(matches) == 0 {
			Print(printInstructions{
				Channel: message.Channel,
				Service: "scramble",
				Scramble: scrambleMode{
					Question: question,
					Matches:  matches,
				},
				Note:     "scramble unknown",
				NoteOnly: true,
				Error:    true,
			})
			w.ScrambleQuestion = question
			w.ScrambleUnknown = true
			config.UpdateStats(w.Channel, "scramble", "unknown")
			return
		}
	case "answer":
		if w.ScrambleUnknown {
			answer = extractAnswer(sentence)
			scramble.AddWord(answer)

			Print(printInstructions{
				Channel: message.Channel,
				Service: "scramble",
				Scramble: scrambleMode{
					Question: question,
					Matches:  []string{},
				},
				Note:     "answer learned",
				NoteOnly: true,
				Error:    false,
			})

			w.ScrambleQuestion = ""
			w.ScrambleUnknown = false
			config.UpdateStats(w.Channel, "scramble", "learned")
			return
		}

		if !w.ScrambleUnknown && question != "" {
			w.ScrambleCancel <- true
		}
	}
}

func (w *worker) playTrivia(message twitch.Message) {
	var split []string
	split = strings.Split(removeCommonEmotes(message.Message), " ")[1:]
	sentence := strings.Join(split, " ")
	var question string
	var answer string

	switch isType(sentence) {
	case "question":
		if isRandomlyRejected() {
			Print(printInstructions{
				Channel: message.Channel,
				Service: "trivia",
				Trivia: triviaMode{
					Question: "",
					Answer:   "",
				},
				Note:     "randomly rejected",
				NoteOnly: true,
				Error:    false,
			})

			config.UpdateStats(w.Channel, "trivia", "rejected")
			return
		}

		question = sentence
		answer, answerFound := trivia.SearchTrivia(question)

		if answerFound {
			go w.answer(printInstructions{
				Channel: message.Channel,
				Service: "trivia",
				Trivia: triviaMode{
					Question: question,
					Answer:   answer,
				},
				NoteOnly: false,
			})
			return
		}

		if !answerFound {
			Print(printInstructions{
				Channel: message.Channel,
				Service: "trivia",
				Trivia: triviaMode{
					Question: question,
					Answer:   answer,
				},
				Note:     "answer unknown",
				NoteOnly: true,
				Error:    true,
			})
			w.TriviaQuestion = question
			w.TriviaUnknown = true
			config.UpdateStats(w.Channel, "trivia", "unknown")
			return
		}

	case "answer":
		if w.TriviaUnknown {
			answer = extractAnswer(sentence)
			trivia.AddTrivia(w.TriviaQuestion, answer)

			Print(printInstructions{
				Channel: message.Channel,
				Service: "trivia",
				Trivia: triviaMode{
					Question: question,
					Answer:   answer,
				},
				Note:     "answer learned",
				NoteOnly: true,
				Error:    false,
			})

			w.TriviaQuestion = ""
			w.TriviaUnknown = false
			config.UpdateStats(w.Channel, "trivia", "learned")
			return
		}

		if !w.TriviaUnknown && question != "" {
			w.TriviaCancel <- true
		}
	}
}

func (w *worker) answer(p printInstructions) {
	var eta int

	if config.Config.AnswerInterval.Min == config.Config.AnswerInterval.Max {
		eta = config.Config.AnswerInterval.Min
	} else {
		eta = RandomNumber(config.Config.AnswerInterval.Min, config.Config.AnswerInterval.Max)
	}

	if p.Service == "trivia" {
		config.UpdateStats(w.Channel, "trivia", "answered")

		ticker := time.NewTicker(time.Duration(eta) * time.Second)
		defer ticker.Stop()

		ticker2 := time.NewTicker(time.Duration(eta/2) * time.Second)
		defer ticker2.Stop()

		p.Note = "answering in " + strconv.Itoa(eta) + " seconds"
		Print(p)

		for {
			select {
			case <-w.TriviaCancel:
				p.Note = "answered by different user"
				p.NoteOnly = true
				p.Error = false
				Print(p)
				return
			case <-ticker.C:
				twitch.Say(w.Channel, strings.ToLower(p.Trivia.Answer))
				return
			case <-ticker2.C:
				if isPartialAnswerFirst() {
					split := strings.Split(p.Trivia.Answer, " ")

					if len(split) > 1 {
						twitch.Say(w.Channel, strings.ToLower(strings.Join(split[:1], " ")))
					} else if len(split) > 2 {
						twitch.Say(w.Channel, strings.ToLower(strings.Join(split[:2], " ")))
					}
				}
				return
			}
		}
	}

	if p.Service == "scramble" {
		config.UpdateStats(w.Channel, "scramble", "answered")

		ticker := time.NewTicker(time.Duration(eta) * time.Second)
		p.Note = "answering in " + strconv.Itoa(eta) + " seconds"
		Print(p)

		match := 0

		for {
			select {
			case <-w.ScrambleCancel:
				p.Note = "answering cancelled"
				p.NoteOnly = true
				p.Error = false
				Print(p)
				return
			case <-ticker.C:
				twitch.Say(w.Channel, strings.ToLower(p.Scramble.Matches[match]))
				match++
				if match == len(p.Scramble.Matches) {
					ticker.Stop()
					return
				}
			}
		}
	}
}
