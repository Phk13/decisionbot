package decisionbot

import (
	"math/rand/v2"
	"sync"
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
	return d.Choices[rand.IntN(len(d.Choices))]
}

func (d *Decision) ChoiceNumber() int {
	return len(d.Choices)
}
