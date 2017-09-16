package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/mgutz/ansi"
	"github.com/rs/xid"
)

// create a fake signal
type punishSignal struct{}

func (p *punishSignal) Signal() {}

func punish(p *punisher) {
	// init
	p.signals = make(chan os.Signal, 1)
	p.wg = sync.WaitGroup{}
	p.stop = make(chan bool)

	if p.loopIncrement == 0 {
		// really shouldn't be zero by default ...
		p.loopIncrement = 1
	}

	if p.cmd != "" {
		// bring the pain
		p.pain()
	}
}

// holds our goods
type punisher struct {
	nice          time.Duration
	duration      time.Duration
	workers       int
	success       int
	failure       int
	verbose       bool
	retry         bool
	stop          chan bool
	signals       chan os.Signal
	wg            sync.WaitGroup
	metricsLock   sync.Mutex
	loopLock      sync.Mutex
	msgLock       sync.Mutex
	cmd           string
	loopIncrement int
	loopName      string
	loopEnd       int
	loopIndex     int
}

func (p *punisher) pain() {
	for id := 1; id <= p.workers; id++ {
		p.wg.Add(1)
		go func(id int, nice time.Duration, cmd, loopName string) {
			// lets get to work
			retry := false
			var loopIndex int
			var loopError error
			for {
				if !retry {
					// save aside our loop value
					loopIndex, loopError = p.getLoopIndex()
				}
				select {
				case <-p.stop:
					// did we get a shutdown sig?
					p.wg.Done()
					return
				default:
					if loopError != nil {
						p.wg.Done()
						return
					}
					// get our cmd ready
					cmdParsed := p.prepCmd(cmd, loopName, loopIndex)
					command := exec.Command("sh", "-c", cmdParsed)
					output, _ := command.CombinedOutput()
					ok := command.ProcessState.Success()

					// count it!
					if !ok {
						p.metricsLock.Lock()
						p.failure++
						p.metricsLock.Unlock()

					} else {
						p.metricsLock.Lock()
						p.success++
						p.metricsLock.Unlock()
					}

					// show the status
					if p.verbose || (!ok && p.retry) {
						lime := ansi.ColorCode("green")
						red := ansi.ColorCode("red")
						reset := ansi.ColorCode("reset")
						// if we retry, and there were failures, let everybody know
						if ok {
							fmt.Print(lime, "[OK] ", reset, cmdParsed, "\n", string(output), "\n")
							retry = false
						} else {
							fmt.Print(red, "[FAILED] ", reset, cmdParsed, "\n", string(output), "\n")
							if p.retry {
								retry = true
							}
						}
					}

					// netflix && chill
					if nice != 0*time.Second {
						<-time.After(nice)
					}

					// do we need to retry? If not, continue on ...
					if ok || !p.retry {
						break
					}
				}
			}
		}(id, p.nice, p.cmd, p.loopName)
	}

	go p.shutdown()
	p.wg.Wait()
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

	// our progress bar
	p.stop <- true
}

func (p *punisher) prepCmd(cmd, loopName string, loopIndex int) string {

	commandOptions := struct {
		UniqID    string
		Date      time.Time
		LoopName  string
		LoopIndex int
	}{
		UniqID:    xid.New().String(),
		Date:      time.Now(),
		LoopName:  loopName,
		LoopIndex: loopIndex,
	}

	// we need to allow the user to replace a uniq string with the loop iter(for nested loops)
	if loopName != "" {
		cmd = strings.Replace(cmd, loopName, strconv.Itoa(loopIndex), -1)
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

func (p *punisher) getLoopIndex() (int, error) {
	// ok, lets lock it up
	p.loopLock.Lock()
	idx := p.loopIndex
	ends := p.loopEnd
	name := p.loopName
	p.loopIndex = idx + p.loopIncrement
	p.loopLock.Unlock()

	// keep on keeping on ...
	if name == "" || ends == 0 {
		return idx, nil
	}

	// check
	if idx > ends {
		// be gone!
		go func() {
			p.signals <- syscall.SIGKILL
		}()
		return idx, fmt.Errorf("Reached the end of the loop")
	}

	// return the goods
	return idx, nil
}
