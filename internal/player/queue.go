package player

import "sync"

type Queue struct {
	mu    sync.Mutex
	songs []string
}

func NewQueue() *Queue {
	return &Queue{songs: []string{}}
}

func (q *Queue) Add(song string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.songs = append(q.songs, song)
}

func (q *Queue) Next() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.songs) == 0 {
		return ""
	}
	song := q.songs[0]
	q.songs = q.songs[1:]
	return song
}
