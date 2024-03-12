package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"main/internal/events/telegram"
	"main/pkg/utils"
)

var (
	admins = initAdmins()
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
	m := make(map[string]struct{})
	administrators := strings.Split(os.Getenv("ADMINS"), ", ")
	for _, a := range administrators {
		m[a] = struct{}{}
	}

	return m
}

func logUserMessage(update tgbotapi.Update) {
	var userName string
	if update.Message.Chat.UserName != "" {
		userName = fmt.Sprintf("@%s", update.Message.Chat.UserName)
	} else {
		userName = fmt.Sprintf("%s %s", update.Message.Chat.FirstName, update.Message.Chat.LastName)
	}

	log.Println(fmt.Sprintf("User %s sended message: %s", userName, update.Message.Text))
}

func userNotAdmin(update tgbotapi.Update) bool {
	_, ok := admins[update.Message.Chat.UserName]
	return !ok
}

func main() {
	debug := mustCheckEnvVars()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalln("error to create bot:", err)
	}
	bot.Debug = debug

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalln("error to get updates", err)
	}

	usersContexts := make(map[int64]telegram.User)
	ch := make(chan int64)
	commandsHandler := telegram.NewCommandsHandler(bot, ch)

	go func() {
		// if opening mode stopped - send message
		for {
			chatId, ok := <-ch
			if !ok {
				log.Println("the channel listening to the events is closed")
				return
			}
			commandsHandler.SendMsgWithChatId(telegram.MsgOpeningModeStopped, chatId)
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		logUserMessage(update)
		if userNotAdmin(update) {
			commandsHandler.SendMsg(telegram.MsgNotAllowedControl, update)
			continue
		}

		commandsHandler.Handle(usersContexts, update)
	}
}
