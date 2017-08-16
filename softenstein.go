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
			buildCommand := regexp.MustCompile("build (.*)")
			buildArgs := buildCommand.FindStringSubmatch(ev.Text)
			if len(buildArgs) > 1 {
				rtm.SendMessage(
					rtm.NewOutgoingMessage(
						"Going to build\n> " + buildArgs[1],
						ev.Channel))
			} else {
				rtm.SendMessage(
					rtm.NewOutgoingMessage(
						"Hmm, I don't know what to do. My supported command is\n> build <any message>",
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
