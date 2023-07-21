package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"main/pkg/events/telegram"
	"os"
	"strings"
)

func mustCheckEnvVars() bool {
	for _, envVar := range []string{
		"BOT_TOKEN",
		"SID",
		"ADMINS",
	} {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Error: env var \"%s\" empty\n", envVar)
		}
	}

	debugMode := false
	if debug := os.Getenv("DEBUG"); debug == "1" || debug == "true" {
		debugMode = true
	}
	return debugMode
}

func initAdmins() map[string]struct{} {
	admins := make(map[string]struct{})
	administrators := strings.Split(os.Getenv("ADMINS"), ", ")
	for _, a := range administrators {
		admins[a] = struct{}{}
	}
	return admins
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

	admins := initAdmins()
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("User @%s sended message: %s\n", update.Message.Chat.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			commandHandler := telegram.NewCommandsHandler(bot, update)
			if _, ok := admins[update.Message.Chat.UserName]; !ok {
				commandHandler.SendMsg(telegram.MsgNotAllowedControl)
				continue
			}
			commandHandler.Handle(openingGateMode)
		}
	}
}
