package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"twitchbot/lib"

	"github.com/joho/godotenv"
)

type Config struct {
	configFilePath string
	DEBUG          int
	VERSION        string
	AppName        string
	LOGDIR         string
	Username       string
	Oauth          string
	Channel        string
}

var Val Config

func init() {
	defer func() {
		r := recover()
		if r != nil {
			log.Print("Possible .env error:", r)
		}
	}()
	flag.StringVar(&Val.configFilePath, "config", "conf.env", "config file path")
	flag.Parse()
	lib.CheckErr(godotenv.Load(Val.configFilePath))
	Val = Config{Val.configFilePath,
		getEnvAsInt("DEBUG", 2),
		getEnv("VERSION", "0.0.1"),
		filepath.Base(os.Args[0]),
		getEnv("LOGDIR", "."),
		getEnv("USERNAME", "botmasterhk"),
		getEnv("OOAUTH", "oauth:rlrasvkbmh57ct9luaesriom1ob5gs"),
		getEnv("CHANNEL", "botmasterhk"),
	}
	lib.Info("Logfile:", Val.AppName)
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		retVal := time.Duration(value)
		return retVal
	}
	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt64(name string, defaultVal int64) int64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}

	return defaultVal
}
