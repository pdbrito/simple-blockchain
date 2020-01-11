package main

import (
	"fmt"
	"github.com/pdbrito/simple-blockchain/blockchain"
	"strconv"
)

func main() {
	bc := blockchain.NewBlockchain()

	bc.AddBlock("Send 100 Credits to Han Solo")
	bc.AddBlock("Send 100 Credits to Ben Solo")

	for _, block := range bc.Blocks {
		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
