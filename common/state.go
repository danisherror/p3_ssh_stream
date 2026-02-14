package common

type ConnState int

const (
	StateDisconnected ConnState = iota
	StateConnecting
	StateConnected
	StateDegraded
)

func (s ConnState) String() string {
	switch s {
	case StateDisconnected:
		return "DISCONNECTED"
	case StateConnecting:
		return "CONNECTING"
	case StateConnected:
		return "CONNECTED"
	case StateDegraded:
		return "DEGRADED"
	default:
		return "UNKNOWN"
	}
}

