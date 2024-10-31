package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/blang/mpv"
	"github.com/zyr4c31/karaoke/client"
	"github.com/zyr4c31/karaoke/oldmpv"
	"github.com/zyr4c31/karaoke/queue"
)

type Playlist struct {
	Current  string `json:"current"`
	Filename string `json:"filename"`
	ID       int    `json:"id"`
	Playing  string `json:"playing"`
}

func main() {
	socket := mpv.NewIPCClient("/tmp/mpvsocket")
	c := mpv.NewClient(socket)
	err := c.Loadfile("https://www.youtube.com/watch?v=jNQXAC9IVRw", mpv.LoadFileModeAppendPlay)
	if err != nil {
		slog.Error(err.Error())
	}

	firstPlaylist, err := c.GetProperty("playlist")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("value: %v\n", firstPlaylist)

	var playlistStruct *Playlist
	err = json.Unmarshal([]byte(firstPlaylist), playlistStruct)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("playlistStruct: %v\n", playlistStruct)

	err = c.Loadfile("https://www.youtube.com/watch?v=VOhdKG24Et8", mpv.LoadFileModeAppend)
	if err != nil {
		log.Fatal(err)
	}

	playlist, err := c.GetProperty("playlist")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("value: %v\n", playlist)

	err = c.SetProperty("playlist", firstPlaylist)
	if err != nil {
		log.Fatal(err)
	}

	isInstalled, err := oldmpv.IsInstalled()
	if err != nil {
		panic(err)
	}

	if !isInstalled {
		panic("mpv not installed")
	}

	flag.Parse()

	var queue queue.Queue
	var pid int
	input := bufio.NewScanner(os.Stdin)
	for {
		// err := tui.Clear()
		fmt.Printf("pid: %v\n", pid)
		// if err != nil {
		// 	panic(err)
		// }

		for idx, song := range queue.Songs {
			fmt.Printf(" %v - %v - %v\n", idx, song.Id, song.Name)
		}
		fmt.Println("add a song to the playlist: ")
		input.Scan()
		result, err := client.Search(input.Text())

		queue.Add(result)
		if err != nil {
			panic(err)
		}
		pid, err = queue.Play()
		if err != nil {
			panic(err)
		}
		process, _ := os.FindProcess(pid)
		fmt.Printf("process: %v\n", process)
	}
}
