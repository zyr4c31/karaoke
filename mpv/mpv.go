package mpv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"
)

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

type Command struct {
	Command []string `json:"command"`
}

type Reply struct {
	Data      []Song `json:"data"`
	RequestID int    `json:"request_id"`
	Error     string `json:"error"`
}

type Song struct {
	FileName string `json:"filename"`
	Title    string `json:"title"`
	ID       int    `json:"id"`
	Current  bool   `json:"current"`
	Playing  bool   `json:"playing"`
}

type Event struct {
	Event     string `json:"event"`
	EventName string `json:"event_name"`
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

func SendAndReceive(conn net.Conn, command string, args ...string) ([]byte, error) {
	message := Command{
		Command: []string{command},
	}

	for _, arg := range args {
		message.Command = append(message.Command, arg)
	}

	cmdBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	newline := []byte("\n")
	cmdBytes = append(cmdBytes, newline...)

	n, err := conn.Write(cmdBytes)
	if err != nil {
		return nil, err
	}

	fmt.Printf("wrote %d number of bytes\n", n)

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	buf := make([]byte, 65535)
	n, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}

	trimmedBuf := bytes.TrimRight(buf, "\x00")

	fmt.Printf("string(buf): %v\n", string(trimmedBuf))
	return buf, nil
}
