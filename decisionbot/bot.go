package decisionbot

import (
	"errors"
	"log"
	"math/rand/v2"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DecisionBot struct {
	API             *tgbotapi.BotAPI
	Updates         tgbotapi.UpdatesChannel
	ActiveDecisions map[int64]*Decision
}

func NewDecisionBot(token string) (*DecisionBot, error) {
	b := &DecisionBot{}
	if token == "" {
		return nil, errors.New("no telegram token found")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	b.API = bot
	b.API.Debug = true
	log.Printf("Authorized on account %s\n", b.API.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	b.Updates = bot.GetUpdatesChan(u)
	b.ActiveDecisions = make(map[int64]*Decision)
	log.Println("Opened updates channel successfully")
	return b, nil
}

func (b *DecisionBot) SendTextMessage(chatId int64, msg string) {
	b.API.Send(tgbotapi.NewMessage(chatId, msg))
}

func (b *DecisionBot) SendCommandKeyboard(chatId int64, msg string) {
	var commandKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/decide"),
			tgbotapi.NewKeyboardButton("/yesno"),
		),
	)
	reply := tgbotapi.NewMessage(chatId, msg)
	reply.ReplyMarkup = commandKeyboard
	b.API.Send(reply)
}

func (b *DecisionBot) RemoveCommandKeyboard(chatId int64, msg string) {
	reply := tgbotapi.NewMessage(chatId, msg)
	reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	b.API.Send(reply)
}

func (b *DecisionBot) HasActiveDecision(chatId int64) bool {
	_, ok := b.ActiveDecisions[chatId]
	return ok
}

func (b *DecisionBot) DecideYesOrNo(chatId int64) {
	decision := rand.IntN(2)
	if decision == 1 {
		b.SendTextMessage(chatId, "Yes")
	} else {
		b.SendTextMessage(chatId, "No")
	}
	log.Printf("Decide yes/no for %d -> %v", chatId, decision != 0)
}
