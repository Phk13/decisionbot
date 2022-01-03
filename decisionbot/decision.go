package decisionbot

import (
	"math/rand"
	"sync"
	"time"
)

type Decision struct {
	Choices []string
	Lock    sync.Mutex
}

/* AddChoice adds a new choice string into the Choices slice */
func (d *Decision) AddChoice(choice string) {
	d.Choices = append(d.Choices, choice)
}

/* Decide returns a random choice from the Choices slice */
func (d *Decision) Decide() string {
	if len(d.Choices) == 0 {
		return "You should let me know the choices before asking me to decide..."
	}
	rand.Seed(time.Now().Unix())
	return d.Choices[rand.Intn(len(d.Choices))]
}

func (d *Decision) ChoiceNumber() int {
	return len(d.Choices)
}
