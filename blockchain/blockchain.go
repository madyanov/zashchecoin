package blockchain

import (
	"sync"

	"zc/bullshit"
)

const (
	blockGenerationIntervalSeconds     = 10 // in seconds
	difficultyAdjustmentIntervalBlocks = 10 // in blocks
)

type Blockchain struct {
	blocks []Block
	mutex  sync.RWMutex
}

func NewBlockchain(blocks []Block) *Blockchain {
	bc := &Blockchain{}

	if len(blocks) > 0 {
		bc.blocks = blocks
	} else {
		bc.blocks = append(bc.blocks, Genesis)
	}

	return bc
}

func (b *Blockchain) AllBlocks() []Block {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.blocks
}

func (b *Blockchain) MineBlock(data []byte) Block {
	b.mutex.RLock()
	lastBlock := b.lastBlock()
	difficulty := b.calculateDifficulty(lastBlock)
	b.mutex.RUnlock()

	return MineBlock(
		lastBlock.Height+1,
		data,
		difficulty,
		lastBlock.Hash,
	)
}

func (b *Blockchain) AddBlock(new Block) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if !b.checkNewBlock(new, b.lastBlock()) {
		return false
	}

	b.blocks = append(b.blocks, new)
	return true
}

func (b *Blockchain) Replace(new *Blockchain) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if !b.shouldReplace(new) {
		return false
	}

	b.blocks = new.blocks
	return true
}

func (b *Blockchain) Height() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.lastBlock().Height
}

func (b *Blockchain) lastBlock() Block {
	return b.blocks[len(b.blocks)-1]
}

func (b *Blockchain) checkChain() bool {
	if !b.blocks[0].isGenesis() {
		return false
	}

	for i := 1; i < len(b.blocks); i++ {
		if !b.checkNewBlock(b.blocks[i], b.blocks[i-1]) {
			return false
		}
	}

	return true
}

func (b *Blockchain) checkNewBlock(new Block, prev Block) bool {
	return prev.checkNextBlock(new) &&
		new.Difficulty == b.calculateDifficulty(prev)
}

func (b *Blockchain) calculateDifficulty(lastBlock Block) int {
	difficulty := lastBlock.Difficulty

	if lastBlock.Height != 0 && lastBlock.Height%difficultyAdjustmentIntervalBlocks == 0 {
		adjustmentBlock := b.blocks[lastBlock.Height-difficultyAdjustmentIntervalBlocks+1]
		timeExpected := int64(blockGenerationIntervalSeconds * difficultyAdjustmentIntervalBlocks)
		timeTaken := lastBlock.Timestamp - adjustmentBlock.Timestamp

		if timeTaken < timeExpected/2 {
			difficulty = adjustmentBlock.Difficulty + 1
		} else if timeTaken > timeExpected*2 {
			difficulty = adjustmentBlock.Difficulty - 1
		} else {
			difficulty = adjustmentBlock.Difficulty
		}
	}

	return bullshit.Max(1, difficulty)
}

func (b *Blockchain) totalDifficulty() int {
	sum := 0

	for _, block := range b.blocks {
		sum += 1 << block.Difficulty
	}

	return sum
}

func (b *Blockchain) shouldReplace(new *Blockchain) bool {
	return new.totalDifficulty() > b.totalDifficulty() &&
		new.checkChain()
}
