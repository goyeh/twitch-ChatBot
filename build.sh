#!/bin/bash
go get github.com/vikpe/twitch-chatbot
go mod tidy
# pi->env GOOS=linux GOARCH=arm GOARM=5 go build
go build
./twitchbot
