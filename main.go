package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {

	quit := make(chan interface{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case t := <-ticker.C:
				println("Hello, world", time.Now().Sub(t).String())
			case <-quit:
				println("Received signal to stop")
				ticker.Stop()
				wg.Done()
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		for _ = range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...")
			quit <- nil
			wg.Done()
		}
	}()

	wg.Wait()

}
