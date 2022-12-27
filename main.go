package main

import (
	"os"
	"os/signal"
	"syscall"
	"twitch-trivia-unscrambler/config"
	"twitch-trivia-unscrambler/twitch"
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	config.GetConfig()
	StartWorkers()

	chat := make(chan twitch.Message)

	go twitch.Start(chat)

	go func() {
		for message := range chat {
			if isTrivia(message) && config.Config.DoTrivia {
				go workers[message.Channel].playTrivia(message)
			}

			if isScramble(message) && config.Config.DoScramble {
				go workers[message.Channel].playScramble(message)
			}
		}
	}()

	page("Started", func() {})

	<-sc
	page("Exited", func() {})
}
