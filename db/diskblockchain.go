package db

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
	"os"
	"sync"

	"zc/blockchain"
	"zc/bullshit"
)

type DiskBlockchain struct {
	mem *blockchain.MemBlockchain
	f   *os.File
	w   *bufio.Writer
	enc *gob.Encoder
	mtx sync.Mutex
}

func NewDiskBlockchain(bcPath string) *DiskBlockchain {
	f, err := os.OpenFile(bcPath, os.O_RDWR|os.O_CREATE, 0755)
	bullshit.FailIf(err)
	// never close file

	r := bufio.NewReader(f)
	dec := gob.NewDecoder(r)

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

	log.Println("Loaded blocks:", len(blocks))

	w := bufio.NewWriter(f)
	bc := &DiskBlockchain{
		mem: blockchain.NewMemBlockchain(blocks),
		f:   f,
		w:   w,
		enc: gob.NewEncoder(w),
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

	_, err = bc.f.Seek(0, 0)
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

func (bc *DiskBlockchain) writeBlock(block blockchain.Block) {
	err := bc.enc.Encode(block)
	bullshit.FailIf(err)
}

func (bc *DiskBlockchain) flush() {
	err := bc.w.Flush()
	bullshit.FailIf(err)
}
