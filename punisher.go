package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/xid"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func punish(nice, duration time.Duration, verbose bool, workers int, cmd string) {
	// setup our punisher
	p := &punisher{
		nice:     nice,
		duration: duration,
		workers:  workers,
		cmd:      cmd,
		signals:  make(chan os.Signal, 1),
		wg:       sync.WaitGroup{},
		stop:     make(chan bool),
		verbose:  verbose,
	}

	if cmd != "" {
		// bring the pain
		p.pain()
	}
}

// holds our goods
type punisher struct {
	nice     time.Duration
	duration time.Duration
	workers  int
	cmd      string
	success  int
	verbose  bool
	failure  int
	lock     sync.Mutex
	stop     chan bool
	signals  chan os.Signal
	wg       sync.WaitGroup
}

func (p *punisher) pain() {
	for id := 1; id <= p.workers; id++ {
		p.wg.Add(1)
		go func(id int, nice time.Duration, cmd string) {
			for {
				select {
				case <-p.stop:
					p.wg.Done()
					return
				default:
					<-time.After(nice)
					cmdParsed := p.prepCmd(cmd)
					command := exec.Command("sh", "-c", cmdParsed)
					output, _ := command.CombinedOutput()
					if !command.ProcessState.Success() {
						p.lock.Lock()
						p.failure++
						p.lock.Unlock()
						if p.verbose {
							fmt.Println("id#", id, ">", string(output))
						}
					} else {
						p.lock.Lock()
						p.success++
						p.lock.Unlock()
					}
				}
			}
		}(id, p.nice, p.cmd)
	}

	go p.track()
	go p.shutdown()
	p.wg.Wait()
}

func (p *punisher) track() {
	p.lock.Lock()
	bar := pb.New(p.failure + p.success)
	p.lock.Unlock()
	bar.ShowTimeLeft = false
	bar.ShowPercent = true
	bar.ShowSpeed = true
	bar.Start()
	for {
		p.lock.Lock()
		bar.Total = int64(p.failure + p.success)
		bar.Set(p.success)
		p.lock.Unlock()
		<-time.After(time.Second)
	}
}

func (p *punisher) shutdown() {
	signal.Notify(p.signals, syscall.SIGINT, syscall.SIGTERM)
	if p.duration == 0*time.Second {
		<-p.signals
	} else {
		select {
		case <-p.signals:
			break
		case <-time.After(p.duration):
			break
		}
	}
	for workers := 1; workers <= p.workers; workers++ {
		p.stop <- true
	}
}

func (p *punisher) prepCmd(cmd string) string {

	commandOptions := struct {
		UniqID string
		Date   time.Time
	}{
		UniqID: xid.New().String(),
		Date:   time.Now(),
	}

	tmpl, parseErr := template.New(commandOptions.UniqID).Parse(cmd)
	if parseErr != nil {
		return cmd
	}

	cmdParsed := new(bytes.Buffer)
	executionErr := tmpl.Execute(cmdParsed, commandOptions)
	if executionErr != nil {
		return cmd
	}

	return cmdParsed.String()
}
