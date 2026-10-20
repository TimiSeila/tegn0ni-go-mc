package state

type State int

const (
	Handshaking State = iota
	Status
	Login
	Configuration
	Play
)
