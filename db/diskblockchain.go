package db

import (
	"bufio"
	"encoding/gob"
	"io"
	"os"
	"sync"

	"zc/blockchain"
	"zc/bullshit"
	"zc/encoding"
)

type DiskBlockchain struct {
	mem *blockchain.MemBlockchain
	f   *os.File
	w   *bufio.Writer
	enc *gob.Encoder
	mtx sync.Mutex
}

func NewDiskBlockchain(bcPath string) *DiskBlockchain {
	f, err := os.OpenFile(bcPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	bullshit.FailIf(err)
	// never close file

	r := bufio.NewReader(f)
	w := bufio.NewWriter(f)

	enc, buf, err := encoding.NewEncoderWithoutHeader(w, &blockchain.Block{})
	bullshit.FailIf(err)

	dec, err := encoding.NewDecoderWithoutHeader(r, buf, &blockchain.Block{})
	bullshit.FailIf(err)

	blocks := []blockchain.Block{}

	for {
		block := blockchain.Block{}

		err := dec.Decode(&block)
		if err == io.EOF {
			break
		}
		bullshit.FailIf(err)

		blocks = append(blocks, block)
	}

	memBc, err := blockchain.NewMemBlockchain(blocks)
	bullshit.FailIf(err)

	bc := &DiskBlockchain{
		mem: memBc,
		f:   f,
		w:   w,
		enc: enc,
	}

	// write the genesis block
	if len(blocks) == 0 {
		bc.writeBlock(blockchain.Genesis)
		bc.flush()
	}

	return bc
}

func (bc *DiskBlockchain) AllBlocks() []blockchain.Block {
	return bc.mem.AllBlocks()
}

func (bc *DiskBlockchain) AddBlock(block blockchain.Block) bool {
	if !bc.mem.AddBlock(block) {
		return false
	}

	bc.mtx.Lock()
	defer bc.mtx.Unlock()

	bc.writeBlock(block)
	bc.flush()

	return true
}

func (bc *DiskBlockchain) Replace(newBc *blockchain.MemBlockchain) bool {
	if !bc.mem.Replace(newBc) {
		return false
	}

	bc.mtx.Lock()
	defer bc.mtx.Unlock()

	err := bc.f.Truncate(0)
	bullshit.FailIf(err)

	_, err = bc.f.Seek(0, io.SeekStart)
	bullshit.FailIf(err)

	blocks := newBc.AllBlocks()
	for _, block := range blocks {
		bc.writeBlock(block)
	}

	bc.flush()

	return true
}

func (bc *DiskBlockchain) MineBlock(data []byte) blockchain.Block {
	return bc.mem.MineBlock(data)
}

func (bc *DiskBlockchain) Weight() int {
	return bc.mem.Weight()
}

func (bc *DiskBlockchain) writeBlock(block blockchain.Block) {
	err := bc.enc.Encode(block)
	bullshit.FailIf(err)
}

func (bc *DiskBlockchain) flush() {
	err := bc.w.Flush()
	bullshit.FailIf(err)
}
