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

	fmt.Println("Name:", tf.Name)
	fmt.Println("Tracker:", tf.Announce)
	fmt.Println("File size:", tf.Length)
	fmt.Println("Piece length:", tf.PieceLength)
	fmt.Println("Number of pieces:", len(tf.PieceHashes))
	fmt.Printf("Info hash: %x\n", tf.InfoHash)
}
