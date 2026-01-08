// What torrentfile.go does

// announce (tracker URL)

// info.name

// info.length

// info.piece length

// info.pieces

// info-hash (SHA-1) of bencoded info dict
package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"os"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte // a list in which each element is a 20 byte hash
	PieceLength int
	Length      int
	Name        string
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func Open(path string) (*TorrentFile, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bt bencodeTorrent
	if err := bencode.Unmarshal(file, &bt); err != nil {
		return nil, err
	}

	infoHash, err := bt.infoHash()
	if err != nil {
		return nil,err
	}

	pieceHashes,err := splitPieceHashes(bt.Info.Pieces)

	if err != nil {
		return nil,err
	}

	tf := TorrentFile{
		Announce: bt.Announce,
		InfoHash: infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bt.Info.PieceLength,
		Length: bt.Info.Length,
		Name: bt.Info.Name,
	}

	return &tf , nil
}

func( bt *bencodeTorrent) infoHash()([20]byte,error){

	var buf bytes.Buffer

	if err := bencode.Marshal(&buf,bt.Info); err != nil {
		return [20]byte{},err
	}

	return sha1.Sum(buf.Bytes()),nil
}

func splitPieceHashes(pieces string)([][20]byte,error){

	const hashLen = 20

	numHashes := len(pieces)/hashLen

	hashes := make([][20]byte,numHashes)

	for i := range numHashes{
		copy(hashes[i][:],pieces[i*hashLen:(i+1)*hashLen])
	}

	return hashes,nil
}