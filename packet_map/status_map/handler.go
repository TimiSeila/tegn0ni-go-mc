package status_map

import(
	"bytes"
	"encoding/binary"
	"io"
	"teng0ni-go-mc/client"
)

func Encode() {

}

func Decode(packet *client.Packet) (any, error) {
	switch packet.ID {
	case 0:
		reader := bytes.NewReader(packet.Data)

		json_data := make([]byte, packet.DataLength)
		if _, err := io.ReadFull(reader, json_data); err != nil {
            return nil, err
        }
		return StatusResponse{JSONResponse: string(json_data)}, nil

	case 1:
		timestamp := int64(binary.BigEndian.Uint64(packet.Data))
		return PongResponse{Timestamp: timestamp}, nil
	}

	return nil, nil
}
