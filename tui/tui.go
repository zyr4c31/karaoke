package tui

import (
	"log"
	"os"
	"os/exec"
)

// clears the screen
func Clear() error {
	out, err := exec.Command("clear").Output()
	if err != nil {
		return err
	}
	if _, err = os.Stdout.Write(out); err != nil {
		return err
	}
	return nil
}

func PrintHelp() {
	commands := map[string]string{
		"help":   "print this",
		"query":  "search",
		"toggle": "toggle play/pause",
	}

	for command := range commands {
		log.Printf("%v %v", command, commands[command])
	}
}
