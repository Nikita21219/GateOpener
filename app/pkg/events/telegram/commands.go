package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	gate_controller "main/pkg/gate-controller"
	"time"
)

const (
	StartCmd                    = "start"
	HelpCmd                     = "help"
	OpenGateEntryCmd            = "open_entry"
	OpenGateExitCmd             = "open_exit"
	OpeningGateEntryModeCmd     = "opening_mode"
	OpeningGateEntryModeStopCmd = "opening_mode_stop"

	OpenGateEntryAction            = "⬅️Въезд⬅️️"
	OpenGateExitAction             = "➡️Выезд➡️"
	OpeningGateEntryModeAction     = "⚠️5 минут⚠️"
	OpeningGateEntryModeStopAction = "✅️Закрыть✅"

	urlAMVideoApi = "https://lk.amvideo-msk.ru/api/api4.php"
)

var actionsKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(OpenGateEntryAction),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(OpenGateExitAction),
		tgbotapi.NewKeyboardButton(OpeningGateEntryModeAction),
		tgbotapi.NewKeyboardButton(OpeningGateEntryModeStopAction),
	),
)

type CommandsHandler struct {
	bot    *tgbotapi.BotAPI
	update tgbotapi.Update
	gc     *gate_controller.GateController
}

func NewCommandsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) *CommandsHandler {
	return &CommandsHandler{
		bot:    bot,
		update: update,
		gc:     gate_controller.NewController(urlAMVideoApi),
	}
}

func (ch *CommandsHandler) Handle(users map[int64]User) {
	if u, ok := users[ch.update.Message.Chat.ID]; ok {
		u.cancel()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	users[ch.update.Message.Chat.ID] = User{
		ctx:    ctx,
		cancel: cancel,
	}

	var action string
	if ch.update.Message.IsCommand() {
		action = ch.update.Message.Command()
	} else {
		action = ch.update.Message.Text
	}

	switch action {
	case StartCmd:
		ch.sendHello()
	case HelpCmd:
		ch.sendHelp()
	case OpenGateEntryCmd, OpenGateEntryAction:
		ch.sendOpenGateEntry()
	case OpenGateExitCmd, OpenGateExitAction:
		ch.sendOpenGateExit()
	case OpeningGateEntryModeCmd, OpeningGateEntryModeAction:
		ch.sendOpeningGateModeCmd(ctx)
	case OpeningGateEntryModeStopCmd, OpeningGateEntryModeStopAction:
		users[ch.update.Message.Chat.ID].cancel()
		ch.sendOpeningGateModeStopCmd()
	default:
		ch.SendMsg(MsgUnknownCommand)
	}
}

func (ch *CommandsHandler) SendMsg(messageText string) {
	msg := tgbotapi.NewMessage(ch.update.Message.Chat.ID, messageText)
	msg.ReplyMarkup = actionsKeyboard
	_, err := ch.bot.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}

func (ch *CommandsHandler) sendHello() {
	ch.SendMsg(msgHello)
}

func (ch *CommandsHandler) sendHelp() {
	ch.SendMsg(msgHelp)
}

func (ch *CommandsHandler) sendOpenGateEntry() {
	if err := ch.gc.OpenGate(true); err != nil {
		log.Printf("cant open gate: %s\n", err)
		msg := fmt.Sprintf("%s: %s", msgCantGateOpen, err)
		ch.SendMsg(msg)
		return
	}
	ch.SendMsg(msgGateOpened)
}

func (ch *CommandsHandler) sendOpenGateExit() {
	if err := ch.gc.OpenGate(false); err != nil {
		log.Printf("cant open gate: %s\n", err)
		msg := fmt.Sprintf("%s: %s", msgCantGateOpen, err)
		ch.SendMsg(msg)
		return
	}
	ch.SendMsg(msgGateOpened)
}

func (ch *CommandsHandler) sendOpeningGateModeCmd(ctx context.Context) {
	ch.gc.OpenGateAlways(ctx)
	ch.SendMsg(msgGateOpeningModeActivated)
}

func (ch *CommandsHandler) sendOpeningGateModeStopCmd() {
	ch.SendMsg(msgGateOpeningModeDeactivated)
}
