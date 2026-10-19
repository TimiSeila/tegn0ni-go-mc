package main

import(
	"bytes"
	"fmt"
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"bufio"
)

const (
	HOST = "one.mods.ovh"
	PORT = 55081
	PROTOCOL_V = 777
	USERNAME = timi.koodaa@gmail.com
)

func write_varint(writer io.Writer, value uint64) {
	buf := make([]byte, 0)
	buf = binary.AppendUvarint(buf, value)
	writer.Write(buf)
}

func read_varint(reader *bufio.Reader) (uint64, error) {
	value, err := binary.ReadUvarint(reader)

	if err != nil {
		return 0, err
	}

	return value, nil
}

func write_string(buf *bytes.Buffer, value string) {
	write_varint(buf, uint64(len(value)))
	buf.WriteString(value)
}

func write_packet(conn net.Conn, data []byte) error {
	packet := new(bytes.Buffer)
	write_varint(packet, uint64(len(data)))
	packet.Write(data)

	_, err := conn.Write(packet.Bytes())
	return err
}

func read_packet(reader *bufio.Reader) ([]byte, error) {
	length, err := read_varint(reader)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	return data, err
}

func main() {
	address := net.JoinHostPort(HOST, strconv.Itoa(PORT))
	conn, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	byte_reader := bufio.NewReader(conn)

	defer conn.Close()

	handshake_buf := new(bytes.Buffer)
	write_varint(handshake_buf, 0x0)
	write_varint(handshake_buf, PROTOCOL_V)
	write_string(handshake_buf, HOST)
	binary.Write(handshake_buf, binary.BigEndian, uint16(PORT))
	write_varint(handshake_buf, 1)

	if err := write_packet(conn, handshake_buf.Bytes()); err != nil {
		fmt.Printf("Error: %v", err)
	}

	status_buf := new(bytes.Buffer)
	write_varint(status_buf, 0x0)

	if err := write_packet(conn, status_buf.Bytes()); err != nil {
		fmt.Printf("Error: %v", err)
	}

	packet, err := read_packet(byte_reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	packet_reader := bufio.NewReader(bytes.NewReader(packet))
	packet_ID, err := read_varint(packet_reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	json_length, err := read_varint(packet_reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	json_data := make([]byte, json_length)
	if _, err := io.ReadFull(packet_reader, json_data); err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("Packet ID: %v\n", packet_ID)

	fmt.Printf("Data: %s\n", json_data)
}
