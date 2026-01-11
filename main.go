// package main

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net"
// 	"time"

// 	"github.com/VishnuKC26/torrent-client-go/handshake"
// 	"github.com/VishnuKC26/torrent-client-go/torrentfile"
// )

// func main() {
// 	// Load torrent
// 	tf, err := torrentfile.Open("torrentfile/testdata/debian.torrent")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Using tracker:", tf.Announce)

// 	// Peer ID (must be 20 bytes)
// 	peerID := generatePeerID()

// 	// Get peers from tracker
// 	peers, err := tf.GetPeers(peerID, 6881)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Peers from tracker:", len(peers))
// 	fmt.Println("Testing handshakes...")

// 	// IMPORTANT: do not try all peers
// 	maxTests := 50
// 	if len(peers) < maxTests {
// 		maxTests = len(peers)
// 	}

// 	success := 0

// 	for i := 0; i < maxTests; i++ {
// 		peer := peers[i]
// 		addr := net.JoinHostPort(peer.IP.String(), fmt.Sprintf("%d", peer.Port))

// 		fmt.Println("→ Dialing", addr)

// 		conn, err := net.DialTimeout("tcp", addr, 6*time.Second)
// 		if err != nil {
// 			fmt.Println("  ✗ dial failed")
// 			continue
// 		}

// 		// Ensure we don't block forever
// 		conn.SetDeadline(time.Now().Add(6 * time.Second))

// 		// Send handshake
// 		h := handshake.Handshake{
// 			InfoHash: tf.InfoHash,
// 			PeerID:   peerID,
// 		}

// 		if _, err := conn.Write(h.Serialize()); err != nil {
// 			fmt.Println("  ✗ write failed")
// 			conn.Close()
// 			continue
// 		}

// 		// Read handshake
// 		resp, err := handshake.Read(conn)
// 		conn.Close()

// 		if err != nil {
// 			fmt.Println("  ✗ handshake read failed")
// 			continue
// 		}

// 		// Validate info hash
// 		if resp.InfoHash != tf.InfoHash {
// 			fmt.Println("  ✗ info_hash mismatch")
// 			continue
// 		}

// 		fmt.Println("  ✓ handshake successful")
// 		success++
// 	}

// 	fmt.Printf("\nHandshake successes: %d / %d\n", success, maxTests)
// }

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/VishnuKC26/torrent-client-go/p2p"
	"github.com/VishnuKC26/torrent-client-go/torrentfile"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <torrent-file>", os.Args[0])
	}

	torrentPath := os.Args[1]

	tf, err := torrentfile.Open(torrentPath)
	if err != nil {
		log.Fatal(err)
	}

	var peerID [20]byte
	copy(peerID[:], []byte("GO-TORRENT_CLIENT-01"))

	peers, err := tf.GetPeers(peerID, 6881)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d peers\n", len(peers))

	t := p2p.Torrent{
		Peers:       peers,
		PeerId:      peerID,
		InfoHash:    tf.InfoHash,
		PieceHashes: tf.PieceHashes,
		PieceLength: tf.PieceLength,
		Length:      tf.Length,
		Name:        tf.Name,
	}

	buf, err := t.Download()
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(tf.Name, buf, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Download completed:", tf.Name)
}
