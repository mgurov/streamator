package main

import (
	"sync"
	"time"
)

type ticker struct {
	quit chan interface{}
}

func startTicker(wg *sync.WaitGroup) *ticker {
	wg.Add(1)

	t := &ticker{quit: make(chan interface{})}

	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case t := <-ticker.C:
				log.WithField("type", "repeating").
					WithField("duration", time.Now().Sub(t)).
					Error("Hello, world ")
			case <-t.quit:
				ticker.Stop()
				wg.Done()
				return
			}
		}
	}()
	return t
}

func (t *ticker) Stop() {
	t.quit <- nil
}
