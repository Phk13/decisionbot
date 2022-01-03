package decisionbot

import (
	"fmt"
	"log"
	"sync"
)

/* StartDecision starts a new decision process for chatId */
func (bot *DecisionBot) StartDecision(chatId int64) {
	bot.ActiveDecisions[chatId] = &Decision{Lock: sync.Mutex{}}
	log.Printf("Starting decision for %d\n", chatId)
	bot.SendTextMessage(chatId, "OK, let me know the choices to decide.")
}

/* StopDecision executes a random selection on all input choices for chatId and finishes the active decision */
func (bot *DecisionBot) StopDecision(chatId int64) {
	decision := bot.ActiveDecisions[chatId]
	decision.Lock.Lock()

	choice := bot.ActiveDecisions[chatId].Decide()
	log.Printf("Decision for %d: %s\n", chatId, choice)

	bot.SendTextMessage(chatId, fmt.Sprintf("Alright, out of the %d choices, you should decide %s", decision.ChoiceNumber(), choice))
	log.Printf("Finishing decision for %d\n", chatId)

	decision.Lock.Unlock()
	delete(bot.ActiveDecisions, chatId)
}

/* AddChoice adds a choice to a currently active decision for chatId */
func (bot *DecisionBot) AddChoice(chatId int64, choice string) {
	if bot.HasActiveDecision(chatId) {
		decision := bot.ActiveDecisions[chatId]
		decision.Lock.Lock()
		decision.AddChoice(choice)
		decision.Lock.Unlock()
	}
}

/* ListenAndDecide starts listening on updates channel and responding to commands and messages. */
func (bot *DecisionBot) ListenAndDecide() {
	for update := range bot.Updates {
		// Ignore any non-message updates
		if update.Message == nil {
			continue
		}
		chatId := update.Message.Chat.ID
		// If message is not a command and the bot has an active decision for update's chat, add as a choice.
		if !update.Message.IsCommand() && bot.HasActiveDecision(chatId) {
			go bot.AddChoice(chatId, update.Message.Text)
		} else {
			// Extract the command from the Message.
			switch update.Message.Command() {
			case "start":
				bot.SendTextMessage(chatId, "Type /decide to begin a decision.")
			case "decide":
				if bot.HasActiveDecision(chatId) {
					go bot.StopDecision(chatId)
				} else {
					go bot.StartDecision(chatId)
				}
			default:
				bot.SendTextMessage(chatId, "I don't know that command")
			}
		}
	}
}
