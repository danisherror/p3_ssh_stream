package main

import (
	"log"
	"sync"
	"p3_ssh_stream/common"
)

type Stream struct {
	id uint32
	cm *ConnManager

	recvCh chan []byte
	window int
	mu     sync.Mutex
}

func NewStream(id uint32, cm *ConnManager) *Stream {
	return &Stream{
		id:     id,
		cm:     cm,
		recvCh: make(chan []byte, 100),
		window: 1024, // initial window
	}
}

func (s *Stream) Send(data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(data) > s.window {
		log.Printf("Stream %d blocked (window=%d)", s.id, s.window)
		return
	}

	frame := common.Frame{
		Type:     common.FrameData,
		StreamID: s.id,
		Payload:  data,
	}

	s.window -= len(data)

	s.cm.Send(common.EncodeFrame(frame))
}


func (s *Stream) handleIncoming(data []byte) {
	select {
	case s.recvCh <- data:
	default:
		log.Printf("Stream %d recv buffer full, dropping", s.id)
	}
}

func (s *Stream) Receive() <-chan []byte {
	return s.recvCh
}

