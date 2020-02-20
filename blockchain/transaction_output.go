package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

// TXOutputs collects TXOutput
type TXOutputs struct {
	Ouputs []TXOutput
}

// DeserializeOutputs deserializes TXOutputs
func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
