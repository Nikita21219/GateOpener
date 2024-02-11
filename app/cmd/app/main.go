package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"main/pkg/events/telegram"
	"main/pkg/utils"
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

	return utils.Debug()
}

func initAdmins() map[string]struct{} {
	admins := make(map[string]struct{})
	administrators := strings.Split(os.Getenv("ADMINS"), ", ")
	for _, a := range administrators {
		admins[a] = struct{}{}
	}
	return admins
}

func logUserMessage(update tgbotapi.Update) {
	var userName string
	if update.Message.Chat.UserName != "" {
		userName = fmt.Sprintf("@%s", update.Message.Chat.UserName)
	} else {
		userName = fmt.Sprintf("%s %s", update.Message.Chat.FirstName, update.Message.Chat.LastName)
	}

	log.Printf("User %s sended message: %s\n", userName, update.Message.Text)
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
	usersContexts := make(map[int64]telegram.User)

	admins := initAdmins()
	for update := range updates {
		if update.Message == nil {
			continue
		}

		logUserMessage(update)

		commandHandler := telegram.NewCommandsHandler(bot, update)
		if _, ok := admins[update.Message.Chat.UserName]; !ok {
			commandHandler.SendMsg(telegram.MsgNotAllowedControl)
			continue
		}
		commandHandler.Handle(usersContexts)
	}
}
