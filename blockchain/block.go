package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"strings"
	"time"
)

var Genesis = newGenesis()

type Block struct {
	Height       int
	Timestamp    int64
	Data         []byte
	Difficulty   int
	Nonce        int64
	PreviousHash []byte
	Hash         []byte
}

func newGenesis() Block {
	return newBlock(
		0,
		1626620944,
		[]byte("genesis"),
		0,
		0,
		[]byte{},
	)
}

func newBlock(
	height int,
	timestamp int64,
	data []byte,
	difficulty int,
	nonce int64,
	previousHash []byte,
) Block {
	block := Block{}
	block.Height = height
	block.Timestamp = timestamp
	block.Data = data
	block.Difficulty = difficulty
	block.Nonce = nonce
	block.PreviousHash = previousHash
	block.Hash = block.calculateHash()
	return block
}

func (b Block) Equal(block Block) bool {
	return b.Height == block.Height &&
		b.Timestamp == block.Timestamp &&
		bytes.Equal(b.Data, block.Data) &&
		b.Difficulty == block.Difficulty &&
		b.Nonce == block.Nonce &&
		bytes.Equal(b.PreviousHash, block.PreviousHash) &&
		bytes.Equal(b.Hash, block.Hash)
}

func (b Block) calculateHash() []byte {
	bytes := []byte{}
	bytes = append(bytes, []byte(strconv.Itoa(b.Height))...)
	bytes = append(bytes, []byte(strconv.FormatInt(b.Timestamp, 10))...)
	bytes = append(bytes, b.Data...)
	bytes = append(bytes, []byte(strconv.Itoa(b.Difficulty))...)
	bytes = append(bytes, []byte(strconv.FormatInt(b.Nonce, 10))...)
	bytes = append(bytes, b.PreviousHash...)
	hash := sha256.Sum256(bytes)
	return hash[:]
}

func (b Block) isGenesis() bool {
	return b.Equal(Genesis)
}

func (b Block) checkBlock() bool {
	return b.checkContent() && b.checkDifficulty()
}

func (b Block) checkContent() bool {
	return bytes.Equal(b.calculateHash(), b.Hash)
}

func (b Block) checkDifficulty() bool {
	return b.Difficulty > 0 &&
		strings.HasPrefix(
			bytesToBinaryString(b.Hash),
			strings.Repeat("0", b.Difficulty),
		)
}

func (b Block) checkNextBlock(next Block) bool {
	if next.Height != b.Height+1 {
		return false
	}

	if !bytes.Equal(next.PreviousHash, b.Hash) {
		return false
	}

	// todo: move 60 in constants
	if b.Timestamp-60 > next.Timestamp || next.Timestamp > time.Now().Unix() {
		return false
	}

	if !next.checkBlock() {
		return false
	}

	return true
}

func bytesToBinaryString(bytes []byte) string {
	str := ""

	for _, byte := range bytes {
		str += strconv.FormatInt(int64(byte), 2)
	}

	return str
}
