package main

import(
	"fmt"
	"net"
	"strconv"
	"teng0ni-go-mc/packet_map/status_map"
	"teng0ni-go-mc/client"
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
	defer client.Conn.Close()

	handshake_buf := make([]byte, 0) 
	handshake_buf = append(handshake_buf, protocol.ITOVI(0)...)
	handshake_buf = append(handshake_buf, protocol.ITOVI(PROTOCOL_V)...)
	handshake_buf = append(handshake_buf, protocol.Encode_String(HOST)...)
	port := uint16(PORT)
	handshake_buf = append(handshake_buf, byte(port>>8), byte(port&0xFF))
	handshake_buf = append(handshake_buf, protocol.ITOVI(1)...)

	if err := client.Write_Packet(handshake_buf); err != nil {
		fmt.Printf("Error: %v", err)
	}

	status_buf := make([]byte, 0)
	status_buf = append(status_buf, protocol.ITOVI(0)...)

	if err := client.Write_Packet(status_buf); err != nil {
		fmt.Printf("Error: %v", err)
	}

	packet, err := client.Read_Packet()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	thisdata, err := status_map.Decode(packet)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("Response: %v", thisdata)
}
