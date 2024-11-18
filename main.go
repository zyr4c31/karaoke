package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/zyr4c31/karaoke/client"
	"github.com/zyr4c31/karaoke/mpv"
)

type MPVResponse struct {
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

func togglePause(conn net.Conn) {
	command := mpv.MPVCommand{
		Command: []string{"cycle", "pause"},
	}

	cmdBytes, err := json.Marshal(command)
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

	trimmedBuf := bytes.TrimRight(buf, "\x00")

	fmt.Printf("string(buf): %v\n", string(trimmedBuf))
}

func main() {
	conn, err := mpv.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	command := mpv.MPVCommand{
		Command: []string{"cycle", "pause"},
	}

	cmdBytes, err := json.Marshal(command)
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

	trimmedBuf := bytes.TrimRight(buf, "\x00")

	fmt.Printf("string(buf): %v\n", string(trimmedBuf))

	var response MPVResponse

	if err := json.Unmarshal(trimmedBuf, &response); err != nil {
		log.Fatal("json.Unmarshal err:", err)
	}

	fmt.Printf("string(response): %v\n", response)
	fmt.Printf("len(response.Data): %v\n", len(response.Data))

	for true {
		// cmd := exec.Command("clear")
		// cmd.Stdout = os.Stdout
		// cmd.Run()

		command := mpv.MPVCommand{
			Command: []string{"get_property", "playlist"},
		}

		cmdBytes, err := json.Marshal(command)
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

		trimmedBuf := bytes.TrimRight(buf, "\x00")

		fmt.Printf("string(buf): %v\n", string(trimmedBuf))

		var response MPVResponse

		if err := json.Unmarshal(trimmedBuf, &response); err != nil {
			log.Fatal("json.Unmarshal err:", err)
		}

		fmt.Printf("string(response): %v\n", response)
		fmt.Printf("len(response.Data): %v\n", len(response.Data))

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

		if len(response.Data) < 1 {
			command = mpv.MPVCommand{
				Command: []string{"loadfile", videoLink, "replace"},
			}
		} else {
			command = mpv.MPVCommand{
				Command: []string{"loadfile", videoLink, "append"},
			}
		}

		cmdBytes, err = json.Marshal(command)
		if err != nil {
			log.Fatal(err)
		}

		newline = []byte("\n")
		cmdBytes = append(cmdBytes, newline...)

		n, err = conn.Write(cmdBytes)
		if err != nil {
			log.Fatal("conn.Write err:", err)
		}

		fmt.Printf("wrote %d number of bytes\n", n)

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		buf = make([]byte, 65535)
		n, err = conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("string(buf): %v\n", string(buf))
	}
}
