package blockchain

import "time"

func MineBlock(
	height int,
	data []byte,
	difficulty int,
	previousHash []byte,
) Block {
	nonce := 0

	for {
		block := newBlock(
			height,
			time.Now().Unix(),
			data,
			difficulty,
			nonce,
			previousHash,
		)

		if block.checkDifficulty() {
			return block
		}

		nonce++
	}
}
