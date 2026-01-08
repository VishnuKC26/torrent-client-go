package main

import (
	"fmt"
	"log"

	"github.com/VishnuKC26/torrent-client-go/torrentfile"
)

func main() {
	tf, err := torrentfile.Open("torrentfile/testdata/debian.torrent")
	if err != nil {
		log.Fatal(err)
	}
	tf.Announce = "http://bt.okmp3.ru:2710/announce"

	var peerID [20]byte
	copy(peerID[:], []byte("GO-TORRENT-CLIENT"))

	peers, err := tf.GetPeers(peerID, 6881)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Peers found:", len(peers))
	for i := 0; i < min(5, len(peers)); i++ {
		fmt.Println(peers[i].IP, peers[i].Port)
	}
}
