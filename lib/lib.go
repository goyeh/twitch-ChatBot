package lib

/* places as a simple wrapper library of self contained functions, no dependancies are required, and these functions will have default values that contain overwride funtions. */

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//const (
//	timeFormat = "2006-01-02 15:04:05"
//)

var (
	debugLevel string
)

/***    _       _ _
 *     (_)     (_) |
 *      _ _ __  _| |_
 *     | | '_ \| | __|
 *     | | | | | | |_
 *     |_|_| |_|_|\__|  * * */
func init() { // Set defaults,
	var err error
	debugLevel = "DEBUG INFO NOTICE WARN ERROR CRIT STDOUT"
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.Println("Log Init:", err)
}

/****   _             _    _      _
 *     | |           | |  | |    | |
 *     | | ___   __ _| |__| | ___| |_ __   ___ _ __ ___
 *     | |/ _ \ / _` |  __  |/ _ \ | '_ \ / _ \ '__/ __|
 *     | | (_) | (_| | |  | |  __/ | |_) |  __/ |  \__ \
 *     |_|\___/ \__, |_|  |_|\___|_| .__/ \___|_|  |___/
 *               __/ |             | |
 *              |___/              |_|                        */
func Debug(msg ...interface{})  { logCore("DEBUG", msg...) }
func Info(msg ...interface{})   { logCore("INFO", msg...) }
func Notice(msg ...interface{}) { logCore("NOTICE", msg...) }
func Warn(msg ...interface{})   { logCore("WARN", msg...) }
func Error(msg ...interface{})  { logCore("ERROR", msg...) }
func Crit(msg ...interface{})   { logCore("CRIT", msg...) }

/* ------------------------------------- */

/***    _              _____
 *     | |            / ____|
 *     | | ___   __ _| |     ___  _ __ ___
 *     | |/ _ \ / _` | |    / _ \| '__/ _ \
 *     | | (_) | (_| | |___| (_) | | |  __/
 *     |_|\___/ \__, |\_____\___/|_|  \___|
 *               __/ |
 *              |___/                           */
func logCore(level string, msg ...interface{}) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Print("Error detected logging:", r)
		}
	}()
	if strings.Contains(debugLevel, level) {
		switch level {
		case "DEBUG":
			log.Debug(fmt.Sprint(msg...))
		case "INFO":
			log.Info(fmt.Sprint(msg...))
		case "WARN":
			log.Warning(fmt.Sprint(msg...))
		case "ERROR":
			log.Error(fmt.Sprint(msg...))
		case "CRIT":
			log.Fatal(fmt.Sprint(msg...))
		}

		if strings.Contains(debugLevel, "STDOUT") {
			fmt.Println(msg...)
		}
	}
}

// ----------------------------------------------------------------------------
//
//	    ____                  __                   ____
//	   / __ \____ _____  ____/ /___  ____ ___     / __ \____ _____  ____ ____
//	  / /_/ / __ `/ __ \/ __  / __ \/ __ `__ \   / /_/ / __ `/ __ \/ __ `/ _ \
//	 / _, _/ /_/ / / / / /_/ / /_/ / / / / / /  / _, _/ /_/ / / / / /_/ /  __/
//	/_/ |_|\__,_/_/ /_/\__,_/\____/_/ /_/ /_/  /_/ |_|\__,_/_/ /_/\__, /\___/
//	                                                             /____/
func RandomRange(min int, max int) (randomRange int) {
	rand.Seed(time.Now().UnixNano())
	randomRange = rand.Intn(max-min) + min
	return
}

// -------------------------------------------------
//
//	  ___ _           _     ___
//	 / __| |_  ___ __| |__ | __|_ _ _ _ ___ _ _
//	| (__| ' \/ -_) _| / / | _|| '_| '_/ _ \ '_|
//	 \___|_||_\___\__|_\_\ |___|_| |_| \___/_|
func CheckErr(err error) (isErr bool) {
	defer func() {
		r := recover()
		if r != nil {
			Error("Error detected:", r)
		}
	}()
	isErr = false
	if err != nil {
		isErr = true
		a, b, c, _ := runtime.Caller(1)
		Error(err, " in ", "Process ID:", a, "In Module:", b, "Line:", c) //Return Error object
	}

	return
}

// -------------------------------------------------
//
//	  ___ _           _  DB ___
//	 / __| |_  ___ __| |__ | __|_ _ _ _ ___ _ _
//	| (__| ' \/ -_) _| / / | _|| '_| '_/ _ \ '_|
//	 \___|_||_\___\__|_\_\ |___|_| |_| \___/_|
func CheckDbErr(err error, db *sql.DB, msg ...interface{}) (isErr bool) {
	defer func() {
		r := recover()
		if r != nil {
			Error("DB Error detected:", r)
		}
	}()
	isErr = false
	if err != nil {
		isErr = true
		a, b, c, _ := runtime.Caller(1)
		Error("DB", db.Ping(), err, "Process ID:", a, "In Module:", b, "Line:", c, msg) //Return the Database Error object
	}
	return
}

//	      _
//	     | |
//	  ___| | ___  ___  ___
//	 / __| |/ _ \/ __|/ _ \
//	| (__| | (_) \__ \  __/
//	 \___|_|\___/|___/\___|
//
// Safe close routine
func DeferClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		Crit(err)
	}
}

