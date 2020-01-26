package blockchain

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	Bc *Blockchain
}

const (
	createBlockchainFlag  string = "createblockchain"
	createBlockchainUsage string = "  %s -address <address> - create a blockchain and send genesis block reward to <address>"
	printChainFlag        string = "printchain"
	printChainUsage       string = "  %s - print all the blocks of the blockchain"
	getBalanceFlag        string = "getbalance"
	getBalanceUsage       string = "  %s -address <address> - calculate the balance of <address>"
)

func (cli *CLI) createBlockchain(address string) {
	bc := CreateBlockchain(address)
	bc.Db.Close()
	fmt.Println("Done!")
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(fmt.Sprintf(createBlockchainUsage, createBlockchainFlag))
	fmt.Println(fmt.Sprintf(printChainUsage, printChainFlag))
	fmt.Println(fmt.Sprintf(getBalanceUsage, getBalanceFlag))
}

func (cli CLI) getBalance(address string) {
	bc := NewBlockchain(address)
	defer bc.Db.Close()

	utxos := bc.FindUTXO(address)
	var balance int
	for _, out := range utxos {
		balance += out.Value
	}
	fmt.Printf("Balance of `%s` is `%d`\n", address, balance)
}

func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}

	createBlockchainCommand := flag.NewFlagSet(createBlockchainFlag, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(printChainFlag, flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet(getBalanceFlag, flag.ExitOnError)

	createBlockchainAddress := createBlockchainCommand.String("address", "", "Recipient of the genesis block reward")
	getBalanceAddress := getBalanceCmd.String("address", "", "Get the balance of this address")

	switch os.Args[1] {
	case createBlockchainFlag:
		err := createBlockchainCommand.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case printChainFlag:
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case getBalanceFlag:
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockchainCommand.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCommand.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}
}

func (cli CLI) printChain() {
	bc := NewBlockchain("")
	defer bc.Db.Close()

	bci := bc.Iterator()

	for bci.currentHash != nil {

		block := bci.Next()

		fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
