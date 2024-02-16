package main

import (
	"flag"
	"log"
	"os"
	"read-adviser-bot/clients/telegram"
)

func main() {
	os.Setenv("TG_BOT_HOST", "api.telegram.org")
	os.Setenv("GET_UPDATES_METHOD", "getUpdates")
	os.Setenv("SEND_MESSAGE_METHOD", "sendMessage")
	os.Setenv("DEFAULT_PERMISSION", "0774")

	tgClient := telegram.New(os.Getenv("TG_BOT_HOST"), mustToken())
}

func mustToken() string {
	token := flag.String("token-bot-token", "", "token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatalln("token is not specified")
	}

	return *token
}
