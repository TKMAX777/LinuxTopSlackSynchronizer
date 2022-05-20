# LinuxTopSlackSynchronizer

## About

This program uses the top command to display programs with high CPU utilization in Slack.


## Install

```sh
go install github.com/TKMAX777/LinuxTopSlackSynchronizer@latest
```

## Usage

Set the following environment variables

```
PROCESS_NUMBER='Number of processes to display'
SLACK_TOKEN='xoxb-*****'
SLACK_REGULAR_CHANNEL='Slack channel ID'
SLACK_ALERT_CHANNEL='Slack channel ID'
SLACK_ALERT_MODE='on / off'
SLACK_ALERT_LEVEL='CPU utilization when sending Alert'
SLACK_SEND_INTERVAL='Interval for updating Slack messages [second]'
```
