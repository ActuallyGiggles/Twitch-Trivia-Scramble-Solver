package twitch

import (
	"twitch-trivia-unscrambler/config"

	"github.com/gempir/go-twitch-irc/v3"
)

var client *twitch.Client

// Start creates a twitch client and connects it.
func Start(in chan Message) {
	client = &twitch.Client{}
	client = twitch.NewClient(config.Config.AccountName, "oauth:"+config.Config.AccountOAuth)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		m := Message{
			Channel: message.Channel,
			Author:  message.User.Name,
			Message: message.Message,
		}

		in <- m
	})

	for _, channel := range config.Config.ChannelsToJoin {
		Join(channel)
	}

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

// Say sends a message to a specific twitch chatroom.
func Say(channel string, message string) {
	client.Say(channel, message)
}

// Join joins a twitch chatroom.
func Join(channel string) {
	client.Join(channel)
}

// Depart departs a twitch chatroom.
func Depart(channel string) {
	client.Depart(channel)
}

// Twitch message struct
type Message struct {
	Channel string
	Author  string
	Message string
}
