package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSimplePunisher(t *testing.T) {
	file := "touch /tmp/punisher.$(date | md5)"
	p := &punisher{
		nice:     time.Second,
		duration: 200 * time.Millisecond,
		verbose:  true,
		workers:  1,
		cmd:      file,
	}

	punish(p)

	//only one should have been created
	files, globErr := filepath.Glob("/tmp/punisher*")
	if len(files) >= 2 || globErr != nil {
		// cleanup test run, failure ... :shrug:
		os.Remove(file)
		t.Fatalf("Expected only 1 file to be created.")
	}

	// cleanup test run, failure ... :shrug:
	os.Remove(file)
}
