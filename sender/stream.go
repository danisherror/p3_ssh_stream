package main

import (
	"log"
	"p3_ssh_stream/common"
)

type Stream struct {
	id uint32
	cm *ConnManager

	recvCh chan []byte
}

func NewStream(id uint32, cm *ConnManager) *Stream {
	return &Stream{
		id:     id,
		cm:     cm,
		recvCh: make(chan []byte, 100),
	}
}

func (s *Stream) Send(data []byte) {
	frame := common.Frame{
		Type:     common.FrameData,
		StreamID: s.id,
		Payload:  data,
	}
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

