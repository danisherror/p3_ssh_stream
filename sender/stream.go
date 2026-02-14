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
	window uint32
	mu     sync.Mutex
	closed bool
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

    if s.closed {
        log.Println("Cannot send, stream closed:", s.id)
        return
    }

    if s.window < uint32(len(data)) {
        log.Printf("Stream %d blocked (window=%d)\n", s.id, s.window)
        return
    }

    // Decrease window BEFORE sending
    s.window -= uint32(len(data))

    frame := common.Frame{
        Type:     common.FrameData,
        StreamID: s.id,
        Payload:  data,
    }

    log.Printf("Sending DATA (Stream %d, %d bytes)\n", s.id, len(data))

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

func (s *Stream) Close() {
    s.mu.Lock()
    defer s.mu.Unlock()

    if s.closed {
        return
    }

    s.closed = true

    frame := common.Frame{
        Type:     common.FrameStreamClose,
        StreamID: s.id,
    }

    log.Println("Sending STREAM_CLOSE:", s.id)

    s.cm.Send(common.EncodeFrame(frame))

    // Optional: remove stream from manager
    s.cm.mu.Lock()
    delete(s.cm.streams, s.id)
    s.cm.mu.Unlock()
}

