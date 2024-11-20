package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zyr4c31/karaoke/client"
	"github.com/zyr4c31/karaoke/mpv"
	"github.com/zyr4c31/karaoke/tui"
)

func main() {
	conn, err := mpv.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	buffer, err := mpv.SendAndReceive(conn, "get_property", "playlist")
	if err != nil {
		log.Fatal(err)
	}

	trimmedBuf := bytes.TrimRight(buffer, "\x00")

	fmt.Printf("string(buf): %v\n", string(trimmedBuf))

	reader := bytes.NewReader(trimmedBuf)

	decoder := json.NewDecoder(reader)

	var reply mpv.Reply
	var event mpv.Event

	for token, err := decoder.Token(); err == nil; {
		if token == "data" {
			if err := json.Unmarshal(trimmedBuf, &reply); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("reply: %v\n", reply)
			break
		}
		if token == "event" {
			if err := json.Unmarshal(trimmedBuf, &event); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("event: %v\n", event)
			break
		}
		token, err = decoder.Token()
	}

	for true {
		tui.PrintHelp()
		time.Sleep(10 * time.Second)
		err := tui.Clear()
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

		var modeFile string

		if len(response.Data) < 1 {
			modeFile = "replacej"
		} else {
			modeFile = "append"
		}

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
