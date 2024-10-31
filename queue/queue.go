package queue

import (
	"github.com/zyr4c31/karaoke/oldmpv"
	"google.golang.org/api/youtube/v3"
)

type Queue struct {
	Songs []Song
}

type Song struct {
	Id   string
	Name string
}

func (q *Queue) Add(item *youtube.SearchResult) {
	q.Songs = append(q.Songs, Song{item.Id.VideoId, item.Snippet.Title})
}

// Plays the first song in the queue, returns an error
func (q *Queue) Play() (int, error) {
	pid, err := oldmpv.Play(q.Songs[0].Id)
	if err != nil {
		return 0, err
	}
	return pid, nil
}
