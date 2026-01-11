package p2p

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/VishnuKC26/torrent-client-go/client"
	"github.com/VishnuKC26/torrent-client-go/message"
	"github.com/VishnuKC26/torrent-client-go/peer"
)

const MaxBlockSize = 16384 // 16KB
const MaxBacklog = 5       //max in-flight requests

type Torrent struct {
	Peers       []peer.Peer
	PeerId      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

type pieceProgress struct {
	index      int
	client     *client.Client
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}

func (state *pieceProgress) readMessage() error {

	msg, err := state.client.Read()
	if err != nil {
		return err
	}

	if msg == nil {
		return nil
	}

	switch msg.ID {

	case message.MsgUnchoke:
		state.client.Choked = false
		log.Printf("unchoked by peer")

	case message.MsgChoke:
		state.client.Choked = true

	case message.MsgHave:
		index, err := message.ParseHave(msg)
		if err != nil {
			return err
		}
		state.client.Bitfield.SetPiece(index)

	case message.MsgPiece:
		n, err := message.ParsePiece(state.index, msg, state.buf)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--

	}

	return nil

}

func attemptDownloadPiece(c *client.Client, pw *pieceWork) ([]byte, error) {

	start := time.Now()

	state := pieceProgress{
		index:  pw.index,
		client: c,
		buf:    make([]byte, pw.length),
	}

	c.Conn.SetDeadline((time.Now().Add(30 * time.Second)))
	defer c.Conn.SetDeadline(time.Time{})

	for state.downloaded < pw.length {
		if time.Since(start) > 30*time.Second {
			return nil, fmt.Errorf("timed out waiting for unchoke")
		}

		if !state.client.Choked {
			for state.backlog < MaxBacklog && state.requested < pw.length {
				blockSize := min(pw.length-state.requested, MaxBacklog)
				err := c.SendRequest(pw.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize

				error := state.readMessage()
				if error != nil {
					return nil, error
				}
			}
		}
	}
	return state.buf, nil
}

func (t *Torrent) startDownloadWorker(peer peer.Peer, workQueue chan *pieceWork, results chan *pieceResult) {

	c, err := client.New(peer, t.PeerId, t.InfoHash)
	if err != nil {
		log.Printf("Could not handshake with %s. Disconnecting\n", peer.IP)
		return
	}
	defer c.Conn.Close()

	log.Printf("completed handshake with %s\n", peer.IP)

	c.SendInterested()

	for pw := range workQueue {

		if c.Bitfield != nil && !c.Bitfield.Has(pw.index) {

			continue
		}

		buf, err := attemptDownloadPiece(c, pw)

		if err != nil {
			log.Println("Worker Failed", err)
			workQueue <- pw
			continue
		}

		err = checkIntegrity(pw, buf)
		if err != nil {
			log.Printf("Piece #%d failed integrity check\n", pw.index)
			workQueue <- pw
			continue
		}

		c.SendHave(pw.index)
		results <- &pieceResult{pw.index, buf}
	}

}

func checkIntegrity(pw *pieceWork, buf []byte) error {

	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], pw.hash[:]) {
		return fmt.Errorf("piece %d failed integrity check", pw.index)
	}
	return nil
}

func (t *Torrent) calculateBoundsForPiece(index int) (begin int, end int) {

	begin = index * t.PieceLength
	end = min(begin+t.PieceLength, t.Length)

	return begin, end
}

func (t *Torrent) calculatePieceSize(index int) int {
	begin, end := t.calculateBoundsForPiece(index)
	return end - begin
}

func (t *Torrent) Download() ([]byte, error) {

	log.Println("starting download for ", t.Name)

	workQueue := make(chan *pieceWork, len(t.PieceHashes))

	results := make(chan *pieceResult)

	for index, hash := range t.PieceHashes {
		length := t.calculatePieceSize((index))
		workQueue <- &pieceWork{index, hash, length}
	}

	maxWorkers := 5
	for i := 0; i < maxWorkers && i < len(t.Peers); i++ {
		go t.startDownloadWorker(t.Peers[i], workQueue, results)
	}

	buf := make([]byte, t.Length)
	donePieces := 0

	for donePieces < len(t.PieceHashes) {
		res := <-results

		begin, end := t.calculateBoundsForPiece((res.index))

		copy(buf[begin:end], res.buf)

		donePieces++

		percent := float64(donePieces) / float64(len(t.PieceHashes)) * 100

		numWorkers := runtime.NumGoroutine() - 1
		log.Printf("%0.2f%% Downloaded piece #%d from %d peers\n", percent, res.index, numWorkers)
	}

	close(workQueue)
	return buf, nil
}
