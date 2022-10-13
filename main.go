package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"twitchbot/config"
	"twitchbot/lib"

	"github.com/Adeithe/go-twitch"
	"github.com/Adeithe/go-twitch/irc"
)

var (
	Username string
	Oauth    string
	Channel  string
)

/***********************
 *      _       _ _
 *     (_)     (_) |
 *      _ _ __  _| |_
 *     | | '_ \| | __|
 *     | | | | | | |_
 *     |_|_| |_|_|\__| */
func init() {
	Username = config.Val.Username
	Oauth = config.Val.Oauth
	Channel = config.Val.Channel
}

/***********************************
 *                      _
 *                     (_)
 *      _ __ ___   __ _ _ _ __
 *     | '_ ` _ \ / _` | | '_ \
 *     | | | | | | (_| | | | | |
 *     |_| |_| |_|\__,_|_|_| |_| */
func main() {
	LoadResponces()
	LoadBot()
}

/*************************************************
 *      _                     _ ____        _
 *     | |                   | |  _ \      | |
 *     | |     ___   __ _  __| | |_) | ___ | |_
 *     | |    / _ \ / _` |/ _` |  _ < / _ \| __|
 *     | |___| (_) | (_| | (_| | |_) | (_) | |_
 *     |______\___/ \__,_|\__,_|____/ \___/ \__| */
func LoadBot() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	writer := &irc.Conn{}
	writer.SetLogin(Username, Oauth)
	if err := writer.Connect(); err != nil {
		panic("failed to start writer")
	}

	Send("BotMaster Has Arrived...")

	reader := twitch.IRC()
	reader.OnShardReconnect(onShardReconnect)
	reader.OnShardLatencyUpdate(onShardLatencyUpdate)
	reader.OnShardMessage(onShardMessage)
	reader.OnShardChannelJoin(onJoined)

	if err := reader.Join(Channel); err != nil {
		panic(err)
	}
	lib.Info("Connected to Twitch IRC")

	<-sc
	fmt.Println("Stopping...")
	reader.Close()
	writer.Close()
}

/**********************************************
 *                       _       _
 *                      | |     (_)
 *       ___  _ __      | | ___  _ _ __
 *      / _ \| '_ \ _   | |/ _ \| | '_ \
 *     | (_) | | | | |__| | (_) | | | | |
 *      \___/|_| |_|\____/ \___/|_|_| |_| *  */
func onJoined(shardID int, msg1 string, msg2 string) {
	Send(fmt.Sprintf("Welcome %s to %s Channel! Don't forget to subscribe. it helps us a lot!", msg2, msg1))
	fmt.Printf("Shard #%s & #%s connected\n", msg1, msg2)
}

/****************************************************************************************************
 *                   _____ _                   _ _____                                      _
 *                  / ____| |                 | |  __ \                                    | |
 *       ___  _ __ | (___ | |__   __ _ _ __ __| | |__) |___  ___ ___  _ __  _ __   ___  ___| |_
 *      / _ \| '_ \ \___ \| '_ \ / _` | '__/ _` |  _  // _ \/ __/ _ \| '_ \| '_ \ / _ \/ __| __|
 *     | (_) | | | |____) | | | | (_| | | | (_| | | \ \  __/ (_| (_) | | | | | | |  __/ (__| |_
 *      \___/|_| |_|_____/|_| |_|\__,_|_|  \__,_|_|  \_\___|\___\___/|_| |_|_| |_|\___|\___|\__| * */
func onShardReconnect(shardID int) {
	fmt.Printf("Shard #%d reconnected\n", shardID)
}

/*************************************************************************************************************************
 *                   _____ _                   _ _           _                        _    _           _       _
 *                  / ____| |                 | | |         | |                      | |  | |         | |     | |
 *       ___  _ __ | (___ | |__   __ _ _ __ __| | |     __ _| |_ ___ _ __   ___ _   _| |  | |_ __   __| | __ _| |_ ___
 *      / _ \| '_ \ \___ \| '_ \ / _` | '__/ _` | |    / _` | __/ _ \ '_ \ / __| | | | |  | | '_ \ / _` |/ _` | __/ _ \
 *     | (_) | | | |____) | | | | (_| | | | (_| | |___| (_| | ||  __/ | | | (__| |_| | |__| | |_) | (_| | (_| | ||  __/
 *      \___/|_| |_|_____/|_| |_|\__,_|_|  \__,_|______\__,_|\__\___|_| |_|\___|\__, |\____/| .__/ \__,_|\__,_|\__\___|
 *                                                                               __/ |      | |
 *                                                                              |___/       |_|                         */
func onShardLatencyUpdate(shardID int, latency time.Duration) {
	fmt.Printf("Shard #%d has %dms ping\n", shardID, latency.Milliseconds())
}

