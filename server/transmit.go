package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"edgefusion-video-push/communication"
	"edgefusion-video-push/service"
)

type TransmitInfo struct {
	remoteAddr *net.UDPAddr
}

type Listener struct {
	conn net.PacketConn
}

func NewLister() Listener {
	return Listener{conn: communication.NewTransient("127.0.0.1:65515")}
}

func (l *Listener) Lister(group *sync.WaitGroup, queue *service.Queue) {
	group.Add(1)
	defer func(conn net.PacketConn) {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close from UDP: %v", err)
		}
	}(l.conn)
	defer group.Done()
	buf := make([]byte, 1500)
	for {
		n, _, err := l.conn.ReadFrom(buf)
		if err != nil {
			log.Printf("Failed to read from UDP: %v", err)
			continue
		}
		// 添加到队列中
		queue.Put(buf[:n])
	}
}

func NewTransmit(remoteAddr string) *TransmitInfo {
	remote, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return nil
	}
	return &TransmitInfo{
		remoteAddr: remote,
	}
}

func (t *Listener) Transmit(data []byte, remoteAddr *net.UDPAddr) {
	if _, err := t.conn.WriteTo(data, remoteAddr); err != nil {
		fmt.Println("Error sending UDP packet:", err)
		return
	}
}
