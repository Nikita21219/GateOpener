package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"main/pkg/events/telegram"
	"os"
)

func mustCheckEnvVars() bool {
	if os.Getenv("BOT_TOKEN") == "" || os.Getenv("SID") == "" {
		log.Fatalln("Bot token or SID empty")
	}

	debugMode := false
	if debug := os.Getenv("DEBUG"); debug == "1" || debug == "true" {
		debugMode = true
	}
	return debugMode
}

func main() {
	debug := mustCheckEnvVars()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = debug

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates, err := bot.GetUpdatesChan(u)
	openingGateMode := make(chan bool)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			commandHandler := telegram.NewCommandsHandler(bot, update)
			commandHandler.Handle(openingGateMode)
		}
	}
}
