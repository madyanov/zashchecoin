package blockchain

import (
	"sync"

	"zc/bullshit"
)

const (
	blockGenerationIntervalSeconds     = 10 // in seconds
	difficultyAdjustmentIntervalBlocks = 10 // in blocks
)

type MemBlockchain struct {
	blocks []Block
	mtx    sync.RWMutex
}

func NewMemBlockchain(blocks []Block) *MemBlockchain {
	bc := &MemBlockchain{}

	if len(blocks) > 0 {
		bc.blocks = blocks
	} else {
		bc.blocks = append(bc.blocks, Genesis)
	}

	return bc
}

func (bc *MemBlockchain) AllBlocks() []Block {
	bc.mtx.RLock()
	defer bc.mtx.RUnlock()

	return bc.blocks
}

func (bc *MemBlockchain) AddBlock(block Block) bool {
	bc.mtx.Lock()
	defer bc.mtx.Unlock()

	if !bc.checkNewBlock(block, bc.lastBlock()) {
		return false
	}

	bc.blocks = append(bc.blocks, block)
	return true
}

func (bc *MemBlockchain) Replace(newBc *MemBlockchain) bool {
	bc.mtx.Lock()
	defer bc.mtx.Unlock()

	if !bc.shouldReplace(newBc) {
		return false
	}

	bc.blocks = newBc.blocks
	return true
}

func (bc *MemBlockchain) MineBlock(data []byte) Block {
	bc.mtx.RLock()
	lastBlock := bc.lastBlock()
	difficulty := bc.calculateDifficulty(lastBlock)
	bc.mtx.RUnlock()

	return mineBlock(
		lastBlock.Height+1,
		data,
		difficulty,
		lastBlock.Hash,
	)
}

func (bc *MemBlockchain) lastBlock() Block {
	return bc.blocks[len(bc.blocks)-1]
}

func (bc *MemBlockchain) checkChain() bool {
	if !bc.blocks[0].isGenesis() {
		return false
	}

	for i := 1; i < len(bc.blocks); i++ {
		if !bc.checkNewBlock(bc.blocks[i], bc.blocks[i-1]) {
			return false
		}
	}

	return true
}

func (bc *MemBlockchain) checkNewBlock(new Block, prev Block) bool {
	return prev.checkNextBlock(new) &&
		new.Difficulty == bc.calculateDifficulty(prev)
}

func (bc *MemBlockchain) calculateDifficulty(lastBlock Block) int {
	difficulty := lastBlock.Difficulty

	if lastBlock.Height != 0 && lastBlock.Height%difficultyAdjustmentIntervalBlocks == 0 {
		adjustmentBlock := bc.blocks[lastBlock.Height-difficultyAdjustmentIntervalBlocks+1]
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

func (bc *MemBlockchain) totalDifficulty() int {
	sum := 0

	for _, block := range bc.blocks {
		sum += 1 << block.Difficulty
	}

	return sum
}

func (bc *MemBlockchain) shouldReplace(newBc *MemBlockchain) bool {
	return newBc.totalDifficulty() > bc.totalDifficulty() &&
		newBc.checkChain()
}
