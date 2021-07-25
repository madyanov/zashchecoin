package blockchain

type Blockchain interface {
	AllBlocks() []Block
	AddBlock(block Block) bool
	Replace(newBc *MemBlockchain) bool
	MineBlock(data []byte) Block
	Weight() int
}
