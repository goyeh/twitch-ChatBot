#!/bin/bash

UPSTREAM=${1:-'@{u}'}
LOCAL=$(git rev-parse @)
REMOTE=$(git rev-parse "$UPSTREAM")
BASE=$(git merge-base @ "$UPSTREAM")

if [ $LOCAL = $REMOTE ]; then
    echo "Up-to-date"
elif [ $LOCAL = $BASE ]; then
    echo "Need to pull"
    git pull
elif [ $REMOTE = $BASE ]; then
    echo "Need to push"
    git commit -m "Autmatic commit"
    git push
    git pull
else
    echo "Diverged; Mmanually solve conflicts"
fi

go get github.com/vikpe/twitch-chatbot
go mod tidy

# force pi build.
# pi->env GOOS=linux GOARCH=arm GOARM=5 go build
go build
./twitchbot
