package main

import (
	"bytes"
	"fmt"
	"launchpad.net/goyaml"
	"log"
	"os/exec"
	"sync"
	"time"
)

var (
	gitRepo     = "/home/mero/src/www-nnev"
	uniLocation = Location{
		Address: "Im Neuenheimer Feld 368, 69120 Heidelberg",
		Lat:     49.41759,
		Lon:     8.66834,
	}
)

type LocationPoller struct {
	loc  Location
	mtx  *sync.RWMutex
	tick *time.Ticker
	done chan struct{}
}

func NewLocationPoller(interval time.Duration) *LocationPoller {
	p := &LocationPoller{}
	p.mtx = &sync.RWMutex{}
	p.tick = time.NewTicker(interval)
	p.done = make(chan struct{})

	p.Poll()

	go func() {
		for {
			select {
			case <-p.tick.C:
				p.Poll()
			case <-p.done:
				return
			}
		}
	}()

	return p
}

func (p *LocationPoller) Poll() {
	var loc *Location

	defer func() {
		if loc == nil {
			loc = &uniLocation
		}

		p.mtx.Lock()
		p.loc = *loc
		p.mtx.Unlock()
	}()

	site, err := p.GetStammtisch()
	if err != nil {
		log.Println(err)
		return
	}

	if site == "" {
		return
	}

	loc, err = p.GetLocation(site)
	if err != nil {
		log.Println(err)
		return
	}
}

func (p *LocationPoller) Get() Location {
	var loc Location
	p.mtx.RLock()
	loc = p.loc
	p.mtx.RUnlock()
	return loc
}

func (p *LocationPoller) Stop() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.tick.Stop()
	close(p.done)
}

func (p *LocationPoller) GetLocation(site string) (loc *Location, err error) {
	buf := new(bytes.Buffer)

	cmd := exec.Command("/usr/bin/git", "cat-file", "-p", fmt.Sprintf("master:stammtisch_%s.md", site))
	cmd.Dir = gitRepo
	cmd.Stdout = buf
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	loc = &Location{}

	err = goyaml.Unmarshal(buf.Bytes(), loc)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

func (p *LocationPoller) GetStammtisch() (site string, err error) {
	buf := new(bytes.Buffer)

	cmd := exec.Command("/usr/bin/termine", "location")
	cmd.Stdout = buf
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(buf.Bytes())), nil
}
