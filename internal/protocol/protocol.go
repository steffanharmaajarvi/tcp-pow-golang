package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

// Header of TCP-message in protocol, means type of message
const (
	Close             = iota // on quit each side (server or client) should close connection
	RequestChallenge         // from client to server - request new challenge from server
	ResponseChallenge        // from server to client - message with challenge for client
	RequestResource          // from client to server - message with solved challenge
	ResponseResource         // from server to client - message with useful info is solution is correct, or with error if not
)

type Message struct {
	Header  int // message type
	Payload string
}

const MESSAGE_SEPARATOR = "|"

func (m *Message) ToString() string {
	return fmt.Sprintf("%d%s%s", m.Header, MESSAGE_SEPARATOR, m.Payload)
}

// ParseMessage - parses Message from str, checks header and payload
func ParseMessage(str string) (*Message, error) {
	str = strings.TrimSpace(str)
	var msgType int

	fmt.Sprintln("parsing message:", str)
	parts := strings.Split(str, MESSAGE_SEPARATOR)

	if len(parts) < 1 || len(parts) > 2 {
		return nil, fmt.Errorf("message doesn't match protocol")
	}
	// try to parse header
	msgType, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("cannot parse header: " + str)
	}
	msg := Message{
		Header: msgType,
	}

	if len(parts) == 2 {
		msg.Payload = parts[1]
	}

	return &msg, nil
}
