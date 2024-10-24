package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/zyr4c31/karaoke/client"
	"github.com/zyr4c31/karaoke/mpv"
	"github.com/zyr4c31/karaoke/queue"
	"github.com/zyr4c31/karaoke/tui"
)

func main() {
	isInstalled, err := mpv.IsInstalled()
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
		err := tui.Clear()
		fmt.Printf("pid: %v\n", pid)
		if err != nil {
			panic(err)
		}

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
