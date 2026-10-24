package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"teng0ni-go-mc/client"
	"teng0ni-go-mc/msauth"
	"teng0ni-go-mc/packet"
	"teng0ni-go-mc/protocol"
	"time"
	"uuid"
)

type JoinRequest struct {
    AccessToken     string `json:"accessToken"`
    SelectedProfile string `json:"selectedProfile"`
    ServerID        string `json:"serverId"`
}

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

	auth, err := msauth.Authenticate()
	if err != nil {
		fmt.Printf("Auth Error: %v", err)
	}

	fmt.Printf("Auth username: %v\n", auth.Username)
	fmt.Printf("Auth uuid: %v\n", auth.UUID)

	handshake_buf := make([]byte, 0) 
	handshake_buf = append(handshake_buf, protocol.Encode_Varint(PROTOCOL_V)...)
	handshake_buf = append(handshake_buf, protocol.Encode_String(HOST)...)
	port := uint16(PORT)
	handshake_buf = append(handshake_buf, byte(port>>8), byte(port&0xFF))
	handshake_buf = append(handshake_buf, protocol.Encode_Varint(2)...)
	handshake_packet := &packet.Packet{ID: 0, Data: handshake_buf}

	if err := client.Connection.Write_Packet(handshake_packet); err != nil {
		fmt.Printf("Error: %v", err)
	}

	login_start_buf := make([]byte, 0)
	login_start_buf = append(login_start_buf, protocol.Encode_String(auth.Username)...)

	uuid_bytes, err := uuid.Parse(auth.UUID)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	login_start_buf = append(login_start_buf, uuid_bytes[:]...)
	login_start_packet := &packet.Packet{ID: 0, Data: login_start_buf}

	if err := client.Connection.Write_Packet(login_start_packet); err != nil {
		fmt.Printf("Error: %v", err)
	}

	encryption_request_packet := &packet.Packet{}
	if err := client.Connection.Read_Packet(encryption_request_packet); err != nil {
		fmt.Printf("Error: %v", err)
	}

	reader := bytes.NewReader(encryption_request_packet.Data)

	server_id, err := protocol.Read_String(reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Println(server_id)

	public_key, _ := protocol.Read_PFArray_Byte(reader)
	verify_token, _ := protocol.Read_PFArray_Byte(reader)
	should_authenticate, _ := protocol.Read_Boolean(reader)
	fmt.Println(should_authenticate)

	shared_secret := make([]byte, 16)
	if _, err := rand.Read(shared_secret); err != nil {
    	panic(err)
	}

	parsed_key, err := x509.ParsePKIXPublicKey(public_key)
	if err != nil {
    	panic(err)
	}

	rsa_public_key, ok := parsed_key.(*rsa.PublicKey)
	if !ok {
    	panic("server public key is not RSA")
	}

	encrypted_secret, err := rsa.EncryptPKCS1v15(
    	rand.Reader,
    	rsa_public_key,
    	shared_secret,
	)
	if err != nil {
    	panic(err)
	}

	encrypted_verify_token, err := rsa.EncryptPKCS1v15(
    	rand.Reader,
    	rsa_public_key,
    	verify_token,
	)
	if err != nil {
    	panic(err)
	}

	

	server_hash := MinecraftServerHash(server_id, shared_secret, public_key)

	ctx := context.Background()
	httpclient := &http.Client{
		Timeout: 10 * time.Second,
	}

	JoinMinecraftServer(ctx, httpclient, auth.AccessToken, auth.UUID, server_hash)
	encrypt_response_buf := make([]byte, 0)
	encrypt_response_buf = append(encrypt_response_buf, encrypted_secret...)
	encrypt_response_buf = append(encrypt_response_buf, encrypted_verify_token...)
	encrypt_response_packet := &packet.Packet{ID: 1, Data: encrypt_response_buf}

	if err := client.Connection.Write_Packet(encrypt_response_packet); err != nil {
		fmt.Printf("Error: %v", err)
	}

	login_ack_buf := make([]byte, 0)
	login_ack_packet := &packet.Packet{ID: 3, Data: login_ack_buf}

	if err := client.Connection.Write_Packet(login_ack_packet); err != nil {
		fmt.Printf("Error: %v", err)
	}
	for {}
}

func JoinMinecraftServer(
    ctx context.Context,
    client *http.Client,
    accessToken string,
    uuid string,
    serverHash string,
) error {
    payload := JoinRequest{
        AccessToken:     accessToken,
        SelectedProfile: uuid,
        ServerID:        serverHash,
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodPost,
        "https://sessionserver.mojang.com/session/minecraft/join",
        bytes.NewReader(body),
    )
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        responseBody, _ := io.ReadAll(resp.Body)
        return fmt.Errorf(
            "Minecraft session join failed: status %d: %s",
            resp.StatusCode,
            responseBody,
        )
    }

    return nil
}

func MinecraftServerHash(
    serverID string,
    sharedSecret []byte,
    publicKey []byte,
) string {
    h := sha1.New()
    h.Write([]byte(serverID))
    h.Write(sharedSecret)
    h.Write(publicKey)

    hash := h.Sum(nil)

    // Minecraft uses Java BigInteger's signed hexadecimal representation.
    negative := hash[0]&0x80 != 0

    if !negative {
        return strings.TrimLeft(hex.EncodeToString(hash), "0")
    }

    // Compute the magnitude of the negative two's-complement integer.
    magnitude := make([]byte, len(hash))
    copy(magnitude, hash)

    for i := range magnitude {
        magnitude[i] = ^magnitude[i]
    }

    for i := len(magnitude) - 1; i >= 0; i-- {
        magnitude[i]++
        if magnitude[i] != 0 {
            break
        }
    }

    return "-" + strings.TrimLeft(hex.EncodeToString(magnitude), "0")
}
