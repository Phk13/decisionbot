package decisionbot

import (
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DecisionBot struct {
	API             *tgbotapi.BotAPI
	Updates         tgbotapi.UpdatesChannel
	ActiveDecisions map[int64]*Decision
}

func (b *DecisionBot) NewDecisionBot(token string) error {
	if token == "" {
		return errors.New("no telegram token found")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	b.API = bot
	b.API.Debug = true
	log.Printf("Authorized on account %s\n", b.API.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	b.Updates = bot.GetUpdatesChan(u)
	b.ActiveDecisions = make(map[int64]*Decision)
	log.Println("Opened updates channel successfully")
	return nil
}

func (b *DecisionBot) SendTextMessage(chatId int64, msg string) {
	b.API.Send(tgbotapi.NewMessage(chatId, msg))
}

func (b *DecisionBot) HasActiveDecision(chatId int64) bool {
	_, ok := b.ActiveDecisions[chatId]
	return ok
}
