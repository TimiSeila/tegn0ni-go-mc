package packet

import(
	"bytes"
	"fmt"
	"teng0ni-go-mc/protocol"
)

type Packet struct {
	ID int
	Data []byte
}

func (packet *Packet) Unpack(packet_length int, raw_data []byte) error {
	reader := bytes.NewReader(raw_data)

	packet_id, err := protocol.Read_Varint(reader)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}

	packet.ID = packet_id
	packet.Data = raw_data
	return nil
}

func (packet *Packet) Pack() ([]byte, error) {
	raw_data := make([]byte, 0)

	id := protocol.Encode_Varint(packet.ID)
	length := protocol.Encode_Varint(len(id) + len(packet.Data))

	raw_data = append(raw_data, length...)
	raw_data = append(raw_data, id...)
	raw_data = append(raw_data, packet.Data...)

	return raw_data, nil
}
