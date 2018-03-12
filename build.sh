#!/bin/sh

GOOS=linux go build src/departures2Slack.go

rm -f departures2Slack.zip
zip departures2Slack.zip departures2Slack
