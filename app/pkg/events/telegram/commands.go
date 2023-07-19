package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	gate_controller "main/pkg/gate-controller"
)

const (
	StartCmd                    = "start"
	HelpCmd                     = "help"
	OpenGateEntryCmd            = "open_entry"
	OpenGateExitCmd             = "open_exit"
	OpeningGateEntryModeCmd     = "opening_mode"
	OpeningGateEntryModeStopCmd = "opening_mode_stop"

	urlAMVideoApi = "https://lk.amvideo-msk.ru/api/api4.php"
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

func (ch *CommandsHandler) Handle(openingGateMode chan bool) {
	switch ch.update.Message.Command() {
	case StartCmd:
		ch.sendHello()
	case HelpCmd:
		ch.sendHelp()
	case OpenGateEntryCmd:
		ch.sendOpenGateEntry()
	case OpenGateExitCmd:
		ch.sendOpenGateExit()
	case OpeningGateEntryModeCmd:
		ch.sendOpeningGateModeCmd(openingGateMode)
	case OpeningGateEntryModeStopCmd:
		ch.sendOpeningGateModeStopCmd(openingGateMode)
	default:
		ch.sendMsg(msgUnknownCommand)
	}
}

func (ch *CommandsHandler) sendMsg(messageText string) {
	msg := tgbotapi.NewMessage(ch.update.Message.Chat.ID, messageText)
	_, err := ch.bot.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}

func (ch *CommandsHandler) sendHello() {
	ch.sendMsg(msgHello)
}

func (ch *CommandsHandler) sendHelp() {
	ch.sendMsg(msgHelp)
}

func (ch *CommandsHandler) sendOpenGateEntry() {
	if err := ch.gc.OpenGate(true); err != nil {
		log.Printf("cant open gate: %s\n", err)
		ch.sendMsg(msgCantGateOpen)
	}
	ch.sendMsg(msgGateOpened)
}

func (ch *CommandsHandler) sendOpenGateExit() {
	if err := ch.gc.OpenGate(false); err != nil {
		log.Printf("cant open gate: %s\n", err)
		ch.sendMsg(msgCantGateOpen)
	}
	ch.sendMsg(msgGateOpened)
}

func (ch *CommandsHandler) sendOpeningGateModeCmd(openingGateMode chan bool) {
	ch.gc.OpenGateAlways(openingGateMode)
	ch.sendMsg(msgGateOpeningModeActivated)
}

func (ch *CommandsHandler) sendOpeningGateModeStopCmd(openingGateMode chan bool) {
	openingGateMode <- false
	ch.sendMsg(msgGateOpeningModeDeactivated)
}
