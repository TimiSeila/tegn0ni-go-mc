package status_map

// - Serverbound
type StatusRequest struct {}

type PingRequest struct {
	Timestamp int64
}

// - Clientbound
type StatusResponse struct {
	JSONResponse string
}

type PongResponse struct {
	Timestamp int64
}
