package mpv

import (
	"fmt"
	"os/exec"
)

func IsInstalled() (bool, error) {
	if err := exec.Command("mpv").Run(); err != nil {
		return false, err
	}
	return true, nil
}

func Play(id string) (int, error) {
	arg := fmt.Sprintf("https://www.youtube.com/watch?v=%v", id)
	mpv := exec.Command("mpv", arg)
	if err := mpv.Start(); err != nil {
		return 0, err
	}
	return mpv.Process.Pid, nil
}

func StartMpv() error {
	mpv := exec.Command("mpv --idle --input-ipc-server=/tmp/mpvsocket")
	// arg := fmt.Sprintf("https://www.youtube.com/watch?v=%v", id)
	// mpv := exec.Command("mpv", arg)
	if err := mpv.Start(); err != nil {
		return err
	}
	return nil
}
