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
				rtm.SendMessage(ThreadedOutgoingMessage(ev.Channel, "Going to build\n> "+buildArgs[1], ev.Timestamp))
				build(buildArgs)
				rtm.SendMessage(ThreadedOutgoingMessage(ev.Channel, "Done!", ev.Timestamp))
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

func build(args []string) {
	time.Sleep(2 * time.Second)
	fmt.Print("Fake build finished ")
	if len(args) > 0 {
		fmt.Sprintf("(args=%v)\n", args[1:])
	} else {
		fmt.Println("(no args)")
	}
}

func ThreadedOutgoingMessage(channel string, text string, timestamp string) *slack.OutgoingMessage {
	msg := rtm.NewOutgoingMessage(text, channel)
	msg.ThreadTimestamp = timestamp
	return msg
}

func buildServer(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		build([]string{})
		io.WriteString(w, "Built!")
	} else {
		w.WriteHeader(405)
	}
}
