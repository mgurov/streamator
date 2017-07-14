package main

import (
	"os"
	"os/signal"
	"sync"

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

	ticker := startTicker(wg)

	restServer := startRestServer(*portFlag, wg, ourHook)

	listenToCtrlC()

	log.Info("Stopping the services")
	ticker.Stop()
	restServer.Stop()
	wg.Wait()
}

func listenToCtrlC() {
	signalChan := make(chan os.Signal, 0)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	signal.Stop(signalChan)
}
