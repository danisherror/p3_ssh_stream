package main

import (
	"log"
	"time"
)

func main() {
	addr := "localhost:9000"

cm := NewConnManager(addr)
cm.Start()

// Wait until connected
for !cm.IsConnected() {
    time.Sleep(100 * time.Millisecond)
}

log.Println("Connection established. Creating streams...")

stream1 := cm.CreateStream()
stream2 := cm.CreateStream()

	go func() {
		for {
			if cm.IsConnected() {
				stream1.Send([]byte("hello from stream1"))
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			if cm.IsConnected() {
				stream2.Send([]byte("hello from stream2"))
			}
			time.Sleep(1500 * time.Millisecond)
		}
	}()

	log.Println("Sender running...")
	select {}
}

