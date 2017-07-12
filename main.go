package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var b *logrus.Hook

type cappedInMemoryRecorderHook struct {
	m sync.Mutex
	records []*logrus.Entry
	wIndex  int
	owerwrites bool
}

func newCappedInMemoryRecorderHook(c int) *cappedInMemoryRecorderHook {
	if c <= 0 {
		panic("cappedInMemoryRecorderHook should be not empty size of")
	}
	return &cappedInMemoryRecorderHook{
		records: make([]*logrus.Entry, c),
	}
}

func (h *cappedInMemoryRecorderHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *cappedInMemoryRecorderHook) Fire(e *logrus.Entry) error {
	h.m.Lock()
	defer h.m.Unlock()

	h.records[h.wIndex] = e

	h.wIndex ++
	if h.wIndex >= len(h.records) {
		h.wIndex = 0
		h.owerwrites = true
	}
	return nil
}

func (h *cappedInMemoryRecorderHook) Copy() []*logrus.Entry {
	h.m.Lock()
	defer h.m.Unlock()

	if (!h.owerwrites) {
		//todo: test me now
		return h.records[:h.wIndex]
	} 

	result := make([]*logrus.Entry, len(h.records))

	copy(result, h.records[h.wIndex:])	
	copy(result[h.wIndex:], h.records[:h.wIndex])	
	return result
}

func main() {

	ourHook := newCappedInMemoryRecorderHook(20)
	var log = logrus.New()
	log.Hooks.Add(ourHook)

	quit := make(chan interface{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case t := <-ticker.C:
				log.WithField("type", "repeating").
					WithField("duration", time.Now().Sub(t)).
					Error("Hello, world ")
			case <-quit:
				ticker.Stop()
				wg.Done()
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		signalChan := make(chan os.Signal, 0)
		signal.Notify(signalChan, os.Interrupt)
		<-signalChan
		signal.Stop(signalChan)
		fmt.Println("\nAborting...")
		quit <- nil
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("%#v\n", ourHook.Copy())
}
