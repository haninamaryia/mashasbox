package player

import "testing"

func TestQueue_AddAndNext(t *testing.T) {
	q := NewQueue()
	q.Add("song1")
	q.Add("song2")

	if got := q.Next(); got != "song1" {
		t.Errorf("expected song1, got %s", got)
	}
	if got := q.Next(); got != "song2" {
		t.Errorf("expected song2, got %s", got)
	}
	if got := q.Next(); got != "" {
		t.Errorf("expected empty, got %s", got)
	}
}
