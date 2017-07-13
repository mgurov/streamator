package main

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"flag"
	"strconv"
)

var log = logrus.New()

func main() {

	portFlag := flag.Int("port", 8080, "port to listen at")
	flag.Parse()

	ourHook := newCappedInMemoryRecorderHook(20)
	log.Hooks.Add(ourHook)

	quit := make(chan interface{})
	var wg sync.WaitGroup

	wg.Add(1)
	go tick(&wg, quit)

	wg.Add(1)
	go listenToCtrlC(&wg, quit)

	go startHTTP(*portFlag)

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

func startHTTP(port int) {
	http.HandleFunc("/ticks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Info("Starting at port ", port)
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
}
