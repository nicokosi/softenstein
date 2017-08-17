package main

import (
	"log"
	"os"
	"regexp"

	"github.com/nlopes/slack"
	"time"
	"fmt"
)

func main() {
	api := slack.New(os.Getenv("SLACK_API_NICOKOSI_TOKEN"))
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	//api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	buildCommand := regexp.MustCompile("build (.*)")

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:
			buildArgs := buildCommand.FindStringSubmatch(ev.Text)
			if len(buildArgs) > 1 {
				msg := rtm.NewOutgoingMessage(
					"Going to build\n> "+buildArgs[1],
					ev.Channel)
				msg.ThreadTimestamp = ev.Timestamp
				rtm.SendMessage(msg)

				time.Sleep(2 * time.Second)

				msg2 := rtm.NewOutgoingMessage("Done!", ev.Channel)
				msg2.ThreadTimestamp = msg.ThreadTimestamp
				rtm.SendMessage(msg2)
				rtm.AddReaction("white_check_mark", slack.NewRefToMessage(ev.Channel, ev.Timestamp))
			}

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return
		}
	}
}
