package telegram

import (
	"context"
	"log"
	"strings"
)

const (
	StartCmd                = "/start"
	HelpCmd                 = "/help"
	OpenGateEntryCmd        = "/openGateEntry"
	OpenGateExitCmd         = "/openGateExit"
	OpeningGateEntryModeCmd = "/openGateEntryMode"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from %s\n", text, username)

	switch text {
	case StartCmd:
		return p.sendHello(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	case OpenGateEntryCmd:
		return p.sendOpenGateEntry(chatID)
	case OpenGateExitCmd:
		return p.sendOpenGateExit(chatID)
	case OpeningGateEntryModeCmd:
		return p.sendOpeningGateModeCmd(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendOpenGateEntry(chatID int) error {
	if err := p.gc.OpenGate(true); err != nil {
		log.Printf("cant open gate: %s\n", err)
		return p.tg.SendMessage(chatID, msgCantGateOpen)
	}
	return p.tg.SendMessage(chatID, msgGateOpened)
}

func (p *Processor) sendOpenGateExit(chatID int) error {
	if err := p.gc.OpenGate(false); err != nil {
		log.Printf("cant open gate: %s\n", err)
		return p.tg.SendMessage(chatID, msgCantGateOpen)
	}
	return p.tg.SendMessage(chatID, msgGateOpened)
}

func (p *Processor) sendOpeningGateModeCmd(chatID int) error {
	p.gc.OpenGateAlways(context.Background())
	return p.tg.SendMessage(chatID, msgGateOpeningModeActivated)
}
