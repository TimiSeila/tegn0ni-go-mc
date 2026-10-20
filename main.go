package main

import(
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"bufio"
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

func write_string(buf *bytes.Buffer, value string) {
	protocol.Write_Varint(bufio.NewWriter(buf), uint64(len(value)))
	buf.WriteString(value)
}

func write_packet(conn net.Conn, data []byte) error {
	length := protocol.ITOVI(uint64(len(data)))

	packet := make([]byte, 0, len(length)+len(data))
    packet = append(packet, length...)
    packet = append(packet, data...)

	_, err := conn.Write(packet)
	return err
}

func read_packet(reader *bufio.Reader) ([]byte, error) {
	length, err := protocol.Read_Varint(bufio.NewReader(reader))
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	return data, err
}

func main() {
	address := net.JoinHostPort(HOST, strconv.Itoa(PORT))
	client, err := client.New(address)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	byte_reader := bufio.NewReader(client.Conn)

	defer client.Conn.Close()

	handshake_buf := make([]byte, 0) 
	handshake_buf = append(handshake_buf, protocol.ITOVI(0)...)
	handshake_buf = append(handshake_buf, protocol.ITOVI(PROTOCOL_V)...)
	handshake_buf = append(handshake_buf, protocol.Encode_String(HOST)...)
	port := uint16(PORT)
	handshake_buf = append(handshake_buf, byte(port>>8), byte(port&0xFF))
	handshake_buf = append(handshake_buf, protocol.ITOVI(1)...)

	if err := write_packet(client.Conn, handshake_buf); err != nil {
		fmt.Printf("Error: %v", err)
	}

	status_buf := make([]byte, 0)
	status_buf = append(status_buf, protocol.ITOVI(0)...)

	if err := write_packet(client.Conn, status_buf); err != nil {
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

	

	packet1, err := read_packet(byte_reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	packet_reader := bufio.NewReader(bytes.NewReader(packet1))
	packet_ID, err := protocol.Read_Varint(packet_reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	json_length, err := protocol.Read_Varint(packet_reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	json_data := make([]byte, json_length)
	if _, err := io.ReadFull(packet_reader, json_data); err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("Packet ID: %v\n", packet_ID)
	fmt.Printf("Packet Name: %s\n", status_map.Map_Clientbound(packet_ID))

	fmt.Printf("Data: %s\n", json_data)
}
