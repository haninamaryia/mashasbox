package main

import (
	"fmt"

	"github.com/haninamaryia/mashasbox/internal/player"
	"github.com/haninamaryia/mashasbox/internal/server"
)

func main() {
	q := player.NewQueue()
	p := player.NewPlayer(q, "./music/background.mp3") // fallback background
	srv := server.NewServer(p, q)

	fmt.Println("Starting jukebox on :8080...")
	if err := srv.Run(":8080"); err != nil {
		panic(err)
	}
}
