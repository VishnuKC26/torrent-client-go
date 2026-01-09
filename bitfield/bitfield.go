package bitfield


type Bitfield []byte

func (bf Bitfield) Has(index int) bool {

	byteIndex := index / 8
	bitOffset := index % 8

	if byteIndex <0 || byteIndex >= len(bf) { 

		return false
	}

	return bf[byteIndex] >> uint(7-bitOffset)&1 != 0

}

func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	bitOffset := index % 8

	
	if byteIndex < 0 || byteIndex >= len(bf) {
		return
	}
	bf[byteIndex] |= 1 << uint(7 - bitOffset)
}

// in BitTorrent each file is split into pieces

// a piece is a unit of meaningful transfer i.e a piece is transmitted by a peer by dividing it further into blocks but still until and unless all blocks of a piece have been received, they have no value of their own, they are only meaningful the peer has transmitted all blocks of the piece

// But first the we need to know which pieces does the peer have which it tells us using bitfields

// a bitfield is a slice of 8 bits and each bit represents a piece so one bitfield represents a set of 8 pieces

// those bits of a bitfield which are set indicate that that piece is available with the peer we are talking to

// Eg if we need piece 13 of a file and we need to check its availability with a particular peer then the Has func comes into play

// byteIndex := index/8 i.e 13/8 which is 1

// this tells us to look in byte 1 sent by the peer

// bitOffset := index % 8 => 13 % 8 which is 5

// eg: bf[1] = 00100000

// bit:   7 6 5 4 3 2 1 0
// value: 0 0 0 0 0 1 0 0   this is the big-endian format 

// so second bit i.e (7-bitOffset) but is set

// now we right-shift the bitfield by the index position which makes the concerned bit position the elast significant (rightmost bit)

// then we AND this with 1 which translates to 00000001 and if the last bit of our modifed bitfield is also 1 then the result would be 1 and != 0 would return true so in this way we can know if out concerned bit is set or not by always making it the rightmost bit and ANDing with 1