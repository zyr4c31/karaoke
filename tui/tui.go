package tui

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/zyr4c31/karaoke/client"
	"github.com/zyr4c31/karaoke/mpv"
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
		fmt.Printf("%v %v\n", command, commands[command])
	}
}

func Run(conn net.Conn) {
	for true {
		PrintHelp()
		time.Sleep(5 * time.Second)
		err := Clear()
		if err != nil {
			log.Fatal(err)
		}

		buffer, err := mpv.SendAndReceive(conn, "get_property", "playlist")
		if err != nil {
			log.Fatal(err)
		}

		trimmedBuf := bytes.TrimRight(buffer, "\x00")

		fmt.Printf("string(buf): %v\n", string(trimmedBuf))

		var response mpv.Reply

		if err := json.Unmarshal(trimmedBuf, &response); err != nil {
			log.Fatal("json.Unmarshal err:", err)
		}

		fmt.Printf("string(response): %v\n", response)
		// fmt.Printf("len(response.Data): %v\n", len(response.Data))

		time.Sleep(time.Second)

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Search for a song: ")
		query, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("text: %v\n", query)
		result, err := client.Search(query)
		if err != nil {
			log.Fatal("net.Dial err:", err)
		}
		fmt.Printf("result.Id.VideoId: %v\n", result.Id.VideoId)

		videoLink := fmt.Sprintf("https://www.youtube.com/watch?v=%v", result.Id.VideoId)
		log.Print(videoLink)

		var modeFile string

		// if len(response.Data) < 1 {
		// 	modeFile = "replace"
		// } else {
		// 	modeFile = "append"
		// }

		modeFile = "replace"
		buffer, err = mpv.SendAndReceive(conn, "loadfile", videoLink, modeFile)

		cmdBytes, err := json.Marshal(buffer)
		if err != nil {
			log.Fatal(err)
		}

		newline := []byte("\n")
		cmdBytes = append(cmdBytes, newline...)

		n, err := conn.Write(cmdBytes)
		if err != nil {
			log.Fatal("conn.Write err:", err)
		}

		fmt.Printf("wrote %d number of bytes\n", n)

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		buf := make([]byte, 65535)
		n, err = conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("string(buf): %v\n", string(buf))
	}
}
