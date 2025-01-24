package decisionbot_test

import (
	"testing"

	"github.com/phk13/decisionbot/decisionbot"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestDecision(t *testing.T) {

	checkDecision := func(t *testing.T, decision *decisionbot.Decision, choices []string, want []string) {
		t.Helper()

		for _, choice := range choices {
			decision.AddChoice(choice)
		}
		got := decision.Decide()
		if !stringInSlice(got, want) {
			t.Errorf("got %s want %s", got, want)
		}
	}

	t.Run("Normal decision", func(t *testing.T) {
		decision := &decisionbot.Decision{}
		choices := []string{"a", "b", "c"}
		checkDecision(t, decision, choices, choices)
	})

	t.Run("Single decision", func(t *testing.T) {
		decision := &decisionbot.Decision{}
		choices := []string{"a"}
		checkDecision(t, decision, choices, choices)
	})

	t.Run("Empty decision", func(t *testing.T) {
		decision := &decisionbot.Decision{}
		choices := []string{}
		checkDecision(t, decision, choices, []string{"You should let me know the choices before asking me to decide..."})
	})
}

func TestAddChoice(t *testing.T) {
	t.Run("Add choice", func(t *testing.T) {
		value := "a"
		decision := &decisionbot.Decision{}
		decision.AddChoice(value)
		got := decision.Choices
		want := []string{value}
		if !stringSlicesEqual(got, want) {
			t.Errorf("got %s want %s", got, want)
		}
	})
}

func TestChoiceNumber(t *testing.T) {
	t.Run("Choice Number", func(t *testing.T) {
		decision := &decisionbot.Decision{Choices: []string{"a", "b", "c", "d"}}
		got := decision.ChoiceNumber()
		want := 4
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
}
