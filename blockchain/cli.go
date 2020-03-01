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
	sendFlag              string = "send"
	sendUsage             string = "  %s -from <from> -to <to> -amount <amount> - Send <amount> from <from> to <to>"
	createWalletFlag      string = "createwallet"
	createWalletUsage     string = "  %s - generates a new key-pair and saves it into the wallet file"
)

func (cli *CLI) createBlockchain(address string) {
	if !ValidateAddress(address) {
		log.Panic("error: Address is not valid")
	}
	bc := CreateBlockchain(address)
	defer bc.Db.Close()

	UTXOSet := UTXOSet{}
	UTXOSet.Reindex()

	fmt.Println("Done!")
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(fmt.Sprintf(createBlockchainUsage, createBlockchainFlag))
	fmt.Println(fmt.Sprintf(printChainUsage, printChainFlag))
	fmt.Println(fmt.Sprintf(getBalanceUsage, getBalanceFlag))
	fmt.Println(fmt.Sprintf(sendUsage, sendFlag))
	fmt.Println(fmt.Sprintf(createWalletUsage, createWalletFlag))
}

func (cli CLI) getBalance(address string) {
	if !ValidateAddress(address) {
		log.Panic("error: address is not valid")
	}
	bc := NewBlockchain()
	UTXOSet := UTXOSet{bc}
	defer bc.Db.Close()

	var balance int
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of `%s` is `%d`\n", address, balance)
}

func (cli CLI) send(from, to string, amount int) {
	bc := NewBlockchain()
	UTXOSet := UTXOSet{bc}
	defer bc.Db.Close()

	tx := NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTX := NewCoinbaseTX(from, "")
	txs := []*Transaction{cbTX, tx}

	newBlock := bc.MineBlock(txs)
	UTXOSet.Update(newBlock)
	fmt.Println("Success!")
}

func (cli CLI) createWallet() {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}

func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}

	createBlockchainCmd := flag.NewFlagSet(createBlockchainFlag, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(printChainFlag, flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet(getBalanceFlag, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(sendFlag, flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet(createWalletFlag, flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "Recipient of the genesis block reward")
	getBalanceAddress := getBalanceCmd.String("address", "", "Get the balance of this address")
	sendFromAddress := sendCmd.String("fromAddress", "", "Address to take funds from")
	sendToAddress := sendCmd.String("toAddress", "", "Address to send funds to")
	sendAmount := sendCmd.Int("amount", 0, "Amount to transfer")

	switch os.Args[1] {
	case createBlockchainFlag:
		err := createBlockchainCmd.Parse(os.Args[2:])
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
	case sendFlag:
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case createWalletFlag:
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
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
	if sendCmd.Parsed() {
		if *sendFromAddress == "" {
			sendCmd.Usage()
			os.Exit(1)
		}
		if *sendToAddress == "" {
			sendCmd.Usage()
			os.Exit(1)
		}
		if *sendAmount == 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFromAddress, *sendToAddress, *sendAmount)
	}
	if createWalletCmd.Parsed() {
		cli.createWallet()
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
