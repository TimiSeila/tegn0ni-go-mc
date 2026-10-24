package protocol

import(
	"errors"
	"fmt"
	"io"
)

func Encode_Varint(value int) []byte {
    v := int32(value)
    var result []byte

    for {
        if v&^int32(0x7F) == 0 {
            result = append(result, byte(v))
            return result
        }

        result = append(result, byte(v&0x7F)|0x80)
        v = int32(uint32(v) >> 7)
    }
}

func Encode_String(value string) []byte {
	string_bytes := []byte(value)
	string_length := Encode_Varint(len(string_bytes))

	buffer := make([]byte, 0)
	buffer = append(buffer, string_length...)
	buffer = append(buffer, string_bytes...)

    return buffer
}

func Read_String(reader io.Reader) (string, error) {
	string_length, err := Read_Varint(reader)
	if err != nil {
		return "", err
	}

	value := make([]byte, string_length)
	if _, err := io.ReadFull(reader, value); err != nil {
		return "", err
	}

	return string(value), nil
}

func Read_Boolean(reader io.Reader) (bool, error) {
	byte := make([]byte, 1)
	if _, err := io.ReadFull(reader, byte); err != nil {
		return false, err
	}

	return byte[0] != 0, nil
}

func Read_PFArray_Byte(reader io.Reader) ([]byte, error) {
	length, err := Read_Varint(reader)
	if err != nil {
		return nil, err
	}

	arr := make([]byte, length)
	if _, err := io.ReadFull(reader, arr); err != nil {
		return nil, err
	}

	return arr, nil
}

func Read_Varint(reader io.Reader) (int, error) {
    var result int32
    var buf [1]byte

    for i := 0; i < 5; i++ {
        if _, err := io.ReadFull(reader, buf[:]); err != nil {
			return 0, fmt.Errorf("reading VarInt byte %d: %w", i+1, err)
        }

        b := buf[0]
        result |= int32(b&0x7F) << (7 * i)

        if b&0x80 == 0 {
            return int(result), nil
        }
    }

	return 0, errors.New("Error: minecraft VarInt exceeds 5 bytes")
}
