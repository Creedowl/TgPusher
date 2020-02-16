package utils

import (
	"log"
	"os"
)

var _config *Config = nil

type Db struct {
	Host     string
	Port     int
	Database string
}

type Config struct {
	DB       Db
	BotToken string
	Debug    bool
}

func GetConfig() *Config {
	if _config == nil {
		token, existed := os.LookupEnv("BOT_TOKEN")
		if !existed {
			log.Fatalln("please set BOT_TOKEN")
		}
		_debug := os.Getenv("DEBUG")
		debug := false
		if _debug == "true" {
			debug = true
		}
		_config = &Config{
			DB: Db{
				Host:     "localhost",
				Port:     27017,
				Database: "TgPusher",
			},
			BotToken: token,
			Debug:    debug,
		}
	}
	return _config
}

func GetHelp(lang string) string {
	switch lang {
	case "zh-hans":
		return `
I can help you send message via sending message to server

*Usage:*
/help \- show usage
/token \- generate new token or show your token
/revoke \- revoke your token
`
	default:
		return "asdf"
	}
}
