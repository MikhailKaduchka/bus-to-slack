#!/bin/sh

mkdir build
GOOS=linux go build -o build/main src/departures/app.go

rm -f build/departures2Slack.zip
zip build/departures2Slack.zip build/main
