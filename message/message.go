package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	MsgChoke         = 0
	MsgUnchoke       = 1
	MsgInterested    = 2
	MsgNotInterested = 3
	MsgHave          = 4
	MsgBitfield      = 5
	MsgRequest       = 6
	MsgPiece         = 7
	MsgCancel        = 8
)

type Message struct {
	ID      uint8
	Payload []byte
}

func (m *Message) Serialize() []byte {

	if m == nil {
		return make([]byte, 4)
	}

	length := uint32(len(m.Payload) + 1)

	buf := make([]byte, 4+length)

	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = m.ID
	copy(buf[5:], m.Payload)

	return buf
}

func Read(r io.Reader) (*Message, error) {

	lengthBuf := make([]byte, 4)

	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)

	if length == 0 {
		return nil, nil
	}

	msgBuf := make([]byte, length)

	if _, err := io.ReadFull(r, msgBuf); err != nil {
		return nil, err
	}

	return &Message{
		ID:      msgBuf[0],
		Payload: msgBuf[1:],
	}, nil
}

func ParseHave(msg *Message) (int, error) {

	if msg.ID != MsgHave {
		return 0, fmt.Errorf("expected HAVE message")
	}

	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("invalid HAVE payload length")
	}

	index := binary.BigEndian.Uint32(msg.Payload)
	return int(index), nil
}

func ParsePiece(index int,msg *Message, buf []byte) (int, error) {
	if msg.ID != MsgPiece {
		return 0, fmt.Errorf("expected PIECE message")
	}

	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("payload too short")
	}

	parsedIndex := binary.BigEndian.Uint32(msg.Payload[0:4])
	if int(parsedIndex) != index {
		return 0,fmt.Errorf("piece index mismatch")
	}
	begin := binary.BigEndian.Uint32(msg.Payload[4:8])

	n := copy(buf[begin:], msg.Payload[8:])
	return n, nil
}

// peer message format
// <length prefix><message id><payload>

// length prefix = 4 bytes
// message id = 1 byte
// payload = variable
