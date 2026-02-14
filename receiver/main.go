package main

import (
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Receiver listening on :9000")

	for {
		conn, _ := ln.Accept()
		log.Println("Client connected")

		go func(c net.Conn) {
			buf := make([]byte, 1024)
			for {
				n, err := c.Read(buf)
				if err != nil {
					log.Println("Client disconnected")
					return
				}
				log.Printf("Received: %s", string(buf[:n]))
			}
		}(conn)
	}
}

