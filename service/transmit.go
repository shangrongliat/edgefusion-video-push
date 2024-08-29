package service

import (
	"fmt"
	"log"
	"net"
)

type Listener struct {
}

func NewLister() *Listener {
	log.Printf("数据接收者实例化")
	return &Listener{}
}

func (l *Listener) Lister(queue *Queue, f *Forward) {
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
		f.Send(buf[12:n])
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
	listener                        net.Listener
	conn                            net.Conn
	transmitAddr, localTransmitAddr *net.UDPAddr
}

func NewForward() *Forward {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening:", err)
	}

	return &Forward{
		listener: listener,
	}
}

func (f *Forward) Start() {
	for {
		// 接受客户端连接
		conn, err := f.listener.Accept()
		if err != nil {
			log.Fatal("Error accepting:", err)
		}
		f.conn = conn
	}
}

func (f *Forward) Send(data []byte) {
	// 创建一个 bufio.Reader 用于读取客户端数据
	// 读取客户端发送的数据
	// 向客户端发送响应消息
	_, err := f.conn.Write(data)
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}
}

func (f *Forward) SetTransmitAddr(transmitAddr *net.UDPAddr, typ int) {
	if transmitAddr == nil {
		return
	}
	if typ == 1 {
		f.transmitAddr = transmitAddr
	} else if typ == 2 {
		f.localTransmitAddr = transmitAddr
	}
}
