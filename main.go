package main 

import (
	"time"
)

func main() {
	ticker := time.Tick(1 * time.Second)

	for _ = range ticker {
		println("hello, world", time.Now().String())
	}
}