package torrentfile

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/jackpal/bencode-go"
)

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

type Peer struct {
	IP   net.IP
	Port uint16
}

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {

	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{

		"info_hash":  []string{string(t.InfoHash[:])},
		"peer-id":    []string{string(peerID[:])},
		"port":       []string{fmt.Sprintf("%d", port)},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{fmt.Sprintf("%d", t.Length)},
		"compact":    []string{"1"},
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}

func parsePeers(peersBin string)([]Peer,error){

	const peerSize = 6 // 4 bytes IP + 2 bytes Port

	numPeers := len(peersBin)/peerSize

	peers := make([]Peer,0,numPeers)

	for i:=0; i+peerSize <= len(peersBin); i+= peerSize {

		ip := net.IP(peersBin[i:i+4])
		port := binary.BigEndian.Uint16([]byte(peersBin[i+4:i+6]))

		peers = append(peers,Peer{
			IP: ip,
			Port: port,
		})
	}

	return peers , nil
}

func (t *TorrentFile) GetPeers(peerID [20]byte , port uint16)([]Peer, error){

	url,err := t.buildTrackerURL(peerID,port)
	if err != nil {
		return nil, err
	}

	resp , err := http.Get(url)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	var tr bencodeTrackerResp
	if err := bencode.Unmarshal(resp.Body,&tr); err != nil {
		return nil,err
	}

	return parsePeers(tr.Peers)
}