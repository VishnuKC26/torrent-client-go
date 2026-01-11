BitTorrent Client in Go

A from-scratch BitTorrent client written in Go, implementing the core BitTorrent protocol including torrent parsing, tracker communication, peer handshakes, message exchange, piece downloading, and integrity verification.

This project is inspired by and structurally aligned with
👉 https://github.com/veggiedefender/torrent-client

and the BitTorrent protocol walkthrough by Jesse Li.

🚀 Features Implemented
Torrent & Metadata

Bencode decoding of .torrent files

Extraction of:

Tracker URL (announce)

File name & size

Piece length

Piece SHA-1 hashes

Computation of info-hash (SHA-1 of info dictionary)

Tracker Communication

HTTP tracker support

Tracker URL construction with correct query parameters

Compact peer list parsing

Peer discovery (IP:Port)

Peer Protocol (Wire Protocol)

TCP connection with peers

BitTorrent handshake implementation

Message framing:

Keep-alive

Choke / Unchoke

Interested / Not Interested

Have

Bitfield

Request

Piece

Defensive message parsing and validation

Download Engine

Concurrent peer workers (one connection per peer)

Piece-based downloading

Request pipelining with configurable backlog

SHA-1 integrity verification per piece

Automatic retry of failed pieces

Safe handling of misbehaving peers

Concurrency & Safety

Goroutines for parallel peer downloads

Channel-based work queue and result collection

Timeouts to avoid hanging peers

Deadlock-safe orchestration

📂 Project Structure
.
├── main.go                 # Entry point
├── torrentfile/            # Torrent parsing & tracker logic
│   ├── torrentfile.go
│   └── tracker.go
├── peers/                  # Peer representation & parsing
│   └── peers.go
├── handshake/              # BitTorrent handshake
│   └── handshake.go
├── message/                # Wire protocol messages
│   └── message.go
├── bitfield/               # Piece availability tracking
│   └── bitfield.go
├── client/                 # Peer connection abstraction
│   └── client.go
└── p2p/                    # Download orchestration
    └── p2p.go

▶️ How It Works (High Level)

Parse .torrent file

Decode bencoded metadata

Compute info-hash

Contact tracker

Retrieve peer list

Connect to peers

Perform handshake

Exchange protocol messages

Download pieces

Request blocks from peers

Assemble pieces

Verify SHA-1 integrity

Assemble file

Combine all verified pieces into final output

🧪 Usage
go run . path/to/file.torrent


Example:

go run . torrentfile/testdata/debian.torrent

⚠️ Important Notes on Swarm Behavior

This client correctly implements the BitTorrent protocol, but:

It does not upload pieces

It does not implement optimistic unchoking

In public swarms (e.g., Debian torrents), many peers may never unchoke this client due to BitTorrent’s tit-for-tat mechanism.

The client works best when:

Tested against a local seed

Used with small or lenient swarms

Extended with optimistic unchoking or upload support

This behavior is expected and protocol-compliant, not a bug.

🧠 What This Project Demonstrates

Deep understanding of network protocols

Binary data parsing and serialization

Concurrency with goroutines and channels

TCP socket programming

Distributed systems behavior (peer incentives)

Defensive programming against unreliable peers

📌 Future Improvements (Optional)

Optimistic unchoking

Upload (seeding) support

Peer prioritization (rarest-first)

Resume support

UDP tracker support

DHT peer discovery

📚 References

BitTorrent Protocol Specification

Jesse Li’s BitTorrent Walkthrough
https://blog.jse.li/posts/torrent/

veggiedefender torrent client
https://github.com/veggiedefender/torrent-client
