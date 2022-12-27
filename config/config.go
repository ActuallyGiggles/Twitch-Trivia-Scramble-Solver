package config

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
)

var Config configuration

type configuration struct {
	AccountName    string
	AccountOAuth   string
	ChannelsToJoin []string

	DoTrivia   bool
	DoScramble bool

	AnswerInterval            numberRange
	RandomRejectionPercentage int
	PartiallyAnswerPercentage int

	Stats []s
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

func GetConfig() {
	f, err := os.Open("./config.json")
	defer f.Close()
	if err != nil {
		configSetup()
		return
	}

	err = json.NewDecoder(f).Decode(&Config)
	if err != nil {
		panic(err)
	}
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
	page(func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Enter the Twitch account name you will be using.\n"))
		pterm.Print(pterm.LightBlue("	--Account Name: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		name := strings.ToLower(scanner.Text())
		Config.AccountName = name
	})

	// Account OAuth
	page(func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Print(pterm.LightBlue("Obtaining your OAuth is necessary to connect to Twitch chat as yourself.\nHere is a link to get it: ", pterm.Underscore.Sprintf("https://twitchapps.com/tmi/\n")))
		pterm.Print(pterm.LightBlue("	--Account OAuth: "), pterm.White("oauth:"))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		oauth := strings.ToLower(scanner.Text())
		Config.AccountOAuth = oauth
	})

	// Channels to join
	page(func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Do you want to answer Trivia questions (t), Scramble questions (s), or both (b)?\n"))
		pterm.Print(pterm.LightBlue("	--Answer: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		which := strings.Split(strings.ToLower(scanner.Text()), " ")
		if len(which) > 1 {
			pterm.Println()
			pterm.Error.Println("Only answer with t for Trivia, s for scramble, or b for both.")
			os.Exit(3)
		}
		if which[0] == "t" || which[0] == "b" {
			Config.DoTrivia = true
		}
		if which[0] == "s" || which[0] == "b" {
			Config.DoScramble = true
		}
	})

	// Channels to join
	page(func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the channels in which the program should act in.\nSeparate channel names with spaces.\n"))
		pterm.Print(pterm.LightBlue("	--Channels To Join: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		channels := strings.Split(strings.ToLower(scanner.Text()), " ")
		Config.ChannelsToJoin = channels
	})

	// Time range
	page(func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.LightBlue("Specify the time interval in which to respond.\nThe time is enacted in seconds.\nSeparate with a space.\n"))
		pterm.Print(pterm.LightBlue("	--Interval: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		r := strings.Split(scanner.Text(), " ")
		min, err := strconv.Atoi(r[0])
		if err != nil {
			pterm.Println()
			pterm.Error.Println(r[0], "is not a number.")
			os.Exit(3)
		}
		max, err := strconv.Atoi(r[1])
		if err != nil {
			pterm.Println()
			pterm.Error.Println(r[1], "is not a number!")
			os.Exit(3)
		}
		Config.AnswerInterval.Min = min
		Config.AnswerInterval.Max = max
	})

	// Random rejection
	page(func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Print(pterm.LightBlue("Specify what percentage of questions should purposefully be ignored."))
		pterm.Print(pterm.LightBlue("	--Percentage: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		num := scanner.Text()
		percentage, err := strconv.Atoi(num)
		if err != nil {
			pterm.Println()
			pterm.Error.Println(num, "is not a number.")
			os.Exit(3)
		}
		Config.RandomRejectionPercentage = percentage
	})

	// Partially answer percentage
	page(func() {
		pterm.DefaultCenter.WithCenterEachLineSeparately().Print(pterm.LightBlue("Specify what percentage of trivia questions should purposefully be partially answered first."))
		pterm.Print(pterm.LightBlue("	--Percentage: "))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		num := scanner.Text()
		percentage, err := strconv.Atoi(num)
		if err != nil {
			pterm.Println()
			pterm.Error.Println(num, "is not a number.")
			os.Exit(3)
		}
		Config.PartiallyAnswerPercentage = percentage
	})

	saveConfig()
}

func page(content func()) {
	print("\033[H\033[2J")
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithFullWidth().Println("Twitch Trivia/Scramble Autosolver by ActuallyGiggles")
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithFullWidth().Println("First Time Setup")
	pterm.Println()
	pterm.Println()
	content()
	print("\033[H\033[2J")
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
