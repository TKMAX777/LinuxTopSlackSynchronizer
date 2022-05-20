package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/TKMAX777/LinuxTopSlackSynchronizer/top"
)

func main() {
	n, err := strconv.Atoi(os.Getenv("PROCESS_NUMBER"))
	if err != nil {
		n = 20
	}

	var Slack = NewSlackHandler(os.Getenv("SLACK_TOKEN"), n)
	var regularChannel = os.Getenv("SLACK_REGULAR_CHANNEL")

	var AlertChannel = os.Getenv("SLACK_ALERT_CHANNEL")
	var AlertMode = os.Getenv("SLACK_ALERT_MODE") == "on"
	AlertLevel, _ := strconv.Atoi(os.Getenv("SLACK_ALERT_LEVEL"))

	SlackSendInterval, _ := strconv.Atoi(os.Getenv("SLACK_ALERT_LEVEL"))
	if SlackSendInterval == 0 {
		SlackSendInterval = 1
	}

	Slack.RemoveLastRegularMessage(regularChannel)

	for {
		ps, err := top.Get()
		if err != nil {
			log.Println("Error: Get: ", err)
			time.Sleep(time.Second)
			continue
		}

		err = Slack.RegularSend(ps, regularChannel)
		if err != nil {
			log.Println("Error: RegularSend: ", err)
		}

		var cpu float32
		for _, p := range ps {
			cpu += p.CPU
		}

		if AlertMode && cpu > float32(AlertLevel) {
			Slack.Send(ps, AlertChannel)
		}

		time.Sleep(time.Second * time.Duration(SlackSendInterval))
	}
}
