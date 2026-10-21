package connection

import(
	"fmt"
	"io"
	"net"
	"teng0ni-go-mc/packet"
	"teng0ni-go-mc/protocol"
)

type Connection struct {
	net.Conn
	State int
}

func (connection *Connection) Read_Packet(packet *packet.Packet) error {
	packet_length, err := protocol.Read_Varint(connection)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}

	raw_data := make([]byte, packet_length)
	_, err = io.ReadFull(connection, raw_data)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}

	err = packet.Unpack(packet_length, raw_data)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}
	return nil
}

func (connection *Connection) Write_Packet(packet *packet.Packet) error {
	packed_data, err := packet.Pack()
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}

	_, err = connection.Write(packed_data)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}

	return nil
}
