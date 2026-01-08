// package handshake

// import "io"

// // this package builds and reads the handshake bytes

// //BitTorrent handshake format

// // <length=19><"BitTorrent protocol"><8 reserved bytes><info_hash><peer_id>

// // InfoHash is the torrent identity
// type Handshake struct {
// 	InfoHash [20]byte
// 	PeerID   [20]byte
// }

// func (h *Handshake) Serialize() []byte {

// 	buf := make([]byte, 49+19)

// 	buf[0] = 19

// 	copy(buf[1:], []byte("BitTorrent protocol"))

// 	copy(buf[28:], h.InfoHash[:])
// 	copy(buf[48:], h.PeerID[:])

// 	return buf
// }

// // r is something that can read bytes in my buffer and return number of bytes read

// // it implements the Reader interface in the io package

// // This interface directs to have functions like ReadFull which

// // io.Reader is like math.Sqrt or fmt.Println

// // so r io.Reader means r can hold any value whose type has a Read([]byte) method

// // r is a variable and its type is an interface

// // in golang when a variable has an interface type, it does not store data itself

// // instead it stores a concrete value like *os.File, *bytes.Buffer etc

// // r = os.stdin means whatever r points to has a Read([]byte) method

// // so in this case: in the implementation of ReadFull func in io package we have nr, err = r.Read(buf)

// // and this r can point to any type implementing the Reader interface

// // in runtime we do handshake.Read(conn) where conn is of type *net.TCPConn which implements io.Reader

// // so this conn is passed as argument to the Read function which has parameter r io.Reader in now r is conn and conn is of type *TCPConn and net.*TCPConn implements the Read([]byte) interface thus the contract is satisfied as r is pointing to a type which implements Read(p []byte) method

// // so this internally calls conn.Read(...)

// func Read(r io.Reader) (*Handshake, error) {

// 	lengthBuf := make([]byte, 1)
// 	if _, err := io.ReadFull(r, lengthBuf); err != nil {
// 		return nil, err
// 	}

// 	pstrlen := int(lengthBuf[0])
// 	if pstrlen == 0 {
// 		return nil, io.ErrUnexpectedEOF
// 	}

// 	handshakeBuf :=
// 		make([]byte, pstrlen+48)

// 	if _, err := io.ReadFull(r, handshakeBuf); err != nil {
// 		return nil, err
// 	}

// 	var infoHash [20]byte
// 	var peerID [20]byte

// 	copy(infoHash[:], handshakeBuf[pstrlen+8:pstrlen+28])
// 	copy(peerID[:], handshakeBuf[pstrlen+28:])

// 	return &Handshake{
// 		InfoHash : infoHash,
// 		PeerID: peerID,
// 	},nil

// }

package handshake

import (
	"fmt"
	"io"
)

const (
	protocolName = "BitTorrent protocol"
	protocolLen  = byte(len(protocolName))
)

// BitTorrent handshake format:
//
// <pstrlen=19><pstr="BitTorrent protocol">
// <8 reserved bytes><info_hash><peer_id>
//
// Total length = 49 + pstrlen

type Handshake struct {
	InfoHash [20]byte
	PeerID   [20]byte
}

/* ---------- Serialize ---------- */

func (h *Handshake) Serialize() []byte {

	buf := make([]byte, 49+len(protocolName))

	// pstrlen
	buf[0] = protocolLen

	// pstr
	copy(buf[1:], protocolName)

	// reserved bytes (already zeroed by make)

	// info_hash
	copy(buf[1+protocolLen+8:], h.InfoHash[:])

	// peer_id
	copy(buf[1+protocolLen+8+20:], h.PeerID[:])

	return buf
}

/* ---------- Read ---------- */

// Read reads a BitTorrent handshake from any io.Reader.
//
// r can be:
//   - *net.TCPConn
//   - *bytes.Reader
//   - *os.File
//
// io.ReadFull is used because handshake framing is strict:
// partial reads would corrupt the protocol.
func Read(r io.Reader) (*Handshake, error) {

	// Read pstrlen
	lengthBuf := make([]byte, 1)
	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return nil, err
	}

	pstrlen := int(lengthBuf[0])
	if pstrlen == 0 {
		return nil, io.ErrUnexpectedEOF
	}

	// Read rest of handshake
	handshakeBuf := make([]byte, pstrlen+48)
	if _, err := io.ReadFull(r, handshakeBuf); err != nil {
		return nil, err
	}

	// Validate protocol string
	pstr := string(handshakeBuf[:pstrlen])
	if pstr != protocolName {
		return nil, fmt.Errorf("unexpected protocol string: %q", pstr)
	}

	// Offsets
	infoHashOffset := pstrlen + 8
	peerIDOffset := infoHashOffset + 20

	if len(handshakeBuf) < peerIDOffset+20 {
		return nil, io.ErrUnexpectedEOF
	}

	var infoHash [20]byte
	var peerID [20]byte

	copy(infoHash[:], handshakeBuf[infoHashOffset:infoHashOffset+20])
	copy(peerID[:], handshakeBuf[peerIDOffset:peerIDOffset+20])

	return &Handshake{
		InfoHash: infoHash,
		PeerID:   peerID,
	}, nil
}
