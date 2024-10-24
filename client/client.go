package client

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	youtubeApi "google.golang.org/api/youtube/v3"
)

var (
	maxResults = flag.Int64("max-results", 10, "Max Youtube results")
	ApiKey     string
)

func GetApiKey() error {
	if err := godotenv.Load(".env"); err != nil {
		err = errors.New(fmt.Sprintf("godotenv.Load err: %v", err))
		return err
	}
	var exist bool
	ApiKey, exist = os.LookupEnv("API_KEY")
	if exist != true {
		panic(exist)
	}
	return nil
}

func Search(query string) (*youtubeApi.SearchResult, error) {
	if err := GetApiKey(); err != nil {
		return nil, err
	}
	queryPostFix := fmt.Sprintf("%v", query)
	service, err := youtubeApi.NewService(context.Background(), option.WithAPIKey(ApiKey))
	if err != nil {
		err = errors.New(fmt.Sprintf("youtubeApi.NewService err: %v", err))
		return nil, err
	}
	call := service.Search.List([]string{"id,snippet"}).Q(queryPostFix).MaxResults(*maxResults)
	response, err := call.Do()
	if err != nil {
		err = errors.New(fmt.Sprintf("call.Do() err: %v", err))
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
