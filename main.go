package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)


var log = logrus.New()

func main() {

	ourHook := newCappedInMemoryRecorderHook(20)
	log.Hooks.Add(ourHook)

	quit := make(chan interface{})
	var wg sync.WaitGroup

	wg.Add(1)
	go tick(&wg, quit)

	wg.Add(1)
	go listenToCtrlC(&wg, quit)

	wg.Wait()

	fmt.Printf("%#v\n", ourHook.Copy())
}

func tick(wg *sync.WaitGroup, quit <-chan interface{}) {
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
}

func listenToCtrlC(wg *sync.WaitGroup, quit chan<- interface{}) {
		signalChan := make(chan os.Signal, 0)
		signal.Notify(signalChan, os.Interrupt)
		<-signalChan
		signal.Stop(signalChan)
		fmt.Println("\nAborting...")
		quit <- nil
		wg.Done()
}
