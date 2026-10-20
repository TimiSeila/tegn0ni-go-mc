package protocol

import(
	"bufio"
	"encoding/binary"
)

func ITOVI(value uint64) []byte {
	buffer := make([]byte, 0)
	buffer = binary.AppendUvarint(buffer, value)
	return buffer
}

func Encode_String(value string) []byte {
	string_bytes := []byte(value)
	string_length := ITOVI(uint64(len(string_bytes)))

	buffer := make([]byte, 0)
	buffer = append(buffer, string_length...)
	buffer = append(buffer, string_bytes...)

    return buffer
}

func Write_Varint(writer *bufio.Writer, value uint64) {
	buffer := make([]byte, 0)
	buffer = binary.AppendUvarint(buffer, value)
	writer.Write(buffer)
}

func Read_Varint(reader *bufio.Reader) (uint64, error) {
	value, err := binary.ReadUvarint(reader)
	if err != nil {
		return 0, err
	}

	return value, nil
}
