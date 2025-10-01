package player

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type player struct {
	ctrl       *beep.Ctrl
	format     beep.Format
	mu         sync.Mutex
	paused     bool
	queue      *Queue
	background string
	done       chan bool
}

type Player interface {
	LoadAndPlay(path string) error
	Pause()
	Resume()
	Stop()
	Next()
}

func NewPlayer(q *Queue, background string) Player {
	return &player{
		queue:      q,
		background: background,
		done:       make(chan bool, 1),
	}
}

func (p *player) LoadAndPlay(path string) error {

	// in case of "next"
	if ok := p.mu.TryLock(); ok {
		defer p.mu.Unlock()
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode mp3: %w", err)
	}

	if p.format.SampleRate == 0 {
		p.format = format
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	}

	if p.ctrl != nil {
		p.ctrl.Paused = true
	}

	p.ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	done := make(chan bool)
	speaker.Play(beep.Seq(p.ctrl, beep.Callback(func() {
		done <- true
	})))

	// listen for track end
	go func() {
		<-done
		p.done <- true
		p.playNext()
	}()

	return nil
}

func (p *player) playNext() {
	next := p.queue.Next()
	if next != "" {
		_ = p.LoadAndPlay(next)
	} else if p.background != "" {
		_ = p.LoadAndPlay(p.background)
	}
}

func (p *player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ctrl != nil {
		p.ctrl.Paused = true
		p.paused = true
	}
}

func (p *player) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ctrl != nil {
		p.ctrl.Paused = false
		p.paused = false
	}
}

func (p *player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ctrl = nil
}

func (p *player) Next() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// stop current track
	if p.ctrl != nil {
		p.ctrl.Paused = true
		p.ctrl = nil
	}

	// play next in queue or fallback
	p.playNext()
}
