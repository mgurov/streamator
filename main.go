package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"flag"
	"strconv"

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

type restServer struct {
	quit chan interface{}
}

func startRestServer(port int, wg *sync.WaitGroup, data dataProvider) *restServer {
	wg.Add(1)

	http.HandleFunc("/ticks", func(w http.ResponseWriter, r *http.Request) {

		formatter := logrus.JSONFormatter{}
		logRecords := data.Get()

		w.Header().Add("Content-type", "text/json")

		fmt.Fprint(w, "[")
		for i, rec := range logRecords {
			recBytes, err := formatter.Format(rec)
			if err != nil {
				log.Info("Could not convert log record ", rec, err)
				continue
			}
			w.Write(recBytes)
			if i != len(logRecords)-1 {
				fmt.Fprint(w, ",")
			}
		}
		fmt.Fprint(w, "]")
	})

	server := &http.Server{Addr: ":" + strconv.Itoa(port), Handler: nil}

	go func() {
		log.Info("Starting at port ", port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
		wg.Done()
	}()

	s := &restServer{
		quit: make(chan interface{}),
	}

	go func() {
		<-s.quit
		log.Info("Stopping the web server")
		server.Close()
	}()

	return s
}

func (s *restServer) Stop() {
	s.quit <- nil
}
