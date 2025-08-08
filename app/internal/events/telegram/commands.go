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
	openingGateEntryModeCmd     = "opening_mode"
	openingGateEntryModeStopCmd = "opening_mode_stop"

	openGateEntryAction            = "⬅️Въезд⬅️️"
	openingGateEntryModeAction     = "⚠️Въезд на 5 минут⚠️"
	openingGateEntryModeStopAction = "✅️Закрыть✅"
)

var (
	actionsKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(openGateEntryAction),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(openingGateEntryModeAction),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(openingGateEntryModeStopAction),
		),
	)
)

type CommandsHandler struct {
	bot *tgbotapi.BotAPI
	gc  *gateController.GateController
	ch  chan int64
}

func NewCommandsHandler(bot *tgbotapi.BotAPI, ch chan int64) *CommandsHandler {
	return &CommandsHandler{
		bot: bot,
		gc:  gateController.NewController(),
		ch:  ch,
	}
}

func (h *CommandsHandler) Handle(users map[int64]User, update tgbotapi.Update) {
	ctx := context.Background()

	switch h.action(update) {
	case startCmd:
		h.sendStart(update)
	case helpCmd:
		h.sendHelp(update)
	case openGateEntryCmd, openGateEntryAction:
		h.sendOpenGateEntry(ctx, update)
	case openingGateEntryModeCmd, openingGateEntryModeAction:
		h.sendOpeningGateMode(ctx, users, update)
	case openingGateEntryModeStopCmd, openingGateEntryModeStopAction:
		h.sendOpeningGateModeStop(users, update)
	default:
		h.SendMsg(MsgUnknownCommand, update)
	}
}

func (h *CommandsHandler) SendMsg(messageText string, update tgbotapi.Update) {
	h.SendMsgWithChatId(messageText, update.Message.Chat.ID)
}

func (h *CommandsHandler) SendMsgWithChatId(messageText string, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, messageText)
	msg.ReplyMarkup = actionsKeyboard
	_, err := h.bot.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
		return
	}
}

func (h *CommandsHandler) action(update tgbotapi.Update) string {
	var action string

	if update.Message.IsCommand() {
		action = update.Message.Command()
	} else {
		action = update.Message.Text
	}

	return action
}

func (h *CommandsHandler) sendStart(update tgbotapi.Update) {
	h.SendMsg(msgStart, update)
}

func (h *CommandsHandler) sendHelp(update tgbotapi.Update) {
	h.SendMsg(msgHelp, update)
}

func (h *CommandsHandler) openGate(ctx context.Context, gateId string, update tgbotapi.Update) {
	if err := h.gc.OpenGate(ctx, gateId); err != nil {
		log.Println("cant open gate:", err)
		msg := fmt.Sprintf("%s: %s", msgCantGateOpen, err)
		h.SendMsg(msg, update)
		return
	}

	h.SendMsg(msgGateOpened, update)
}

func (h *CommandsHandler) sendOpenGateEntry(ctx context.Context, update tgbotapi.Update) {
	h.openGate(ctx, gateController.EntryGateId, update)
}

func (h *CommandsHandler) sendOpenGateExit(ctx context.Context, update tgbotapi.Update) {
	h.openGate(ctx, gateController.ExitGateId, update)
}

func (h *CommandsHandler) stopGateOpening(users map[int64]User, chatId int64) {
	if u, ok := users[chatId]; ok {
		u.cancelGateMode()
	}
}

func (h *CommandsHandler) startGatesOpening(
	ctx context.Context,
	users map[int64]User,
	update tgbotapi.Update,
	gates []string,
	ticker *time.Ticker,
	msg string,
) {
	chatId := update.Message.Chat.ID
	h.stopGateOpening(users, chatId)

	chErr := make(chan error)

	ctxWithTimeout, cancel := context.WithCancel(ctx)
	users[chatId] = NewUser(cancel)

	h.gc.OpenGateForTimePeriod(ctxWithTimeout, chErr, ticker, gates)
	h.SendMsg(msg, update)

	// start checking errors
	go func() {
		for {
			select {
			case <-ctxWithTimeout.Done():
				return
			case <-ticker.C:
				h.ch <- chatId
				return
			case err := <-chErr:
				h.SendMsg(fmt.Sprintf("%s: %v", msgCantGateOpen, err), update)
			}
		}
	}()
}

func (h *CommandsHandler) sendOpeningGateMode(ctx context.Context, users map[int64]User, update tgbotapi.Update) {
	ticker := time.NewTicker(5 * time.Minute)
	gates := []string{gateController.EntryGateId}
	h.startGatesOpening(ctx, users, update, gates, ticker, msgGateOpeningModeActivated)
}

func (h *CommandsHandler) sendOpeningGateModeStop(users map[int64]User, update tgbotapi.Update) {
	h.stopGateOpening(users, update.Message.Chat.ID)
	h.SendMsg(msgGateOpeningModeDeactivated, update)
}
