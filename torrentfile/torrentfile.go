package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

/* ---------- bencode structs ---------- */

type bencodeTorrent struct {
	Announce     string      `bencode:"announce"`
	AnnounceList [][]string  `bencode:"announce-list"`
	Info         bencodeInfo `bencode:"info"`
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

/* ---------- Open ---------- */

func Open(path string) (*TorrentFile, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var bt bencodeTorrent
	if err := bencode.Unmarshal(bytes.NewReader(data), &bt); err != nil {
		return nil, err
	}

	// Choose announce URL
	announce := bt.Announce
	if announce == "" {
		if len(bt.AnnounceList) > 0 && len(bt.AnnounceList[0]) > 0 {
			announce = bt.AnnounceList[0][0]
		}
	}

	if announce == "" {
		return nil, fmt.Errorf("no announce or announce-list found in torrent")
	}

	infoHash, err := bt.Info.infoHash()
	if err != nil {
		return nil, err
	}

	pieceHashes, err := splitPieceHashes(bt.Info.Pieces)
	if err != nil {
		return nil, err
	}

	tf := TorrentFile{
		Announce:    announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bt.Info.PieceLength,
		Length:      bt.Info.Length,
		Name:        bt.Info.Name,
	}

	return &tf, nil
}

/* ---------- helpers ---------- */

func (i *bencodeInfo) infoHash() ([20]byte, error) {
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, *i); err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buf.Bytes()), nil
}

func splitPieceHashes(pieces string) ([][20]byte, error) {
	const hashLen = 20

	if len(pieces)%hashLen != 0 {
		return nil, fmt.Errorf("invalid pieces length")
	}

	numHashes := len(pieces) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], pieces[i*hashLen:(i+1)*hashLen])
	}

	return hashes, nil
}
