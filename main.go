package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/zyr4c31/karaoke/queue"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	maxResults = flag.Int64("max-results", 5, "Max Youtube results")
)

var apiKey string

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	var exist bool
	apiKey, exist = os.LookupEnv("API_KEY")
	if exist != true {
		panic(exist)
	}

	flag.Parse()

	var queue queue.Queue
	input := bufio.NewScanner(os.Stdin)
	for {
		err := systemClear()
		if err != nil {
			panic(err)
		}

		for idx, song := range queue.Songs {
			fmt.Printf(" %v - %v - %v\n", idx, song.Id, song.Name)
		}
		fmt.Println("add a song to the playlist: ")
		input.Scan()
		result, err := search(input.Text())
		queue.Add(*result)
		if err != nil {
			panic(err)
		}
		queue.Play()
	}
}

func search(query string) (*youtube.SearchResult, error) {
	queryPostFix := fmt.Sprintf("%v", query)
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	call := service.Search.List([]string{"id,snippet"}).Q(queryPostFix).MaxResults(*maxResults)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			return item, nil
		default:
			continue
		}
	}
	err = errors.New("no video found")
	return nil, err
}

func startMpv() error {
	mpv := exec.Command("mpv --idle --input-ipc-server=/tmp/mpvsocket")
	// arg := fmt.Sprintf("https://www.youtube.com/watch?v=%v", id)
	// mpv := exec.Command("mpv", arg)
	if err := mpv.Start(); err != nil {
		return err
	}
	return nil
}

func systemClear() error {
	out, err := exec.Command("clear").Output()
	if err != nil {
		return err
	}
	if _, err = os.Stdout.Write(out); err != nil {
		return err
	}
	return nil
}
