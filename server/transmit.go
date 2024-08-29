package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"edgefusion-video-push/communication"
	"edgefusion-video-push/service"
)

type Listener struct {
	conn net.PacketConn
}

func NewLister() *Listener {
	log.Printf("数据接收者实例化")
	return &Listener{conn: communication.NewTransient("127.0.0.1:65505")}
}

func (l *Listener) Lister(queue *service.Queue) {
	defer func(conn net.PacketConn) {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close from UDP: %v", err)
		}
	}(l.conn)
	file, err := os.Create("/home/edgefusion-video-push/22222222222.h264")
	if err != nil {
		panic(err)
	}
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
		file.Write(buf[12:n])
	}
}

func NewTransmit(remoteAddr string) (*net.UDPAddr, error) {
	remote, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return nil, err
	}
	log.Println("UDP转发客户端初始化", remote)
	return remote, nil
}

func (t *Listener) Transmit(data []byte, transmitAddr *net.UDPAddr) {
	if transmitAddr != nil {
		if _, err := t.conn.WriteTo(data, transmitAddr); err != nil {
			fmt.Println("Error sending UDP packet:", err)
			return
		}
	}
}

func (t *Listener) Live(data []byte, file *os.File, localTransmitAddr *net.UDPAddr) {
	if localTransmitAddr != nil {
		if _, err := file.Write(data[12:]); err != nil {
			fmt.Println("Error sending UDP packet:", err)
		}
	}
}
