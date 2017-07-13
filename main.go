package main

import (
	"fmt"
	"html"
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
	flag.Parse()

	ourHook := newCappedInMemoryRecorderHook(20)
	log.Hooks.Add(ourHook)

	wg := &sync.WaitGroup{}

	quit := make(chan interface{})
	wg.Add(1)
	go tick(wg, quit)

	quitServer := make(chan interface{})
	startHTTP(*portFlag, wg, quitServer)

	listenToCtrlC()

	log.Info("Stopping the services")
	quit <- nil
	quitServer <- nil
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

func listenToCtrlC() {
	signalChan := make(chan os.Signal, 0)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	signal.Stop(signalChan)
}

func startHTTP(port int, wg *sync.WaitGroup, quit <-chan interface{}) {
	wg.Add(1)

	http.HandleFunc("/ticks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	server := &http.Server{Addr: ":" + strconv.Itoa(port), Handler: nil}

	go func() {
		<-quit
		log.Info("Stopping the web server")
		server.Close()
	}()

	go func() {
		log.Info("Starting at port ", port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
		wg.Done()
	}()

}
