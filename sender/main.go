package main

import (
	"log"
	"time"
        "p3_ssh_stream/common"
)

func main() {
	addr := "localhost:9000"

	cm := NewConnManager(addr)
	cm.Start()

	// Simulated stream traffic
go func() {
	for {
		if cm.IsConnected() {
			frame := common.Frame{
				Type:    common.FrameData,
				Payload: []byte("hello"),
			}

			log.Println("Sending DATA frame")
			cm.Send(common.EncodeFrame(frame))
		} else {
			log.Println("Connection not established, skipping DATA send")
		}

		time.Sleep(1 * time.Second)
	}
}()

	log.Println("Sender running...")
	select {}
}

