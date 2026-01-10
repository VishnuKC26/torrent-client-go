package torrentfile

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/VishnuKC26/torrent-client-go/peer"
	"github.com/jackpal/bencode-go"
)

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

/* ---------- Tracker URL ---------- */

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {

	if t.Announce == "" {
		return "", fmt.Errorf("empty announce URL")
	}

	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	// Only HTTP trackers supported here
	if base.Scheme != "http" && base.Scheme != "https" {
		return "", fmt.Errorf("unsupported tracker scheme: %s", base.Scheme)
	}

	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{fmt.Sprintf("%d", port)},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{fmt.Sprintf("%d", t.Length)},
		"compact":    []string{"1"},
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}

/* ---------- Parse compact peers ---------- */

func parsePeers(peersBin string) ([]peer.Peer, error) {

	const peerSize = 6 // 4 bytes IP + 2 bytes port

	if len(peersBin)%peerSize != 0 {
		return nil, fmt.Errorf("invalid peers binary length")
	}

	numPeers := len(peersBin) / peerSize
	peers := make([]peer.Peer, 0, numPeers)

	for i := 0; i < len(peersBin); i += peerSize {

		ip := net.IP([]byte(peersBin[i : i+4]))
		port := binary.BigEndian.Uint16([]byte(peersBin[i+4 : i+6]))

		peers = append(peers, peer.Peer{
			IP:   ip,
			Port: port,
		})
	}

	return peers, nil
}

/* ---------- Get peers from tracker ---------- */

func (t *TorrentFile) GetPeers(peerID [20]byte, port uint16) ([]peer.Peer, error) {

	trackerURL, err := t.buildTrackerURL(peerID, port)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(trackerURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tracker returned status %s", resp.Status)
	}

	var tr bencodeTrackerResp
	if err := bencode.Unmarshal(resp.Body, &tr); err != nil {
		return nil, err
	}

	return parsePeers(tr.Peers)
}
