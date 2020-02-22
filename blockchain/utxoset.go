package blockchain

import (
	"encoding/hex"
	bolt "go.etcd.io/bbolt"
	"log"
)

// UTXOSet represents the UTXO set
type UTXOSet struct {
	Blockchain *Blockchain
}

const utxoBucket = "chainstate"

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.Db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)

			for outIdx, output := range outs.Ouputs {
				if output.IsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += output.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

// FindUTXO finds UTXO for a public key hash
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.Db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)

			for _, output := range outs.Ouputs {
				if output.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, output)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}
