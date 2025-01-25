package decisionbot

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

/* StartDecision starts a new decision process for chatId */
func (bot *DecisionBot) StartDecision(chatId int64) {
	bot.ActiveDecisions[chatId] = &Decision{Lock: sync.Mutex{}}
	log.Printf("Starting decision for %d\n", chatId)
	bot.SendTextMessage(chatId, "OK, let me know the choices to decide. \nWrite each choice in a separate message.")
}

func (bot *DecisionBot) StartDecisionWithChoice(chatId int64, choice string) {
	bot.StartDecision(chatId)
	bot.AddChoice(chatId, choice)
}

/* StopDecision executes a random selection on all input choices for chatId and finishes the active decision */
func (bot *DecisionBot) StopDecision(chatId int64) {
	decision := bot.ActiveDecisions[chatId]
	decision.Lock.Lock()
	defer decision.Lock.Unlock()

	choice := bot.ActiveDecisions[chatId].Decide()
	log.Printf("Decision for %d: %s\n", chatId, choice)

	if decision.ChoiceNumber() > 0 {
		choice = fmt.Sprintf("Decision (%d choices): %s", decision.ChoiceNumber(), choice)
	}
	bot.SendTextMessage(chatId, choice)
	log.Printf("Finishing decision for %d\n", chatId)

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
		if !update.Message.IsCommand() {
			if update.Message.Text == "" {
				continue
			}
			if strings.ToLower(update.Message.Text) == "decide" {
				if bot.HasActiveDecision(chatId) {
					go bot.StopDecision(chatId)
				} else {
					go bot.StartDecision(chatId)
				}
				continue
			}
			if bot.HasActiveDecision(chatId) {
				go bot.AddChoice(chatId, update.Message.Text)
			} else {
				go bot.StartDecisionWithChoice(chatId, update.Message.Text)
			}
		} else {
			// Extract the command from the Message.
			switch update.Message.Command() {
			case "start":
				bot.SendTextMessage(chatId, "Type /yesno to begin a yes or no decision.")
				bot.SendTextMessage(chatId, "Type /decide to begin a multiple choice decision.")
				bot.SendCommandKeyboard(chatId, "I've opened a keyboard with the commands. \nType /closekeyboard if you ever want to close it.")
			case "decide":
				if bot.HasActiveDecision(chatId) {
					go bot.StopDecision(chatId)
				} else {
					go bot.StartDecision(chatId)
				}
			case "yesno":
				go bot.DecideYesOrNo(chatId)
			case "closekeyboard":
				bot.RemoveCommandKeyboard(chatId, "Alright, keyboard removed. Type /start to open it again.")
			default:
				bot.SendTextMessage(chatId, "I don't know that command")
			}
		}
	}
}
