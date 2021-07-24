package p2p

import (
	"encoding/binary"
	"encoding/gob"
	"io"

	"zc/blockchain"
	"zc/bullshit"
)

const (
	respTypeBlock = byte(iota)
	respTypeChain
)

// Block
type respBlock struct {
	block blockchain.Block
}

func (resp *respBlock) encode(w io.Writer) {
	err := binary.Write(w, binary.BigEndian, respTypeBlock)
	bullshit.WarnIf(err)

	enc := gob.NewEncoder(w)
	err = enc.Encode(resp.block)
	bullshit.WarnIf(err)
}

func (resp *respBlock) decode(r io.Reader) {
	dec := gob.NewDecoder(r)
	err := dec.Decode(&resp.block)
	bullshit.WarnIf(err)
}

// Blockchain
type respChain struct {
	blocks []blockchain.Block
}

func (resp *respChain) encode(w io.Writer) {
	err := binary.Write(w, binary.BigEndian, respTypeChain)
	bullshit.WarnIf(err)

	enc := gob.NewEncoder(w)
	err = enc.Encode(resp.blocks)
	bullshit.WarnIf(err)
}

func (resp *respChain) decode(r io.Reader) {
	dec := gob.NewDecoder(r)
	err := dec.Decode(&resp.blocks)
	bullshit.WarnIf(err)
}
