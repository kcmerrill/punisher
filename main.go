package main

import (
	"flag"
	"strings"
	"time"
)

func main() {
	nice := flag.Duration("nice", 0*time.Second, "A duration to pause after each command")
	duration := flag.Duration("duration", 0*time.Second, "A duration for length of time to run")
	status := flag.Duration("status", 10*time.Second, "Interval to show the status")
	workers := flag.Int("workers", 20, "The amount of workers to spawn to run this task")
	loopName := flag.String("loop-name", "", "Name of the loop to use")
	loopEnd := flag.Int("loop-end", 0, "Size to end loop at (requires loop-name flag)")
	loopStartsAt := flag.Int("loop-starts-at", 0, "Starting value of where the loop should start at (requires loop-name flag")
	loopIncrement := flag.Int("loop-increment", 1, "The amount to increment loop iterator by (requires loop-name flag)")
	verbose := flag.Bool("verbose", false, "Display the results of each command")
	retry := flag.Bool("retry", false, "Retry on failure?")
	flag.Parse()

	punish(&punisher{
		nice:          *nice,
		duration:      *duration,
		status:        *status,
		verbose:       *verbose,
		retry:         *retry,
		workers:       *workers,
		cmd:           strings.Join(flag.Args(), " "),
		loopName:      *loopName,
		loopIncrement: *loopIncrement,
		loopEnd:       *loopEnd,
		loopIndex:     *loopStartsAt,
	})
}
