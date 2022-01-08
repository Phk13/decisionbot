package decisionbot

import (
	"errors"
	"log"
	"math/rand"
	"time"

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
	rand.Seed(time.Now().Unix())
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

func (b *DecisionBot) DecideYesOrNo(chatId int64) {
	decision := rand.Intn(2)
	if decision == 1 {
		b.SendTextMessage(chatId, "Yes")
	} else {
		b.SendTextMessage(chatId, "No")
	}
	log.Printf("Decide yes/no for %d -> %v", chatId, decision != 0)
}
