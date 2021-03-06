package main

import (
	"os/exec"
	"sync"
	"time"
)

type Pinger struct {
	available Tristate
	mtx       *sync.RWMutex
	tick      *time.Ticker
	done      chan struct{}
}

func NewPinger(host string, interval time.Duration) *Pinger {
	p := &Pinger{}
	p.available = Undefined
	p.mtx = &sync.RWMutex{}
	p.done = make(chan struct{})
	p.tick = time.NewTicker(interval)

	go func() {
		state_change_pings := 0

		for {
			select {
			case <-p.tick.C:
				err := exec.Command("ping", "-q", "-l 3", "-c 3", "-w 1", host).Run()
				var available Tristate
				if err == nil {
					available = True
				} else {
					available = False
				}

				p.mtx.Lock()
				if p.available != available {
					state_change_pings++
				} else {
					state_change_pings = 0
				}
				if state_change_pings > 4 {
					p.available = available
					state_change_pings = 0
				}
				p.mtx.Unlock()
			case <-p.done:
				return
			}
		}
	}()

	return p
}

func (p *Pinger) Stop() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.tick.Stop()
	close(p.done)
}

func (p *Pinger) GetState() Tristate {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	return p.available
}
