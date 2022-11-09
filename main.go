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

/***********************************
 *                      _
 *                     (_)
 *      _ __ ___   __ _ _ _ __
 *     | '_ ` _ \ / _` | | '_ \
 *     | | | | | | (_| | | | | |
 *     |_| |_| |_|\__,_|_|_| |_| */
func main() {
	Username = config.Val.Username
	Oauth = config.Val.Oauth
	Channel = config.Val.Channel
	lib.Info("Connecting ", Username, " to ", Channel)

	LoadJokes()
	LoadResponces()
	LoadBot()
}

func MonitorServer(reader *irc.Conn) {
	//for {
	lib.Debug("Channel:", reader.UserState.UserState)
	//	time.Sleep(time.Duration(5) * time.Second)
	//}
}

/*************************************************
 *      _                     _ ____        _
 *     | |                   | |  _ \      | |
 *     | |     ___   __ _  __| | |_) | ___ | |_
 *     | |    / _ \ / _` |/ _` |  _ < / _ \| __|
 *     | |___| (_) | (_| | (_| | |_) | (_) | |_
 *     |______\___/ \__,_|\__,_|____/ \___/ \__| */
func LoadBot() {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error Loading Bot:", r)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	writer := &irc.Conn{}
	writer.SetLogin(Username, Oauth)
	if err := writer.Connect(); err != nil {
		panic("failed to start writer")
	}
	writer.Close()

	Send("BotMaster Has Arrived...")

	reader := twitch.IRC()
	defer reader.Close()
	reader.OnShardReconnect(onShardReconnect)
	reader.OnShardLatencyUpdate(onShardLatencyUpdate)
	reader.OnShardMessage(onShardMessage)
	reader.OnShardChannelJoin(onJoined)
	reader.OnShardChannelLeave(onLeave)
	reader.OnShardChannelUpdate(onServerUpdate)
	reader.OnShardServerNotice(onServerNotice)
	if err := reader.Join(Channel); err != nil {
		panic(err)
	}
	lib.Info("Connected to Twitch IRC")

	<-sc
	fmt.Println("Stopping...")
}

func onServerUpdate(serverID int, msg irc.RoomState) {
	lib.Debug("Server Update:", msg)
}

func onServerNotice(serverID int, msg irc.ServerNotice) {
	lib.Debug("Server Notices:", msg)
}

/**********************************************
 *                       _       _
 *                      | |     (_)
 *       ___  _ __      | | ___  _ _ __
 *      / _ \| '_ \ _   | |/ _ \| | '_ \
 *     | (_) | | | | |__| | (_) | | | | |
 *      \___/|_| |_|\____/ \___/|_|_| |_| *  */
func onJoined(shardID int, msg1 string, msg2 string) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error When user Joins:", r)
		}
	}()

	for i := range TwitchMsg.Events {
		if strings.ToUpper(TwitchMsg.Events[i].Msg) == "JOIN" {
			textToSend := TwitchMsg.Events[i].Reply[lib.RandomRange(0, len(TwitchMsg.Events[i].Reply))]
			textToSend = strings.Replace(textToSend, "{name}", msg2, -1)
			textToSend = strings.Replace(textToSend, "{channel}", msg1, -1)
			Send(textToSend)
			break
		}
	}
	lib.Info("Shard", msg1, " & ", msg2, " connected\n")
}

/**********************************************
 *                       _       _
 *                      | |     (_)
 *       ___  _ __      | | ___  _ _ __
 *      / _ \| '_ \ _   | |/ _ \| | '_ \
 *     | (_) | | | | |__| | (_) | | | | |
 *      \___/|_| |_|\____/ \___/|_|_| |_| *  */
func onLeave(shardID int, msg1 string, msg2 string) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error When user Leavs:", r)
		}
	}()

	for i := range TwitchMsg.Events {
		if strings.ToUpper(TwitchMsg.Events[i].Msg) == "LEAVE" {
			textToSend := TwitchMsg.Events[i].Reply[lib.RandomRange(0, len(TwitchMsg.Events[i].Reply))]
			textToSend = strings.Replace(textToSend, "{name}", msg2, -1)
			textToSend = strings.Replace(textToSend, "{channel}", msg1, -1)
			Send(textToSend)
			break
		}
	}
	lib.Info("Shard ", msg1, " & ", msg2, " connected\n")
}

