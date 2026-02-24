package cc

import (
	"fmt"
)

// GetGreeting generates a greeting message.
// isFormal: whether to use a formal tone.
// isMorning: whether it is morning.
// isEvening: whether it is evening.
func GetGreeting(name string, isFormal bool, isMorning bool, isEvening bool) string {
	if isFormal && isMorning {
		return fmt.Sprintf("Good morning, %s.", name)
	}
	if isFormal && isEvening {
		return fmt.Sprintf("Good evening, %s.", name)
	}
	if isFormal {
		return fmt.Sprintf("Hello, %s.", name)
	}
	if isMorning {
		return fmt.Sprintf("Hey %s, good morning!", name)
	}
	if isEvening {
		return fmt.Sprintf("Hey %s, good evening!", name)
	}
	return fmt.Sprintf("Hey, %s!", name)
}
