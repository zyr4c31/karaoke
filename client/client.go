package client

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"google.golang.org/api/option"
	youtubeApi "google.golang.org/api/youtube/v3"
)

var (
	maxResults = flag.Int64("max-results", 1, "Max Youtube results")
	ApiKey     string
)

func Search(query string) (*youtubeApi.SearchResult, error) {
	queryPostFix := fmt.Sprintf("%v", query)
	service, err := youtubeApi.NewService(context.Background(), option.WithAPIKey(ApiKey))
	if err != nil {
		err = errors.New(fmt.Sprintf("youtubeApi.NewService err: %v", err))
		return nil, err
	}
	call := service.Search.List([]string{"id", "snippet"}).Q(queryPostFix).Type("video").MaxResults(*maxResults)
	response, err := call.Do()
	if err != nil {
		err = errors.New(fmt.Sprintf("call.Do() err: %v", err))
		return nil, err
	}
	for _, item := range response.Items {
		fmt.Printf("item: %v\n", item)
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