/****************************************************************************************************
 *                   _____ _                   _ _____                                      _
 *                  / ____| |                 | |  __ \                                    | |
 *       ___  _ __ | (___ | |__   __ _ _ __ __| | |__) |___  ___ ___  _ __  _ __   ___  ___| |_
 *      / _ \| '_ \ \___ \| '_ \ / _` | '__/ _` |  _  // _ \/ __/ _ \| '_ \| '_ \ / _ \/ __| __|
 *     | (_) | | | |____) | | | | (_| | | | (_| | | \ \  __/ (_| (_) | | | | | | |  __/ (__| |_
 *      \___/|_| |_|_____/|_| |_|\__,_|_|  \__,_|_|  \_\___|\___\___/|_| |_|_| |_|\___|\___|\__| * */
func onShardReconnect(shardID int) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error Reconnecting:", r)
		}
	}()

	lib.Info("Shard ", shardID, " reconnected\n")
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
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error Getting Ping:", r)
		}
	}()

	lib.Info("Shard ", shardID, " has ", latency.Milliseconds(), "ping\n")
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
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error When Receiving Message:", r)
		}
	}()

	if config.Val.Username != msg.Sender.Username {
		lib.Info("Channel", msg.Channel, msg.Sender.DisplayName, msg.Text)

		if len(msg.Text) > 0 && strings.ToUpper(msg.Text[0:1]) == ":" {
			lib.Debug("Bot Command")
			BotCommands(msg)
		} else if FullMessage(msg) {
		} else if SomeWords(msg) {
			lib.Debug("Word Check")
		} else {
			lib.Debug("Nothing to do")
		}
	} else {
		lib.Info("Ignore message from Self", msg.Text)
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
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error BotCommands:", r)
		}
	}()

	command := msg.Text[1:]
	lib.Debug("Command Requested:", command)
	switch strings.ToUpper(command) {
	case "HELLO":
		Send(fmt.Sprint("Hello, This is BotMaster, What can I do for you ", msg.Sender.Username))
	case "BYE":
		Send("You want me to leave? What did I do wrong? You are mean! Very Mean!!")
	case "COMMANDS", "HELP":
		Send("To enter a Bot command, type the following:")
		Send(":[command] [press enter]")
		Send("Example: :hello ⤶")
		Send("List:")
		Send("Commands or help")
		Send("Joke for a Joke")
		Send("Insult for a Insults")
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
	defer func() {
		r := recover()
		if r != nil {
			lib.Error("Error Sending:", r)
		}
	}()

	lib.Debug("Message to Send:", message)
	sendWriter := &irc.Conn{}
	sendWriter.SetLogin(Username, Oauth)
	if err := sendWriter.Connect(); err != nil {
		panic("failed to start writer")
	}
	defer sendWriter.Close()
	sendWriter.Say(Channel, message)
	time.Sleep(200 * time.Millisecond)
	lib.Debug("Message Sent:", message)
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
	Msg   string   `json:"msg"`
	Reply []string `json:"reply"`
}
type FullMsg struct {
	Sentences []MsgFormat `json:"sentences"`
	Words     []MsgFormat `json:"words"`
	Events    []MsgFormat `json:"events"`
}

var TwitchMsg FullMsg

func LoadResponces() bool {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error Load Responces:", r)
		}
	}()
	fileContent, err := os.Open("responces.json")
	if err != nil {
		lib.Error(err)
		return false
	}
	lib.Info("Json Responces File loaded!")
	defer fileContent.Close()
	byteResult, _ := ioutil.ReadAll(fileContent)
	json.Unmarshal(byteResult, &TwitchMsg)

	return true
}

/************_*******_********************
 *          | |     | |
 *          | | ___ | | _____  ___
 *      _   | |/ _ \| |/ / _ \/ __|
 *     | |__| | (_) |   <  __/\__ \
 *      \____/ \___/|_|\_\___||___/ * * */
type MsgJokes struct {
	Jokes []string `json:"jokes"`
}

var FullJokes MsgJokes

