package main

import (
	"github.com/pdbrito/simple-blockchain/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer bc.Db.Close()

	cli := blockchain.CLI{Bc: bc}
	cli.Run()
}
