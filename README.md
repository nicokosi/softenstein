# softenstein

Another Slack bot, for fun (sure) and profit (maybe)! 🤖

## Setup

Create a bot in your Slack team https://your_Slack_team.slack.com.

Install [Go](https://golang.org/dl/)  then:
```bash
go get
go install

export SOFTENSTEIN_SLACK_API_TOKEN=<your secret Slack API token>
softenstein
```
The bot:
- replies to messages like `build this please`
- executes commands like `/build this please`


In order to test commands on localhost, open a secured tunnel to localhost via ngrok:
```bash
ngrok http 8080
```
Then copy/paste forwarded URL (e.g. http://f96a4adc.ngrok.io)
and configure it in Slack (`Custom Integrations` > `Slash Commands`):
http://f96a4adc.ngrok.io/build