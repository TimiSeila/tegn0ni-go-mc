package client

import (
	"fmt"
	"net"
	"teng0ni-go-mc/connection"
	"teng0ni-go-mc/packet"
)

type Client struct {
	Connection *connection.Connection
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Connected successfully\n")

	return &Client {
		Connection: &connection.Connection{
			Conn: conn,
			State: connection.Handshaking,
		},
	}, nil
}

func (client *Client) Read_Loop() error {
	for {
		packet := &packet.Packet{}

		err := client.Connection.Read_Packet(packet)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return err
		}

		fmt.Printf("Received packet with ID: %v", packet.ID)
	}
}
