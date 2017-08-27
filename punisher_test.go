package main

import (
	"path/filepath"
	"testing"
	"time"
)

func TestSimplePunisher(t *testing.T) {
	punish(time.Second, 200*time.Millisecond, 1, "touch /tmp/punisher.$(date | md5)")
	//only one should have been created
	files, globErr := filepath.Glob("/tmp/punisher*")
	if len(files) >= 2 || globErr != nil {
		t.Fatalf("Expected only 1 file to be created.")
	}
}
