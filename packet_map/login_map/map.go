package loginmap

Clientbound_packet_mappings := map[uint64]string(
	0: "disconnect",
	1: "encryption-request",
	2: "login-success",
	3: "set-compression",
	4: "login-plugin-reset",
	5: "cookie-request",
)

Serverbound_packet_mappings := map[uint64]string(
	0: "login-start",
	1: "encryption-response",
	2: "login-plugin-response",
	3: "login-acknowledged",
	4: "cookie-response",
)
