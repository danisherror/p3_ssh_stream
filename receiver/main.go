package main

import (
	"bufio"
	"log"
	"net"
	"strings"
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

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Client disconnected")
			return
		}

		line = strings.TrimSpace(line)

		if line == "PING" {
			conn.Write([]byte("PONG\n"))
			continue
		}

		log.Println("Received:", line)
	}
}

