package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"teng0ni-go-mc/packet"
	"teng0ni-go-mc/packet_map/status_map"
	"teng0ni-go-mc/protocol"
	"teng0ni-go-mc/state"
)

type Client struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
	Conn net.Conn
	State state.State
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Connected successfully\n")

	return &Client {
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
		Conn: conn,
	}, nil
}

func (client *Client) Write_Packet(data []byte) error {
	length := protocol.ITOVI(uint64(len(data)))

	packet := make([]byte, 0)
	packet = append(packet, length...)
	packet = append(packet, data...)

	_, err := client.Conn.Write(packet)
	return err
}

func (client *Client) Read_Packet() (*packet.Packet, error) {
	fmt.Printf("Received a packet!\n")
	length, err := protocol.Read_Varint(client.Reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}

	data_and_id := make([]byte, length)
	if _, err := io.ReadFull(client.Reader, data_and_id); err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}

	reader := bufio.NewReader(bytes.NewReader(data_and_id))

	id, err := protocol.Read_Varint(reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}
	fmt.Printf("Packet ID: %v\n", id)
	fmt.Printf("Packet Name: %v\n", status_map.Map_Clientbound(id))

	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}

	return &packet.Packet{
		ID: id,
		Data: data,
		DataLength: uint64(len(data)),
	}, nil
}

func (client *Client) Decode_Packet(packet *packet.Packet) (any, error) {
	switch client.State {
	case state.Status:
		return status_map.Decode(packet)
	}
	return nil, nil
}
