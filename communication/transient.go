package communication

import (
	"log"
	"net"
)

func NewTransient(addr string) net.PacketConn {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalf("Failed to bind to address %s: %v", addr, err)
	}
	return conn
}
