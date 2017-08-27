package main

import (
	"flag"
	"strings"
	"time"
)

func main() {
	nice := flag.Duration("nice", 0*time.Second, "A duration to pause after each command")
	duration := flag.Duration("duration", 0*time.Second, "A duration for length of time to run")
	workers := flag.Int("workers", 10, "The amount of workers to spawn to run this task")
	flag.Parse()

	punish(*nice, *duration, *workers, strings.Join(flag.Args(), " "))
}
