package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haninamaryia/mashasbox/internal/player"
)

type Server struct {
	player player.Player
	queue  *player.Queue
	engine *gin.Engine
}

func NewServer(p player.Player, q *player.Queue) *Server {
	engine := gin.Default()
	s := &Server{player: p, queue: q, engine: engine}

	// register routes
	engine.POST("/play", s.handlePlay)
	engine.POST("/pause", s.handlePause)
	engine.POST("/resume", s.handleResume)
	engine.POST("/queue", s.handleQueue)
	s.engine.POST("/next", func(c *gin.Context) {
		p.Next()
		c.JSON(http.StatusOK, gin.H{"status": "skipped to next track"})
	})

	return s
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}

func (s *Server) handlePlay(c *gin.Context) {
	var body struct {
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	if err := s.player.LoadAndPlay(body.Path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "playing", "path": body.Path})
}

func (s *Server) handlePause(c *gin.Context) {
	s.player.Pause()
	c.JSON(http.StatusOK, gin.H{"message": "paused"})
}

func (s *Server) handleResume(c *gin.Context) {
	s.player.Resume()
	c.JSON(http.StatusOK, gin.H{"message": "resumed"})
}

func (s *Server) handleQueue(c *gin.Context) {
	var body struct {
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	s.queue.Add(body.Path)
	c.JSON(http.StatusOK, gin.H{"message": "queued", "path": body.Path})
}