/*****************************************************************************************
 *                   _____ _                   _ __  __
 *                  / ____| |                 | |  \/  |
 *       ___  _ __ | (___ | |__   __ _ _ __ __| | \  / | ___  ___ ___  __ _  __ _  ___
 *      / _ \| '_ \ \___ \| '_ \ / _` | '__/ _` | |\/| |/ _ \/ __/ __|/ _` |/ _` |/ _ \
 *     | (_) | | | |____) | | | | (_| | | | (_| | |  | |  __/\__ \__ \ (_| | (_| |  __/
 *      \___/|_| |_|_____/|_| |_|\__,_|_|  \__,_|_|  |_|\___||___/___/\__,_|\__, |\___|
 *                                                                           __/ |
 *                                                                          |___/       */
func onShardMessage(shardID int, msg irc.ChatMessage) {
	fmt.Printf("#%s %s: %s\n", msg.Channel, msg.Sender.DisplayName, msg.Text)
	switch strings.ToUpper(msg.Text) {
	case "HELLO":
		Send(fmt.Sprint("Welcome, and Hello back to you ", msg.Sender.Username))
	case "BYE":
		Send("You leaving now, We will miss you, Hope your back soon!")
	case "GG":
		Send("Was a great game, Loved playing with you too!")
	default:
		if len(msg.Text) > 3 && strings.ToUpper(msg.Text[0:3]) == "BOT" {
			BotCommands(msg)
		} else if strings.Contains(strings.ToUpper(msg.Text), "I WAS THE KILLER") {
			Send("Great run, Enjoyed the Challenge!")
		} else if strings.Contains(strings.ToUpper(msg.Text), "SUCK") {
			Send("Sorry, but Sucking is not in my Service agreements, Try a different message!")
		}
	}
}

/***********************************************************************************
 *      ____        _    _____                                          _
 *     |  _ \      | |  / ____|                                        | |
 *     | |_) | ___ | |_| |     ___  _ __ ___  _ __ ___   __ _ _ __   __| |___
 *     |  _ < / _ \| __| |    / _ \| '_ ` _ \| '_ ` _ \ / _` | '_ \ / _` / __|
 *     | |_) | (_) | |_| |___| (_) | | | | | | | | | | | (_| | | | | (_| \__ \
 *     |____/ \___/ \__|\_____\___/|_| |_| |_|_| |_| |_|\__,_|_| |_|\__,_|___/  * * */
func BotCommands(msg irc.ChatMessage) {
	command := msg.Text[4:]
	switch strings.ToUpper(command) {
	case "HELLO":
		Send(fmt.Sprint("Hello, This is BotMaster, What can I do for you ", msg.Sender.Username))
	case "BYE":
		Send("You want me to leave? What did I do wrong? You are mean! Very Mean!!")
	case "COMMANDS", "HELP":
		Send("To enter a Bot command, type the following:")
		Send("bot [command] [press enter]")
		Send("Example: bot hello ⤶")
		Send("List:")
		Send("Commands or help")
	default:
		Send("Command not found!! Type Commands for a list!")
	}
}

/*****************************************
 *       _____                _
 *      / ____|              | |
 *     | (___   ___ _ __   __| |
 *      \___ \ / _ \ '_ \ / _` |
 *      ____) |  __/ | | | (_| |
 *     |_____/ \___|_| |_|\__,_|  * *  */
func Send(message string) {
	sendWriter := &irc.Conn{}
	sendWriter.SetLogin(Username, Oauth)
	if err := sendWriter.Connect(); err != nil {
		panic("failed to start writer")
	}
	defer sendWriter.Close()
	sendWriter.Say(Channel, message)
	time.Sleep(200 * time.Millisecond)
}

/********************************************************************************
 *      _                     _ _____
 *     | |                   | |  __ \
 *     | |     ___   __ _  __| | |__) |___  ___ _ __   ___  _ __   ___ ___  ___
 *     | |    / _ \ / _` |/ _` |  _  // _ \/ __| '_ \ / _ \| '_ \ / __/ _ \/ __|
 *     | |___| (_) | (_| | (_| | | \ \  __/\__ \ |_) | (_) | | | | (_|  __/\__ \
 *     |______\___/ \__,_|\__,_|_|  \_\___||___/ .__/ \___/|_| |_|\___\___||___/
 *                                             | |
 *                                             |_|                               */
type MsgFormat struct {
	Msg   string `json:"msg"`
	Reply string `json:"reply"`
}
type FullMsg struct {
	Text   []MsgFormat `json:"text"`
	Intext []MsgFormat `json:"intext"`
}

func LoadResponces() {
	fileContent, err := os.Open("responces.json")
	if err != nil {
		lib.Error(err)
		return
	}
	lib.Info("Json Responces File loaded!")
	defer fileContent.Close()
	byteResult, _ := ioutil.ReadAll(fileContent)
	var fullMsg FullMsg
	json.Unmarshal(byteResult, &fullMsg)
	for i := 0; i < len(fullMsg.Text); i++ {
		lib.Debug("Message: " + fullMsg.Text[i].Msg)
		lib.Debug("Reply: " + fullMsg.Text[i].Reply)
	}
	for i := 0; i < len(fullMsg.Intext); i++ {
		lib.Debug("Message: " + fullMsg.Intext[i].Msg)
		lib.Debug("Reply: " + fullMsg.Intext[i].Reply)
	}
}
