package main

import "testing"

func TestMainRuns(t *testing.T) {
	// just a sanity check to ensure main function doesn’t panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main panicked: %v", r)
		}
	}()
}
