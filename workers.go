package main

import (
	"strconv"
	"strings"
	"time"
	"twitch-trivia-unscrambler/config"
	"twitch-trivia-unscrambler/print"
	"twitch-trivia-unscrambler/scramble"
	"twitch-trivia-unscrambler/trivia"
	"twitch-trivia-unscrambler/twitch"
)

var (
	workers = make(map[string]*worker)
)

type worker struct {
	Channel string

	TriviaKnown    bool
	TriviaQuestion string
	TriviaCancel   chan bool

	ScrambleKnown    bool
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
		question = split[len(split)-1]
		matches := scramble.Unscramble(question)

		if len(matches) > 0 {
			if isRandomlyRejected() {
				print.Print(print.Instructions{
					Channel: message.Channel,
					Service: "scramble",
					Scramble: print.ScrambleMode{
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

			go w.answer(print.Instructions{
				Channel: message.Channel,
				Service: "scramble",
				Scramble: print.ScrambleMode{
					Question: question,
					Matches:  matches,
				},
				NoteOnly: false,
			})
			w.ScrambleQuestion = ""
			w.ScrambleKnown = true
			return
		}

		if len(matches) == 0 {
			print.Print(print.Instructions{
				Channel: message.Channel,
				Service: "scramble",
				Scramble: print.ScrambleMode{
					Question: question,
					Matches:  matches,
				},
				Note:     "scramble unknown",
				NoteOnly: true,
				Error:    true,
			})
			w.ScrambleQuestion = question
			w.ScrambleKnown = false
			config.UpdateStats(w.Channel, "scramble", "unknown")
			return
		}
	case "answer":
		if !w.ScrambleKnown {
			answer = extractAnswer(sentence)
			scramble.AddWord(answer)

			print.Print(print.Instructions{
				Channel: message.Channel,
				Service: "scramble",
				Scramble: print.ScrambleMode{
					Question: question,
					Matches:  []string{},
				},
				Note:     "answer learned",
				NoteOnly: true,
				Error:    false,
			})

			w.ScrambleQuestion = ""
			w.ScrambleKnown = true
			config.UpdateStats(w.Channel, "scramble", "learned")
			return
		}

		if w.ScrambleKnown && w.ScrambleQuestion == "" {
			w.ScrambleCancel <- true
			return
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
		question = sentence
		answer, answerFound := trivia.SearchTrivia(question)

		if answerFound {
			if isRandomlyRejected() {
				print.Print(print.Instructions{
					Channel: message.Channel,
					Service: "trivia",
					Trivia: print.TriviaMode{
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

			go w.answer(print.Instructions{
				Channel: message.Channel,
				Service: "trivia",
				Trivia: print.TriviaMode{
					Question: question,
					Answer:   answer,
				},
				NoteOnly: false,
			})
			w.TriviaQuestion = ""
			w.TriviaKnown = true
			return
		}

		if !answerFound {
			print.Print(print.Instructions{
				Channel: message.Channel,
				Service: "trivia",
				Trivia: print.TriviaMode{
					Question: question,
					Answer:   answer,
				},
				Note:     "answer unknown",
				NoteOnly: true,
				Error:    true,
			})
			w.TriviaQuestion = question
			w.TriviaKnown = false
			config.UpdateStats(w.Channel, "trivia", "unknown")
			return
		}

	case "answer":
		if !w.TriviaKnown {
			answer = extractAnswer(sentence)
			trivia.AddTrivia(w.TriviaQuestion, answer)

			print.Print(print.Instructions{
				Channel: message.Channel,
				Service: "trivia",
				Trivia: print.TriviaMode{
					Question: question,
					Answer:   answer,
				},
				Note:     "answer learned",
				NoteOnly: true,
				Error:    false,
			})

			w.TriviaQuestion = ""
			w.TriviaKnown = true
			config.UpdateStats(w.Channel, "trivia", "learned")
			return
		}

		if w.TriviaKnown && w.TriviaQuestion == "" {
			w.TriviaCancel <- true
			return
		}
	}
}

func (w *worker) answer(p print.Instructions) {
	var eta int

	if config.Config.AnswerInterval.Min == config.Config.AnswerInterval.Max {
		eta = config.Config.AnswerInterval.Max
	} else {
		eta = RandomNumber(config.Config.AnswerInterval.Min, config.Config.AnswerInterval.Max)
	}

	p.Note = "answering in " + strconv.Itoa(eta) + " seconds"
	print.Print(p)

	if p.Service == "trivia" {
		config.UpdateStats(w.Channel, "trivia", "answered")

		if eta == 1 {
			time.Sleep(1 * time.Second)
			twitch.Say(w.Channel, strings.ToLower(p.Trivia.Answer))
			return
		}

		ticker := time.NewTicker(time.Duration(eta) * time.Second)
		defer ticker.Stop()

		ticker2 := time.NewTicker(time.Duration(eta/2) * time.Second)
		defer ticker2.Stop()

		for {
			select {
			case <-w.TriviaCancel:
				p.Note = "answered by different user"
				p.NoteOnly = true
				p.Error = false
				print.Print(p)
				return
			case <-ticker.C:
				twitch.Say(w.Channel, strings.ToLower(p.Trivia.Answer))
				return
			case <-ticker2.C:
				if isPartialAnswerFirst(p.Trivia.Answer) {
					split := strings.Split(p.Trivia.Answer, " ")

					if len(split) > 1 {
						twitch.Say(w.Channel, strings.ToLower(strings.Join(split[:1], " ")))
					} else if len(split) > 2 {
						twitch.Say(w.Channel, strings.ToLower(strings.Join(split[:2], " ")))
					}
				}
				continue
			}
		}
	}

	if p.Service == "scramble" {
		config.UpdateStats(w.Channel, "scramble", "answered")

		if eta == 1 {
			time.Sleep(1 * time.Second)
			twitch.Say(w.Channel, strings.ToLower(p.Trivia.Answer))
			return
		}

		ticker := time.NewTicker(time.Duration(eta) * time.Second)

		match := 0

		for {
			select {
			case <-w.ScrambleCancel:
				p.Note = "answering cancelled"
				p.NoteOnly = true
				p.Error = false
				print.Print(p)
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
