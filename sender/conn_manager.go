package main

import (
	"log"
	"net"
	"sync"
	"time"

	"p3_ssh_stream/common"
)

type ConnManager struct {
	addr   string
	state  common.ConnState
	conn   net.Conn
	mu     sync.RWMutex
	sendCh chan []byte
}

func NewConnManager(addr string) *ConnManager {
	return &ConnManager{
		addr:   addr,
		state:  common.StateDisconnected,
		sendCh: make(chan []byte, 100),
	}
}

func (cm *ConnManager) Start() {
	go cm.connectionLoop()
	go cm.writerLoop()
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
		cm.mu.Unlock()

		cm.setState(common.StateConnected)
		log.Println("Connected")

		// Wait until connection breaks
		buf := make([]byte, 1)
		_, err = conn.Read(buf)
		if err != nil {
			log.Println("Connection lost:", err)
		}

		cm.mu.Lock()
		cm.conn.Close()
		cm.conn = nil
		cm.mu.Unlock()

		cm.setState(common.StateDisconnected)
		time.Sleep(2 * time.Second)
	}
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

func (cm *ConnManager) Send(data []byte) {
	cm.sendCh <- data
}

func (cm *ConnManager) setState(s common.ConnState) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.state != s {
		log.Println("State:", s.String())
		cm.state = s
	}
}

