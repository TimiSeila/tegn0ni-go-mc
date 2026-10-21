package main

import(
	"fmt"
	"net"
	"strconv"
	"teng0ni-go-mc/client"
	"teng0ni-go-mc/packet"
	"teng0ni-go-mc/protocol"
)

const (
	HOST = "one.mods.ovh"
	PORT = 55081
	PROTOCOL_V = 777
	USERNAME = "TimppaEZ"
)

func main() {
	address := net.JoinHostPort(HOST, strconv.Itoa(PORT))
	client, err := client.New(address)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	defer client.Connection.Close()

	handshake_buf := make([]byte, 0) 
	handshake_buf = append(handshake_buf, protocol.ITOVI(PROTOCOL_V)...)
	handshake_buf = append(handshake_buf, protocol.Encode_String(HOST)...)
	port := uint16(PORT)
	handshake_buf = append(handshake_buf, byte(port>>8), byte(port&0xFF))
	handshake_buf = append(handshake_buf, protocol.ITOVI(1)...)
	handshake_packet := &packet.Packet{ID: 0, Data: handshake_buf}

	if err := client.Connection.Write_Packet(handshake_packet); err != nil {
		fmt.Printf("Error: %v", err)
	}

	status_packet := &packet.Packet{ID: 0, Data: make([]byte, 0)}

	if err := client.Connection.Write_Packet(status_packet); err != nil {
		fmt.Printf("Error: %v", err)
	}

	response_packet := &packet.Packet{}
	if err := client.Connection.Read_Packet(response_packet); err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("ID: %X\n", response_packet.ID)
	fmt.Printf("Data: %v\n", response_packet.Data)
}
