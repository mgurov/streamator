package main

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type ticker struct {
	quit chan interface{}
}

func startTicker(wg *sync.WaitGroup, logger *logrus.Entry) *ticker {
	wg.Add(1)

	t := &ticker{quit: make(chan interface{})}

	ticker := time.NewTicker(1 * time.Second)

	go func() {
		odd := true
		var counter int64
		for {
			select {
			case t := <-ticker.C:
				logger.WithField("type", "repeating").
					WithField("duration", time.Now().Sub(t)).
					Error("Hello, world ", counter, odd)
				counter++
				odd = !odd	
			case <-t.quit:
				ticker.Stop()
				logger.Info("stopping ticker")
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
