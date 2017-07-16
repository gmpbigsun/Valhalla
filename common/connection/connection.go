package connection

import (
	"fmt"

	"github.com/Hucaru/Valhalla/common/packet"
)

// Connection -
type Connection interface {
	Write(p packet.Packet) error
	Read(p packet.Packet) error
	Close()
	String() string
}

// PacketHandler -
type PacketHandler func(conn Connection, p packet.Packet, isHeader bool) int

// HandleNewConnection -
func HandleNewConnection(conn Connection, handler PacketHandler, sizeOfRead int) {
	sizeToRead := sizeOfRead
	isHeader := true
	fmt.Println("New connection from", conn)
	for {
		buffer := make([]byte, sizeToRead)

		err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Error in reading from", conn, ", closing the connection")
			conn.Close()
			return
		}

		p := packet.NewPacket()
		p.Append(buffer)
		sizeToRead = handler(conn, p, isHeader)
		isHeader = !isHeader
	}
}