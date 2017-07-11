package main 

import (
	"time"
)

func main() {
	ticker := time.Tick(1 * time.Second)

	go func() {
		for _ = range ticker {
			println("hello, world", time.Now().String())
		}
	}()

	time.Sleep(20 * time.Second)
}