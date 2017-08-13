package main

import (
	"os"
	"os/signal"
	"sync"

	"flag"

	"github.com/sirupsen/logrus"
)

func main() {

	portFlag := flag.Int("port", 8080, "port to listen at")
	memoryCapFlag := flag.Int("cap", 100, "memory cap, e.g. how much log messages to retain")
	appName := flag.String("app", "generic", "app name to put to the log field app")
	flag.Parse()

	var ourHook = newCappedInMemoryRecorderHook(*memoryCapFlag)

	var log = logrus.New()	
	log.Hooks.Add(ourHook)

	esHook, err := newEsHook()

	if err != nil {
		log.Fatal("Couldn't create es hook:", err)
	}

	log.Hooks.Add(esHook)

	wg := &sync.WaitGroup{}

	logWithApp := log.WithField("app", *appName)

	ticker := startTicker(wg, logWithApp)

	restServer := startRestServer(*portFlag, wg, ourHook, logWithApp)

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
