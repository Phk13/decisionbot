package main

import (
	"log"
	"os"

	"github.com/phk13/decisionbot/decisionbot"
)

func main() {
	log.Println("Starting decision bot...")
	log.Println("Fetching API token...")

	dbot, err := decisionbot.NewDecisionBot(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	dbot.ListenAndDecide()
}
