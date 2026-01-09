package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/VishnuKC26/torrent-client-go/handshake"
	"github.com/VishnuKC26/torrent-client-go/torrentfile"
)

func main() {
	// Load torrent
	tf, err := torrentfile.Open("torrentfile/testdata/debian.torrent")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Using tracker:", tf.Announce)

	// Peer ID (must be 20 bytes)
	peerID := generatePeerID()

	// Get peers from tracker
	peers, err := tf.GetPeers(peerID, 6881)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Peers from tracker:", len(peers))
	fmt.Println("Testing handshakes...")

	// IMPORTANT: do not try all peers
	maxTests := 50
	if len(peers) < maxTests {
		maxTests = len(peers)
	}

	success := 0

	for i := 0; i < maxTests; i++ {
		peer := peers[i]
		addr := net.JoinHostPort(peer.IP.String(), fmt.Sprintf("%d", peer.Port))

		fmt.Println("→ Dialing", addr)

		conn, err := net.DialTimeout("tcp", addr, 6*time.Second)
		if err != nil {
			fmt.Println("  ✗ dial failed")
			continue
		}

		// Ensure we don't block forever
		conn.SetDeadline(time.Now().Add(6 * time.Second))

		// Send handshake
		h := handshake.Handshake{
			InfoHash: tf.InfoHash,
			PeerID:   peerID,
		}

		if _, err := conn.Write(h.Serialize()); err != nil {
			fmt.Println("  ✗ write failed")
			conn.Close()
			continue
		}

		// Read handshake
		resp, err := handshake.Read(conn)
		conn.Close()

		if err != nil {
			fmt.Println("  ✗ handshake read failed")
			continue
		}

		// Validate info hash
		if resp.InfoHash != tf.InfoHash {
			fmt.Println("  ✗ info_hash mismatch")
			continue
		}

		fmt.Println("  ✓ handshake successful")
		success++
	}

	fmt.Printf("\nHandshake successes: %d / %d\n", success, maxTests)
}
func generatePeerID() [20]byte {
	var pid [20]byte

	// Azureus-style client ID (8 bytes)
	copy(pid[:8], []byte("-GO0001-"))

	// Fill remaining 12 bytes securely
	if _, err := rand.Read(pid[8:]); err != nil {
		panic(err)
	}

	return pid
}