func LoadJokes() bool {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error Load Jokes:", r)
		}
	}()
	fileContent, err := os.Open("jokes.json")
	if err != nil {
		lib.Error(err)
		return false
	}
	lib.Info("Json Jokes File loaded!")
	defer fileContent.Close()
	byteResult, _ := ioutil.ReadAll(fileContent)
	json.Unmarshal(byteResult, &FullJokes)
	return true
}

/*******_____*****************_**********************_*******_**************
 *     |  __ \               | |                    | |     | |
 *     | |__) |__ _ _ __   __| | ___  _ __ ___      | | ___ | | _____
 *     |  _  // _` | '_ \ / _` |/ _ \| '_ ` _ \ _   | |/ _ \| |/ / _ \
 *     | | \ \ (_| | | | | (_| | (_) | | | | | | |__| | (_) |   <  __/
 *     |_|  \_\__,_|_| |_|\__,_|\___/|_| |_| |_|\____/ \___/|_|\_\___|    */
func RandomJokes() string {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error Choosing Random Joke:", r)
		}
	}()

	return FullJokes.Jokes[lib.RandomRange(0, len(FullJokes.Jokes))]
}

/*******______*****_*_*__**__*************************************
 *     |  ____|   | | |  \/  |
 *     | |__ _   _| | | \  / | ___  ___ ___  __ _  __ _  ___
 *     |  __| | | | | | |\/| |/ _ \/ __/ __|/ _` |/ _` |/ _ \
 *     | |  | |_| | | | |  | |  __/\__ \__ \ (_| | (_| |  __/
 *     |_|   \__,_|_|_|_|  |_|\___||___/___/\__,_|\__, |\___|
 *                                                 __/ |
 *                                                |___/       */
func FullMessage(message irc.ChatMessage) (retVal bool) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error Reading Full Message:", r)
		}
	}()
	retVal = false
	for i := range TwitchMsg.Sentences {
		textToSee := strings.ToUpper(message.Text)
		if strings.ToUpper(TwitchMsg.Sentences[i].Msg) == textToSee {
			textToSend := TwitchMsg.Sentences[i].Reply[lib.RandomRange(0, len(TwitchMsg.Sentences[i].Reply))]
			Send(Substitution(textToSend, message))
			retVal = true
			break
		}

	}
	return retVal
}

/********_____*******************__**********__***********_**********
 *      / ____|                  \ \        / /          | |
 *     | (___   ___  _ __ ___   __\ \  /\  / /__  _ __ __| |___
 *      \___ \ / _ \| '_ ` _ \ / _ \ \/  \/ / _ \| '__/ _` / __|
 *      ____) | (_) | | | | | |  __/\  /\  / (_) | | | (_| \__ \
 *     |_____/ \___/|_| |_| |_|\___| \/  \/ \___/|_|  \__,_|___/ */
func SomeWords(message irc.ChatMessage) (retVal bool) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error Reading Some Words:", r)
		}
	}()

	retVal = false
	for i := range TwitchMsg.Words {
		textToSee := strings.ToUpper(message.Text)

		if strings.Contains(textToSee, "JOKE") { // Capute joke cue
			Send(RandomJokes())
			retVal = true
			break
		}

		if strings.Contains(textToSee, TwitchMsg.Words[i].Msg) {
			textToSend := TwitchMsg.Words[i].Reply[lib.RandomRange(0, len(TwitchMsg.Words[i].Reply))]
			Send(Substitution(textToSend, message))
			retVal = true
			break
		}

	}
	return retVal
}

/********_____*******_*********_***_*_*********_***_******************
 *      / ____|     | |       | | (_) |       | | (_)
 *     | (___  _   _| |__  ___| |_ _| |_ _   _| |_ _  ___  _ __
 *      \___ \| | | | '_ \/ __| __| | __| | | | __| |/ _ \| '_ \
 *      ____) | |_| | |_) \__ \ |_| | |_| |_| | |_| | (_) | | | |
 *     |_____/ \__,_|_.__/|___/\__|_|\__|\__,_|\__|_|\___/|_| |_|  */
func Substitution(textToSend string, msg irc.ChatMessage) (textToReturn string) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error Replacing Tags:", r)
		}
	}()

	textToReturn = textToSend
	textToReturn = strings.Replace(textToReturn, "{name}", msg.Sender.DisplayName, -1)
	textToReturn = strings.Replace(textToReturn, "{channel}", msg.Sender.DisplayName, -1)
	return textToReturn
}
