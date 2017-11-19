package main

import (
	"log"
	"time"
	"os"
	"os/signal"
	"sync"

	"flag"

	"github.com/sirupsen/logrus"
)

func main() {

	myHostname, err := os.Hostname()
	if err != nil {
		log.Panic("Could not determine own hostname to report to es", err)
	}

	tickDuration := flag.Duration("duration", 5 * time.Second, "duration between log ticks")
	portFlag := flag.Int("port", 8080, "port to listen at")
	elasticURL := flag.String("elasticURL", "http://localhost:9200", "where to send the elastic search logs to")
	elasticIndex := flag.String("elasticIndex", "mylog", "elastic search index")
	elasticReportHost := flag.String("elasticReportHost", myHostname, "host to report to elastic")
	memoryCapFlag := flag.Int("cap", 100, "memory cap, e.g. how much log messages to retain")
	appName := flag.String("app", "generic", "app name to put to the log field app")
	flag.Parse()

	var ourHook = newCappedInMemoryRecorderHook(*memoryCapFlag)

	var log = logrus.New()	
	log.Hooks.Add(ourHook)

	esHook, err := newEsHook(*elasticURL, *elasticReportHost, *elasticIndex)
	if err != nil {
		log.Fatal("Couldn't create es hook:", err)
	}
	log.Hooks.Add(esHook)

	wg := &sync.WaitGroup{}

	logWithApp := log.WithField("app", *appName)

	ticker := startTicker(*tickDuration, wg, logWithApp)

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
