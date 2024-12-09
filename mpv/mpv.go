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
	PlaybackControlStop   = "stop"
	PlaylistManipLoadFile = "loadfile"
	PlaylistManipNext     = "playlist-next"
	PlaylistManipPrev     = "playlist-prev"
	PlaylistManipRemove   = "playlist-clear"
	PropertyManipCycle    = "cycle"
	// Append the file to the playlist
	PlaylistManipLoadfileFlagAppend = "append"
	PropertyNamePause               = "pause"
	argIPC                          = "--input-ipc-server=/tmp/mpvsocket"
	argIdle                         = "--idle"
	argNoTerminal                   = "--no-terminal"
	socketFileName                  = "/tmp/mpvsocket"
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
func Send(conn net.Conn, command string, args ...string) error {
	message := Command{
		Command: []string{command},
	}

	for _, arg := range args {
		message.Command = append(message.Command, arg)
	}

	cmdBytes, err := json.Marshal(message)
	if err != nil {
		return nil
	}

	newline := []byte("\n")
	cmdBytes = append(cmdBytes, newline...)

	_, err = conn.Write(cmdBytes)
	if err != nil {
		return nil
	}

	return nil
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

	_, err = conn.Write(cmdBytes)
	if err != nil {
		return nil, err
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	buf := make([]byte, 65535)
	_, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}

	fmt.Printf("string(buf): %v\n", string(buf))
	trimmedBuf := bytes.TrimRight(buf, "{")
	fmt.Printf("string(trimmedBuf): %v\n", string(trimmedBuf))

	return trimmedBuf, nil
}
