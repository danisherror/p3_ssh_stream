package main

import (
	"log"
	"time"
)

func main() {
	addr := "localhost:9000"

	cm := NewConnManager(addr)
	cm.Start()

	// Simulated stream traffic
	go func() {
		for {
			cm.Send([]byte("hello\n"))
			time.Sleep(1 * time.Second)
		}
	}()

	log.Println("Sender running...")
	select {}
}

