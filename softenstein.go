package main

import (
	"log"
	"os"
	"regexp"

	"github.com/nlopes/slack"
	"time"
	"fmt"
	"net/http"
	"io"
)

var rtm *slack.RTM
func main() {
	api := slack.New(os.Getenv("SOFTENSTEIN_SLACK_API_TOKEN"))
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	//api.SetDebug(true)

	http.HandleFunc("/build", buildServer)

	rtm = api.NewRTM()
	go rtm.ManageConnection()
	go http.ListenAndServe(":8080", nil)


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

func buildServer(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		io.WriteString(w, "Built!")
	} else {
		w.WriteHeader(405)
	}
}