package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	gateController "main/pkg/gate-controller"
)

const (
	startCmd                    = "start"
	helpCmd                     = "help"
	openGateEntryCmd            = "open_entry"
	openGateExitCmd             = "open_exit"
	openingGateEntryModeCmd     = "opening_mode"
	openingGateEntryModeStopCmd = "opening_mode_stop"

	openGateEntryAction            = "⬅️Въезд⬅️️"
	openGateExitAction             = "➡️Выезд➡️"
	openingGateEntryModeAction     = "⚠️5 минут⚠️"
	openingGateEntryModeStopAction = "✅️Закрыть✅"
)

var (
	actionsKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(openGateEntryAction),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(openGateExitAction),
			tgbotapi.NewKeyboardButton(openingGateEntryModeAction),
			tgbotapi.NewKeyboardButton(openingGateEntryModeStopAction),
		),
	)
)

type CommandsHandler struct {
	bot *tgbotapi.BotAPI
	gc  *gateController.GateController
}

func NewCommandsHandler(bot *tgbotapi.BotAPI) *CommandsHandler {
	return &CommandsHandler{
		bot: bot,
		gc:  gateController.NewController(),
	}
}

func (ch *CommandsHandler) Handle(users map[int64]User, update tgbotapi.Update) {
	ctx := context.Background()

	switch ch.action(update) {
	case startCmd:
		ch.sendHello(update)
	case helpCmd:
		ch.sendHelp(update)
	case openGateEntryCmd, openGateEntryAction:
		ch.sendOpenGateEntry(ctx, update)
	case openGateExitCmd, openGateExitAction:
		ch.sendOpenGateExit(ctx, update)
	case openingGateEntryModeCmd, openingGateEntryModeAction:
		ch.sendOpeningGateModeCmd(ctx, users, update)
	case openingGateEntryModeStopCmd, openingGateEntryModeStopAction:
		ch.sendOpeningGateModeStopCmd(users, update)
	default:
		ch.SendMsg(MsgUnknownCommand, update)
	}
}

func (ch *CommandsHandler) SendMsg(messageText string, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	msg.ReplyMarkup = actionsKeyboard
	_, err := ch.bot.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}

func (ch *CommandsHandler) action(update tgbotapi.Update) string {
	var action string

	if update.Message.IsCommand() {
		action = update.Message.Command()
	} else {
		action = update.Message.Text
	}

	return action
}

func (ch *CommandsHandler) sendHello(update tgbotapi.Update) {
	ch.SendMsg(msgHello, update)
}

func (ch *CommandsHandler) sendHelp(update tgbotapi.Update) {
	ch.SendMsg(msgHelp, update)
}

func (ch *CommandsHandler) openGate(ctx context.Context, gateId string, update tgbotapi.Update) {
	if err := ch.gc.OpenGate(ctx, gateId); err != nil {
		log.Println("cant open gate:", err)
		msg := fmt.Sprintf("%s: %s", msgCantGateOpen, err)
		ch.SendMsg(msg, update)
		return
	}

	ch.SendMsg(msgGateOpened, update)
}

func (ch *CommandsHandler) sendOpenGateEntry(ctx context.Context, update tgbotapi.Update) {
	ch.openGate(ctx, gateController.EntryGateId, update)
}

func (ch *CommandsHandler) sendOpenGateExit(ctx context.Context, update tgbotapi.Update) {
	ch.openGate(ctx, gateController.ExitGateId, update)
}

func (ch *CommandsHandler) stopGateOpening(users map[int64]User, chatId int64) {
	if u, ok := users[chatId]; ok {
		u.cancelGateMode()
	}
}

func (ch *CommandsHandler) sendOpeningGateModeCmd(ctx context.Context, users map[int64]User, update tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	ch.stopGateOpening(users, chatId)

	openingGateDuration := 5 * time.Minute
	chErr := make(chan error)

	ctxWithTimeout, cancel := context.WithCancel(ctx)
	users[chatId] = NewUser(cancel)

	ch.gc.OpenGateForTimePeriod(ctxWithTimeout, chErr, openingGateDuration)
	ch.SendMsg(msgGateOpeningModeActivated, update)

	// start checking errors
	go func() {
		for {
			select {
			case <-ctxWithTimeout.Done():
				return
			case err := <-chErr:
				ch.SendMsg(fmt.Sprintf("%s: %v", msgCantGateOpen, err), update)
			}
		}
	}()
}

func (ch *CommandsHandler) sendOpeningGateModeStopCmd(users map[int64]User, update tgbotapi.Update) {
	ch.stopGateOpening(users, update.Message.Chat.ID)
	ch.SendMsg(msgGateOpeningModeDeactivated, update)
}
