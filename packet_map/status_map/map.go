package status_map

var clientbound_packet_mappings = map[uint64]string{
	0: "status-response",
	1: "pong-response",
}

var serverbound_packet_mappings = map[uint64]string{
	0: "status-request",
	1: "ping-request",
}

func Map_Clientbound(id uint64) string {
	packet_name := clientbound_packet_mappings[id]
	return packet_name
}

func Map_Serverbound(id uint64) string {
	packet_name := serverbound_packet_mappings[id]
	return packet_name
}
