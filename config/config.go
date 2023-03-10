package config

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"twitch-trivia-unscrambler/print"

	"github.com/pterm/pterm"
)

var Config configuration

type configuration struct {
	AccountName  string
	AccountOAuth string

	Play whichToPlay

	ChannelsToJoin []string

	AnswerInterval            numberRange
	RandomRejectionPercentage int
	PartiallyAnswerPercentage int

	Stats []s
}

type whichToPlay struct {
	Trivia   bool
	Scramble bool
}

type numberRange struct {
	Min int
	Max int
}

type s struct {
	Channel  string
	Trivia   statistics
	Scramble statistics
}

type statistics struct {
	Answered int
	Rejected int
	Unknown  int
}

func GetConfig() bool {
	f, err := os.Open("./config.json")
	defer f.Close()
	if err != nil {
		configSetup()
		return false
	}

	err = json.NewDecoder(f).Decode(&Config)
	if err != nil {
		panic(err)
	}

	return true
}

func saveConfig() {
	f, err := os.OpenFile("./config.json", os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	c, err := json.MarshalIndent(Config, "", " ")
	if err != nil {
		panic(err)
	}

	_, err = f.Write(c)
	if err != nil {
		panic(err)
	}
}

func configSetup() {
	// Account name
	print.Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Enter the Twitch account name you will be using.\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Account Name: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		name := strings.ToLower(scanner.Text())
		Config.AccountName = name
		return true
	})

	// Account OAuth
	print.Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Print(pterm.LightBlue("Obtaining your OAuth is necessary to connect to Twitch chat as yourself.\nHere is a link to get it: ", pterm.Underscore.Sprintf("https://twitchapps.com/tmi/\n")))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Account OAuth: "), pterm.White("oauth:"))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		oauth := strings.ToLower(scanner.Text())
		Config.AccountOAuth = oauth
		return true
	})

	// Play which games
	print.Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Do you want to answer Trivia questions (t), Scramble questions (s), or both (b)?\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Answer: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		which := strings.Split(strings.ToLower(scanner.Text()), " ")
		if len(which) > 1 || (which[0] != "s" && which[0] != "b" && which[0] != "t") {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println("Only answer with t for Trivia, s for scramble, or b for both.")
			return false
		}
		if which[0] == "t" || which[0] == "b" {
			Config.Play.Trivia = true
		}
		if which[0] == "s" || which[0] == "b" {
			Config.Play.Scramble = true
		}
		return true
	})

	// Channels to join
	print.Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the channels in which the program should act in.\nSeparate channel names with spaces.\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Channels To Join: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		channels := strings.Split(strings.ToLower(scanner.Text()), " ")
		if strings.Trim(channels[0], " ") == "" {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println("No channels entered!")
			return false
		}
		Config.ChannelsToJoin = channels
		return true
	})

	// Time range
	print.Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the time interval in which to respond.\nThe time is enacted in seconds.\nSeparate with a space.\nMinimum of 1 second.\n"))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Interval: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		r := strings.Split(scanner.Text(), " ")
		min, err := strconv.Atoi(r[0])
		if err != nil {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println(r[0], "is not a number!")
			return false
		}
		max, err := strconv.Atoi(r[1])
		if err != nil {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println(r[1], "is not a number!")
			return false
		}
		if min < 1 || max < 1 {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println("Cannot be less than 1!")
			return false
		}
		if min > max {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println("Minimum cannot be greater than maximum!")
			return false
		}
		Config.AnswerInterval.Min = min
		Config.AnswerInterval.Max = max
		return true
	})

	// Random rejection
	print.Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Print(pterm.LightBlue("Specify what percentage of questions should purposefully be ignored.\nMust be between 0 and 100."))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Percentage: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		num := scanner.Text()
		percentage, err := strconv.Atoi(num)
		if err != nil {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println(num, "is not a number.")
			return false
		}
		if percentage < 0 || percentage > 100 {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println("Number must be between 0 and 100!")
			return false
		}
		Config.RandomRejectionPercentage = percentage
		return true
	})

	// Partially answer percentage
	print.Page("Set Up", func() bool {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Print(pterm.LightBlue("Specify what percentage of trivia questions should purposefully be partially answered first.\nMust be between 0 and 100."))
		pterm.Println()
		pterm.Print(pterm.LightBlue("	--Percentage: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		num := scanner.Text()
		percentage, err := strconv.Atoi(num)
		if err != nil {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println(num, "is not a number.")
			return false
		}
		if percentage < 0 || percentage > 100 {
			pterm.Println()
			pterm.Println()
			pterm.Error.Println("Number must be between 0 and 100!")
			return false
		}
		Config.PartiallyAnswerPercentage = percentage
		return true
	})

	saveConfig()
	pterm.Println()
}

func UpdateStats(channel string, service string, stat string) {
	exists := false

	for i := 0; i < len(Config.Stats); i++ {
		c := &Config.Stats[i]

		if c.Channel == channel {
			exists = true

			if stat == "answered" && service == "trivia" {
				c.Trivia.Answered++
			}
			if stat == "rejected" && service == "trivia" {
				c.Trivia.Rejected++
			}
			if stat == "unknown" && service == "trivia" {
				c.Trivia.Unknown++
			}

			if stat == "answered" && service == "scramble" {
				c.Scramble.Answered++
			}
			if stat == "rejected" && service == "scramble" {
				c.Scramble.Rejected++
			}
			if stat == "unknown" && service == "scramble" {
				c.Scramble.Unknown++
			}
		}
	}

	if !exists {
		if stat == "answered" && service == "trivia" {
			Config.Stats = append(Config.Stats, s{Channel: channel, Trivia: statistics{Answered: 1}})
		}
		if stat == "rejected" && service == "trivia" {
			Config.Stats = append(Config.Stats, s{Channel: channel, Trivia: statistics{Rejected: 1}})
		}
		if stat == "unknown" && service == "trivia" {
			Config.Stats = append(Config.Stats, s{Channel: channel, Trivia: statistics{Unknown: 1}})
		}

		if stat == "answered" && service == "scramble" {
			Config.Stats = append(Config.Stats, s{Channel: channel, Scramble: statistics{Answered: 1}})
		}
		if stat == "rejected" && service == "scramble" {
			Config.Stats = append(Config.Stats, s{Channel: channel, Scramble: statistics{Rejected: 1}})
		}
		if stat == "unknown" && service == "scramble" {
			Config.Stats = append(Config.Stats, s{Channel: channel, Scramble: statistics{Unknown: 1}})
		}
	}

	saveConfig()
}
