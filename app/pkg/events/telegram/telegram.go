package telegram

import (
	"errors"
	"fmt"
	"main/pkg/clients/telegram"
	"main/pkg/events"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type Processor struct {
	tg     *telegram.Client
	offset int
}

type Meta struct {
	ChatID   int
	Username string
}

func New(client *telegram.Client) *Processor {
	return &Processor{
		tg: client,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(e events.Event) error {
	switch e.Type {
	case events.Message:
		return p.processMessage(e)
	default:
		return ErrUnknownEventType
	}
}

func (p *Processor) processMessage(e events.Event) error {
	meta, err := meta(e)
	if err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}

	if err = p.doCmd(e.Text, meta.ChatID, meta.Username); err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}
	return nil
}

func meta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, ErrUnknownMetaType
	}
	return res, nil
}

func event(u telegram.Update) events.Event {
	uType := fetchType(u)

	res := events.Event{
		Type: uType,
		Text: fetchText(u),
	}

	if uType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			Username: u.Message.From.Username,
		}
	}

	return res
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}
	return u.Message.Text
}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}
	return events.Message
}
