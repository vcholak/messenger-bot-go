package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/vcholak/messenger-bot/lib/errp"
	"github.com/vcholak/messenger-bot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *EventProcessor) doCmd(ctx context.Context, text string, chatID int, firstName string) error {
	text = strings.TrimSpace(text)

	log.Printf("Got a new command '%s' from user '%s'", text, firstName)

	if isAddCmd(text) {
		return p.savePage(ctx, chatID, text, firstName)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(ctx, chatID, firstName)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *EventProcessor) savePage(ctx context.Context, chatID int, pageURL string, firstName string) (err error) {
	defer func() { err = errp.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:       pageURL,
		FirstName: firstName,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(ctx, chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(ctx, chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *EventProcessor) sendRandom(ctx context.Context, chatID int, firstName string) (err error) {
	defer func() { err = errp.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(ctx, firstName)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(ctx, chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(ctx, chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(ctx, page)
}

func (p *EventProcessor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *EventProcessor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
