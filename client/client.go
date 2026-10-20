package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"teng0ni-go-mc/protocol"
)

type Client struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
	Conn net.Conn
}

type Packet struct {
	ID uint64
	Data []byte
	DataLength uint64
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

func (client *Client) Read_Packet() (*Packet, error) {
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

	reader := bytes.NewReader(data_and_id)

	id, err := protocol.Read_Varint(bufio.NewReader(reader))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}

	return &Packet{
		ID: id,
		Data: data,
		DataLength: uint64(len(data)),
	}, nil
}
