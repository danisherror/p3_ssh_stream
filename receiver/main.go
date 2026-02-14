package main

import (
	"log"
	"net"

	"p3_ssh_stream/common"
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

	for {
		frame, err := common.DecodeFrame(conn)
		if err != nil {
			log.Println("Client disconnected")
			return
		}

		switch frame.Type {

		case common.FramePing:
			pong := common.Frame{
				Type:    common.FramePong,
				Payload: nil,
			}
			conn.Write(common.EncodeFrame(pong))

		case common.FrameData:
			log.Println("Received DATA:", string(frame.Payload))
		}
	}
}

