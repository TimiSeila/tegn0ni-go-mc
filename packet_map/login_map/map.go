package login_map

var clientbound_packet_mappings = map[uint64]string{
	0: "disconnect",
	1: "encryption-request",
	2: "login-success",
	3: "set-compression",
	4: "login-plugin-reset",
	5: "cookie-request",
}

var serverbound_packet_mappings = map[uint64]string{
	0: "login-start",
	1: "encryption-response",
	2: "login-plugin-response",
	3: "login-acknowledged",
	4: "cookie-response",
}

func Map_Clientbound(id uint64) string {
	packet_name := clientbound_packet_mappings[id]
	return packet_name
}

func Map_Serverbound(id uint64) string {
	packet_name := serverbound_packet_mappings[id]
	return packet_name
}
