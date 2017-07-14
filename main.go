package main

import (
	"os"
	"os/signal"
	"sync"
	"time"

	"flag"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {

	portFlag := flag.Int("port", 8080, "port to listen at")
	memoryCapFlag := flag.Int("cap", 100, "memory cap, e.g. how much log messages to retain")
	flag.Parse()

	var ourHook = newCappedInMemoryRecorderHook(*memoryCapFlag)

	log.Hooks.Add(ourHook)

	wg := &sync.WaitGroup{}

	quit := make(chan interface{})
	wg.Add(1)
	go tick(wg, quit)

	restServer := startRestServer(*portFlag, wg, ourHook)

	listenToCtrlC()

	log.Info("Stopping the services")
	quit <- nil
	restServer.Stop()
	wg.Wait()
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

func listenToCtrlC() {
	signalChan := make(chan os.Signal, 0)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	signal.Stop(signalChan)
}
