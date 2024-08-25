package server

import (
	"edgefusion-video-push/communication"
	"edgefusion-video-push/service"
	"fmt"
	"log"
	"net"
)

type TransmitInfo struct {
	RemoteAddr *net.UDPAddr
}

type Listener struct {
	conn net.PacketConn
}

func NewLister() *Listener {
	log.Printf("数据接收者实例化")
	return &Listener{conn: communication.NewTransient("127.0.0.1:65515")}
}

func (l *Listener) Lister(queue *service.Queue) {
	defer func(conn net.PacketConn) {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close from UDP: %v", err)
		}
	}(l.conn)
	buf := make([]byte, 1500)
	log.Println("数据接收启动。。。。。")
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
	log.Println("UDP转发客户端初始化", remote)
	return &TransmitInfo{
		RemoteAddr: remote,
	}
}

func (t *Listener) Transmit(data []byte, remoteAddrs ...*net.UDPAddr) {
	for i := range remoteAddrs {
		if _, err := t.conn.WriteTo(data, remoteAddrs[i]); err != nil {
			fmt.Println("Error sending UDP packet:", err)
			return
		}
	}
}