/*************************************************
 *            _ _  _____ _        _
 *           (_) |/ ____| |      (_)
 *      _ __  _| | (___ | |_ _ __ _ _ __   __ _
 *     | '_ \| | |\___ \| __| '__| | '_ \ / _` |
 *     | | | | | |____) | |_| |  | | | | | (_| |
 *     |_| |_|_|_|_____/ \__|_|  |_|_| |_|\__, |
 *                                         __/ |
 *                                        |___/
 * * * * * * * * * * * * * * * * * * * * * * * * *
 * prevents a nil string from infecting data
 * --------------------------------------------- */
func NilString(s string) string {
	if len(s) == 0 {
		return ""
	}
	return s
}

/*****************************************************
 *      _        _           _
 *     | |      (_)         | |
 *     | |_ _ __ _ _ __ ___ | |     ___ _ __
 *     | __| '__| | '_ ` _ \| |    / _ \ '_ \
 *     | |_| |  | | | | | | | |___|  __/ | | |
 *      \__|_|  |_|_| |_| |_|______\___|_| |_|
 * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * Trim the length of a string based on the max length
 * -------------------------------------------------- */
func TrimLen(str string, size int) (splited string) {
	strLength := len(str)
	splitedLength := int(math.Ceil(float64(strLength) / float64(size)))
	var start, stop int
	for i := 0; i < splitedLength; i += 1 {
		start = i * size
		stop = start + size
		if stop > strLength {
			stop = strLength
		}
		splited = str[start:stop]
	}
	return splited
}

/****************************************************************************************
 *                                _ _______    _____                        _
 *                               | |__   __|  |  __ \                      (_)
 *      _ __ ___ _ __   ___  _ __| |_ | | ___ | |  | | ___  _ __ ___   __ _ _ _ __
 *     | '__/ _ \ '_ \ / _ \| '__| __|| |/ _ \| |  | |/ _ \| '_ ` _ \ / _` | | '_ \
 *     | | |  __/ |_) | (_) | |  | |_ | | (_) | |__| | (_) | | | | | | (_| | | | | |
 *     |_|  \___| .__/ \___/|_|   \__||_|\___/|_____/ \___/|_| |_| |_|\__,_|_|_| |_|
 *              | |
 *              |_|
 * --------------------------------------------------------------------------------------
 * Function to report to domains, this will require paramters, since we can not predict
 * how each application is managing its configuration control.
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */
type Services struct {
	Service  string `json:"service"`
	Dns      string `json:"dns"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Cb       string `json:"callback"`
	Expected string `json:"expected"`
	Deps     string `json:"deps"`
	Meta     string `json:"meta"`
}

func ReportToDomains(option string, message string, myServiceName string, myDns string, myIp string, myPort string, url string) {
	defer func() {
		r := recover()
		if r != nil {
			Error("Unable to reach Domains:", r, "You may need to add http:// to the domains url")
		}
	}()
	switch option {
	case "start":
		todo := Services{myServiceName, myDns, myIp, myPort, "", "", "", "News Service"}
		jsonReq, _ := json.Marshal(todo)
		Debug("Starting monitor at:", url+"/start")
		resp, err := http.Post(url+"/start", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
		CheckErr(err)
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		Info("INFO", "Reported to Domains Result:", string(bodyBytes))
	case "ping":
		resp, err := http.Get(url + "/ping/" + myServiceName)
		CheckErr(err)
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		Info("Ping to Domains Result:", string(bodyBytes))
	case "fail":
		resp, err := http.Get(url + "/fail/" + myServiceName)
		CheckErr(err)
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		Error("Fail to Domains Result:", string(bodyBytes))
	case "error": // special handling since we want to pass the error message to domains
		resp, err := http.Get(url + "/fail/" + myServiceName + "/" + message)
		CheckErr(err)
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		Error("Fail to Domains Result:", string(bodyBytes))
	}
}

/******************************************************************************
 *      _   _               _   _   _               _ _
 *     | | | |             | | | | | |        /\   | | |
 *     | |_| |__  _ __ ___ | |_| |_| | ___   /  \  | | | _____      __
 *     | __| '_ \| '__/ _ \| __| __| |/ _ \ / /\ \ | | |/ _ \ \ /\ / /
 *     | |_| | | | | | (_) | |_| |_| |  __// ____ \| | | (_) \ V  V /
 *      \__|_| |_|_|  \___/ \__|\__|_|\___/_/    \_\_|_|\___/ \_/\_/
 * ----------------------------------------------------------------------------
 * This function will temp store the value in a map and then remove it, it will
 * return true or false if the item is in the map, Now sets delay on second response
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */
var throttle = make(map[string]bool)

func ThrottleAllow(ip string, timeout int) (retVal bool) {
	if throttle[ip] {
		Warn("-=Throttle=-To frequent calls from:", ip)
		time.Sleep(time.Duration(timeout) * time.Second) //Random next cycle.
		retVal = true                                    // false will result is receiging to frequent message
	} else {
		throttle[ip] = true
		go func() {
			time.Sleep(time.Duration(timeout) * time.Nanosecond) //Random next cycle.
			delete(throttle, ip)
		}()
		retVal = true
	}
	return
}

func StayAlive() {
	Debug("Stay Alive Started")
	for {
		time.Sleep(time.Duration(60) * time.Second)
		//status = Status{1, "OK"}
	}
}

func Status() string {

	return "OK"
}
