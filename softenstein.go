package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New(os.Getenv("SLACK_API_NICOKOSI_TOKEN"))
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	//api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:
			helloCmd, _ := regexp.MatchString("hello *", ev.Text)
			if helloCmd {
				rtm.SendMessage(rtm.NewOutgoingMessage("Namaste", ev.Channel))
			} else {
				rtm.SendMessage(rtm.NewOutgoingMessage(
					"Hmm, I don't know what to do. My supported command is\n> hello <any message>",
					ev.Channel))
			}
		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return
		}
	}
}
