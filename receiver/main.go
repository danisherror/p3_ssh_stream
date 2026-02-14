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

	streams := make(map[uint32]bool)

	for {
		frame, err := common.DecodeFrame(conn)
		if err != nil {
			log.Println("Client disconnected")
			return
		}

		switch frame.Type {

		case common.FramePing:
			pong := common.Frame{
				Type:     common.FramePong,
				StreamID: 0,
			}
			conn.Write(common.EncodeFrame(pong))

		case common.FrameStreamOpen:
			streams[frame.StreamID] = true
			log.Println("Stream opened:", frame.StreamID)

		case common.FrameStreamClose:
			delete(streams, frame.StreamID)
			log.Println("Stream closed:", frame.StreamID)

		case common.FrameData:
			if !streams[frame.StreamID] {
				log.Println("DATA for unopened stream:", frame.StreamID)
				continue
			}

			log.Printf("Received DATA (Stream %d): %s",
				frame.StreamID,
				string(frame.Payload),
			)
		}
	}
}

