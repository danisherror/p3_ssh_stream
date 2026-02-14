package main

import (
	"log"
	"net"
	"sync"
	"time"

	"p3_ssh_stream/common"
)

type ConnManager struct {
	addr string

	state common.ConnState
	conn  net.Conn

	mu sync.RWMutex

	sendCh chan []byte

	lastPong time.Time
 	streams map[uint32]*Stream
	nextID  uint32
}

func NewConnManager(addr string) *ConnManager {
	return &ConnManager{
		addr:     addr,
		state:    common.StateDisconnected,
		sendCh:   make(chan []byte, 100),
		lastPong: time.Now(),
		streams:  make(map[uint32]*Stream),
		nextID:   1,
	}
}

func (cm *ConnManager) Start() {
	go cm.connectionLoop()
	go cm.writerLoop()
	go cm.heartbeatLoop()
	go cm.healthMonitorLoop()
}

func (cm *ConnManager) connectionLoop() {
	for {
		cm.setState(common.StateConnecting)
		log.Println("Connecting to", cm.addr)

		conn, err := net.Dial("tcp", cm.addr)
		if err != nil {
			log.Println("Connect failed:", err)
			cm.setState(common.StateDisconnected)
			time.Sleep(2 * time.Second)
			continue
		}

		cm.mu.Lock()
		cm.conn = conn
		cm.lastPong = time.Now()
		cm.mu.Unlock()

		cm.setState(common.StateConnected)
		log.Println("Connected")

		cm.readLoop(conn)

		cm.cleanupConnection()
		time.Sleep(2 * time.Second)
	}
}

func (cm *ConnManager) readLoop(conn net.Conn) {
	for {
		frame, err := common.DecodeFrame(conn)
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		switch frame.Type {
		case common.FramePong:
			cm.mu.Lock()
			cm.lastPong = time.Now()
			cm.mu.Unlock()

		case common.FrameData:
			cm.mu.RLock()
			stream := cm.streams[frame.StreamID]
			cm.mu.RUnlock()

			if stream != nil {
				stream.handleIncoming(frame.Payload)
			} else {
				log.Println("Unknown stream:", frame.StreamID)
			}


		}
	}
}

func (cm *ConnManager) cleanupConnection() {
	cm.mu.Lock()
	if cm.conn != nil {
		cm.conn.Close()
		cm.conn = nil
	}
	cm.mu.Unlock()

	cm.setState(common.StateDisconnected)
	log.Println("Disconnected")
}

func (cm *ConnManager) writerLoop() {
	for msg := range cm.sendCh {
		cm.mu.RLock()
		conn := cm.conn
		cm.mu.RUnlock()

		if conn == nil {
			continue
		}

		_, err := conn.Write(msg)
		if err != nil {
			log.Println("Write error:", err)
			cm.setState(common.StateDegraded)
		}
	}
}

func (cm *ConnManager) heartbeatLoop() {
	ticker := time.NewTicker(3 * time.Second)

	for range ticker.C {
		frame := common.Frame{
	Type:     common.FramePing,
	StreamID: 0, // Control stream
	Payload:  nil,
}
		cm.Send(common.EncodeFrame(frame))
	}
}

func (cm *ConnManager) healthMonitorLoop() {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		cm.mu.RLock()
		state := cm.state
		last := cm.lastPong
		cm.mu.RUnlock()

		if state == common.StateConnected {
			if time.Since(last) > 6*time.Second {
				cm.setState(common.StateDegraded)
			}
		}

		if state == common.StateDegraded {
			if time.Since(last) <= 6*time.Second {
				cm.setState(common.StateConnected)
			}
		}
	}
}

func (cm *ConnManager) Send(data []byte) {
	cm.mu.RLock()
	state := cm.state
	cm.mu.RUnlock()

	if state != common.StateConnected {
		log.Println("Dropping frame, not connected")
		return
	}

	select {
	case cm.sendCh <- data:
	default:
		log.Println("Send buffer full, dropping message")
	}
}

func (cm *ConnManager) setState(s common.ConnState) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.state != s {
		log.Println("State:", s.String())
		cm.state = s
	}
}
func (cm *ConnManager) IsConnected() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.state == common.StateConnected
}

func (cm *ConnManager) CreateStream() *Stream {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	id := cm.nextID
	cm.nextID++

	stream := NewStream(id, cm)
	cm.streams[id] = stream

	log.Println("Created stream", id)

	return stream
}

