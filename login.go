package main

import(
	"bytes"
	"fmt"
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"bufio"
	"github.com/maxsupermanhd/go-mc-ms-auth"
	"log"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
)

const (
	HOST = "one.mods.ovh"
	PORT = 55081
	PROTOCOL_V = 777
	USERNAME = "TimppaEZ"
)

func parse_public_key(data []byte) (*rsa.PublicKey, error) {
    parsed, err := x509.ParsePKIXPublicKey(data)
    if err != nil {
        return nil, err
    }

    publicKey, ok := parsed.(*rsa.PublicKey)
    if !ok {
        return nil, errors.New("server public key is not RSA")
    }

    return publicKey, nil
}

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

func read_boolean(reader *bufio.Reader) (bool, error) {
	value := make([]byte, 1)

	_, err := io.ReadFull(reader, value)
	if err != nil {
		return false, err
	}

	return value[0] != 0, nil
}

func write_string(buf *bytes.Buffer, value string) {
	write_varint(buf, uint64(len(value)))
	buf.WriteString(value)
}

func read_string(reader *bufio.Reader) (string, error) {
	length, err := read_varint(reader)
	if err != nil {
		return "", err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)

	return string(data), err
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

func write_pfarray(buf *bytes.Buffer, data []byte) {
	write_varint(buf, uint64(len(data)))
	buf.Write(data)
}

func read_pfarray(reader *bufio.Reader) ([]byte, error) {
	length, err := read_varint(reader)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	return data, err
}

func auth() (GMMAuth.BotAuth, error){
	mauth, err := GMMAuth.GetMCcredentials("./auth.cache", "88650e7e-efee-4857-b9a9-cf580a00ef43")
	if err != nil {
    	log.Print(err)
    	return mauth, err
	}
	return mauth, nil
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
	write_varint(handshake_buf, 2)

	if err := write_packet(conn, handshake_buf.Bytes()); err != nil {
		fmt.Printf("Error: %v", err)
	}

	login_buf := new(bytes.Buffer)
	write_varint(login_buf, 0x0)
	write_string(login_buf, USERNAME)
	login_buf.Write(make([]byte, 16))

	if err := write_packet(conn, login_buf.Bytes()); err != nil {
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

	server_ID, err := read_string(packet_reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	public_key, err := read_pfarray(packet_reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	verify_token, err := read_pfarray(packet_reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	should_authenticate, err := read_boolean(packet_reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("Packet ID: %v\n", packet_ID)
	fmt.Printf("Server ID: %s\n", server_ID)
	fmt.Printf("Public Key: %x\n", public_key)
	fmt.Printf("Verify Token: %x\n", verify_token)
	fmt.Printf("Should Authenticate: %v\n", should_authenticate)

	mauth, err := auth()

	fmt.Printf("Auth: %v", mauth)

	parsed_key, err := parse_public_key(public_key)
	if err != nil {
    	log.Fatal(err)
	}

	shared_secret := make([]byte, 16)
	if _, err := rand.Read(shared_secret); err != nil {
		fmt.Printf("Error: %v", err)
	}

	encrypted_secret, err := rsa.EncryptPKCS1v15(
		rand.Reader,
		parsed_key,
		shared_secret,
		)

	encrypted_token, err := rsa.EncryptPKCS1v15(
		rand.Reader,
		parsed_key,
		verify_token,
		)

	login_buf_2 := new(bytes.Buffer)
	write_varint(login_buf_2, 0x0)
	write_pfarray(login_buf_2, encrypted_secret)
	write_pfarray(login_buf_2, encrypted_token)

	if err := write_packet(conn, login_buf_2.Bytes()); err != nil {
		fmt.Printf("Error: %v", err)
	}


	fmt.Printf("Packet ID: %v\n", packet_ID)
	fmt.Printf("Server ID: %s\n", server_ID)
	fmt.Printf("Public Key: %x\n", public_key)
	fmt.Printf("Verify Token: %x\n", verify_token)
	fmt.Printf("Should Authenticate: %v\n", should_authenticate)
}
