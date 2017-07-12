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
					Error("Hello, world ", )
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
}