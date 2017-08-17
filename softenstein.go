package main

import (
	"log"
	"os"
	"regexp"

	"github.com/nlopes/slack"
	"time"
	"fmt"
)

var rtm *slack.RTM
func main() {
	api := slack.New(os.Getenv("SLACK_API_NICOKOSI_TOKEN"))
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	//api.SetDebug(true)

	rtm = api.NewRTM()
	go rtm.ManageConnection()
	buildCommand := regexp.MustCompile("build (.*)")
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			buildArgs := buildCommand.FindStringSubmatch(ev.Text)
			if len(buildArgs) > 1 {
				build(buildArgs, ev.Channel, ev.Timestamp)
			}
		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())
		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return
		}
	}
}
func build(buildArgs []string, channel string, timestamp string) {
	rtm.SendMessage(ThreadedOutgoingMessage(channel, "Going to build\n> "+buildArgs[1], timestamp))
	time.Sleep(2 * time.Second)
	rtm.SendMessage(ThreadedOutgoingMessage(channel, "Done!", timestamp))
	rtm.AddReaction("white_check_mark", slack.NewRefToMessage(channel, timestamp))
}
func ThreadedOutgoingMessage(channel string, text string, timestamp string) *slack.OutgoingMessage {
	msg := rtm.NewOutgoingMessage(text, channel)
	msg.ThreadTimestamp = timestamp
	return msg
}
