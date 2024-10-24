package tui

import (
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
