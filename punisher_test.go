package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSimplePunisher(t *testing.T) {
	file := "touch /tmp/punisher.$(date | md5)"
	punish(time.Second, 200*time.Millisecond, true, 1, file)
	//only one should have been created
	files, globErr := filepath.Glob("/tmp/punisher*")
	if len(files) >= 2 || globErr != nil {
		t.Fatalf("Expected only 1 file to be created.")
	}

	// cleanup test run, failure ... :shrug:
	os.Remove(file)
}
