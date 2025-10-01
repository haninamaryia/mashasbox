package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/haninamaryia/mashasbox/internal/player"
	"github.com/stretchr/testify/assert"
)

type mockPlayer struct {
	LoadAndPlayFunc func(path string) error
	PauseFunc       func()
	ResumeFunc      func()
	StopFunc        func()
	NextFunc        func()
}

func (m *mockPlayer) LoadAndPlay(path string) error {
	if m.LoadAndPlayFunc != nil {
		return m.LoadAndPlayFunc(path)
	}
	return nil
}

func (m *mockPlayer) Pause() {
	if m.PauseFunc != nil {
		m.PauseFunc()
	}
}

func (m *mockPlayer) Resume() {
	if m.ResumeFunc != nil {
		m.ResumeFunc()
	}
}

func (m *mockPlayer) Stop() {
	if m.StopFunc != nil {
		m.StopFunc()
	}
}

func (m *mockPlayer) Next() {
	if m.NextFunc != nil {
		m.NextFunc()
	}
}

func TestHandlePlay(t *testing.T) {
	mockP := &mockPlayer{
		LoadAndPlayFunc: func(path string) error {
			if path == "invalid.mp3" {
				return assert.AnError
			}
			return nil
		},
	}
	q := player.NewQueue()
	srv := NewServer(mockP, q)

	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{"Valid Path", `{"path":"commas.mp3"}`, http.StatusOK},
		{"Missing Path", `{}`, http.StatusBadRequest},
		{"Invalid Path", `{"path":"invalid.mp3"}`, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/play", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			srv.engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestHandlePause(t *testing.T) {
	mockP := &mockPlayer{
		PauseFunc: func() {},
	}
	q := player.NewQueue()
	srv := NewServer(mockP, q)

	req, _ := http.NewRequest("POST", "/pause", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleResume(t *testing.T) {
	mockP := &mockPlayer{
		ResumeFunc: func() {},
	}
	q := player.NewQueue()
	srv := NewServer(mockP, q)

	req, _ := http.NewRequest("POST", "/resume", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleQueue(t *testing.T) {
	mockP := &mockPlayer{}
	q := player.NewQueue()
	srv := NewServer(mockP, q)

	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{"Valid Path", `{"path":"commas.mp3"}`, http.StatusOK},
		{"Missing Path", `{}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/queue", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			srv.engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
