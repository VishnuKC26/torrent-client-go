package client

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/VishnuKC26/torrent-client-go/bitfield"
	"github.com/VishnuKC26/torrent-client-go/handshake"
	"github.com/VishnuKC26/torrent-client-go/message"
	"github.com/VishnuKC26/torrent-client-go/peer"
)

type Client struct {
	Conn     net.Conn
	Bitfield bitfield.Bitfield
	Choked   bool
}

func New(peer peer.Peer, peerID, infoHash [20]byte) (*Client, error) {

	addr := net.JoinHostPort(peer.IP.String(), fmt.Sprintf("%d", peer.Port))

	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {

		return nil, err
	}

	h := handshake.Handshake{
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	_, err = conn.Write(h.Serialize())

	if err != nil {
		conn.Close()
		return nil, err
	}

	resp, err := handshake.Read((conn))

	if err != nil {
		conn.Close()
		return nil, err
	}

	if resp.InfoHash != infoHash {
		conn.Close()
		return nil, fmt.Errorf("infohash mismatch")
	}

	c := &Client{
		Conn:   conn,
		Choked: true,
	}

	msg, err := message.Read(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if msg != nil && msg.ID == message.MsgBitfield {
		c.Bitfield = bitfield.Bitfield(msg.Payload)
	} else {
		// No bitfield received → optimistic assumption
		c.Bitfield = nil
	}

	return c, nil

}

func (c *Client) Read() (*message.Message, error) {
	return message.Read(c.Conn)
}

func (c *Client) SendInterested() error {
	msg := &message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendUnchoke() error {
	msg := &message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendHave(index int) error {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))

	msg := &message.Message{
		ID:      message.MsgHave,
		Payload: payload,
	}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendRequest(index, begin, length int) error {

	log.Printf(
		"SEND REQUEST → piece=%d begin=%d length=%d",
		index, begin, length,
	)
	payload := make([]byte, 12)

	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))

	msg := &message.Message{
		ID:      message.MsgRequest,
		Payload: payload,
	}

	_, err := c.Conn.Write(msg.Serialize())
	return err
}
