BitTorrent Client in Go
A from-scratch BitTorrent client written in Go, implementing the core BitTorrent protocol. This project handles everything from bencode parsing to concurrent piece downloading and SHA-1 integrity verification.

This project is inspired by and structurally aligned with the torrent-client by veggiedefender and the BitTorrent protocol walkthrough by Jesse Li.

🚀 Features Implemented
Torrent & Metadata
Bencode Decoding: Full support for parsing .torrent files.

Metadata Extraction: Extracts Tracker URL, file metadata, piece length, and SHA-1 hashes.

Info-Hash Computation: Accurate SHA-1 calculation of the info dictionary.

Tracker & Peer Discovery
HTTP Tracker Support: Full communication with trackers via GET requests.

Compact Peer Parsing: Efficiently parses binary peer lists into IP:Port format.

Peer Protocol (Wire Protocol)
TCP Handshake: Reliable implementation of the BitTorrent handshake.

Message Framing: Support for Keep-alive, Choke/Unchoke, Interested, Have, Bitfield, Request, and Piece.

Validation: Defensive parsing to handle malformed peer data.

Download Engine
Concurrency: High-performance peer workers using Goroutines.

Pipelining: Configurable request backlog to maximize throughput.

Integrity Checks: Per-piece SHA-1 verification before writing to disk.

Fault Tolerance: Automatic retries for failed or corrupted pieces.

📂 Project Structure
Plaintext

.
├── main.go                 # Entry point: Orchestrates the download
├── torrentfile/            # Torrent parsing & Tracker communication logic
├── peers/                  # Peer representation & binary address parsing
├── handshake/              # Peer-to-peer handshake implementation
├── message/                # Wire protocol message serialization
├── bitfield/               # Tracking which pieces peers have
├── client/                 # TCP connection abstraction & message handling
└── p2p/                    # Download orchestration & concurrency logic
▶️ How It Works
Parse: Decodes the .torrent file and extracts the cryptographic info-hash.

Discover: Contacts the tracker to receive a list of active peers in the swarm.

Handshake: Establishes a TCP connection and verifies the protocol handshake with peers.

Download: Workers request blocks, assemble them into pieces, and verify them against the info-hash.

Assemble: Once all pieces are verified, they are joined to create the final output file.

🧪 Usage
Ensure you have Go installed.

Bash

# Clone the repository
git clone https://github.com/your-username/your-repo-name.git
cd your-repo-name

# Run the client
go run . path/to/file.torrent
Example:

Bash

go run . torrentfile/testdata/debian.torrent
⚠️ Important Notes on Swarm Behavior
This client implements the core download protocol but does not currently implement seeding (uploading) or optimistic unchoking.

[!IMPORTANT] Because BitTorrent relies on a "tit-for-tat" mechanism, some peers in public swarms (like Debian) may refuse to unchoke this client because it is not uploading back to them. This is expected behavior for a download-only educational client.

🧠 Lessons Learned
Network Protocols: Deep dive into binary framing and stateful TCP communication.

Concurrency Patterns: Using Go channels to manage a work queue across multiple workers.

Data Integrity: Implementing cryptographic verification in a streaming context.

Defensive Programming: Handling flaky network connections and unreliable peers.

📌 Future Improvements
[ ] UDP Tracker Support: Support for the udp:// tracker protocol.

[ ] Seeding: Implement uploading and optimistic unchoking.

[ ] Rarest-First: Optimize piece selection to improve swarm health.

[ ] Resume: Save progress to disk to allow resuming interrupted downloads.

[ ] DHT: Implement distributed hash tables for trackerless discovery.

📚 References
BitTorrent Protocol Specification

Jesse Li’s BitTorrent Walkthrough

veggiedefender/torrent-client
