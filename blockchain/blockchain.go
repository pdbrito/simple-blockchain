package blockchain

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := newBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func newGenesisBlock() *Block {
	return newBlock("A long time ago in a galaxy far far away....", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{newGenesisBlock()}}
}
