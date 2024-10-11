package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"edgefusion-video-push/service"
)

type Listener struct {
}

func NewLister() *Listener {
	log.Printf("数据接收者实例化")
	return &Listener{}
}

func (l *Listener) Lister(queue *service.Queue, f *Forward) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 65505})
	if err != nil {
		log.Fatalf("Failed to bind to address %v", err)
	}
	buf := make([]byte, 1500)
	log.Println("数据接收启动。。。。。")
	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("Failed to read from UDP: %v", err)
			continue
		}
		// 添加到队列中
		queue.Put()
		f.Live(buf[n:])
		f.Transmit(buf[n:])
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

type Forward struct {
	conn                            *net.UDPConn
	transmitAddr, localTransmitAddr *net.UDPAddr
	TransmitAddrIp                  string
}

func NewForward() *Forward {
	transmit, err := NewTransmit("127.0.0.1:65506")
	if err != nil {
		fmt.Println("Error creating connection:", err)
		return nil
	}
	// 创建一个 UDP 连接
	conn, err := net.ListenUDP("udp", transmit)
	if err != nil {
		fmt.Println("Error creating connection:", err)
		return nil
	}
	return &Forward{
		conn: conn,
	}
}

func (f *Forward) Transmit(data []byte) {
	if f.transmitAddr != nil {
		if _, err := f.conn.WriteTo(data, f.transmitAddr); err != nil {
			fmt.Println("Error sending UDP packet:", err)
			return
		}
	}
}

func (f *Forward) Live(data []byte) {
	if f.localTransmitAddr != nil {
		if _, err := f.conn.WriteTo(data, f.localTransmitAddr); err != nil {
			fmt.Println("Error sending UDP packet:", err)
			return
		}
	}
}

func (f *Forward) SetTransmitAddr(transmitAddr *net.UDPAddr, typ int) {
	if transmitAddr == nil {
		return
	}
	if typ == 1 {
		f.transmitAddr = transmitAddr
		f.TransmitAddrIp = strings.Split(transmitAddr.String(), ":")[0]
	} else if typ == 2 {
		f.localTransmitAddr = transmitAddr
	}
}
