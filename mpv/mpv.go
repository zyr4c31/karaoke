package mpv

import (
	"fmt"
	"net"
	"os/exec"
)

type MPVCommand struct {
	Command []string `json:"command"`
}

func Connect() (net.Conn, error) {
	conn, err := net.Dial("unix", "/tmp/mpvsocket")
	if err != nil {
		return nil, err
	}
	return conn, nil

}

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

const (
	socketFileName = "/tmp/mpvsocket"
	argIdle        = "--idle"
	argNoTerminal  = "--no-terminal"
	argIPC         = "--input-ipc-server=/tmp/mpvsocket"
)

func StartMpv() (*exec.Cmd, error) {
	cmd := exec.Command("mpv", argIdle, argIPC)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}
