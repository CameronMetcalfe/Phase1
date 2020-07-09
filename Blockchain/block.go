package Blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct{
	Hash []byte
	Data []byte
	PrevHash []byte
	Nonce int
}



func CreateBlock(data string, prevHash []byte)*Block{
	block := &Block{[]byte{},[]byte(data),prevHash,0}
	proof := NewProof(block)
	nonce,hash := proof.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}



func Genesis() *Block{
	return CreateBlock("Genesis", []byte{})
}


/*Need to serialize and deserialize, blocks for the database

 */

func (block *Block) serial() []byte{
	var res bytes.Buffer
	enc := gob.NewEncoder(&res)
	err := enc.Encode(block)

	Handle(err)

	return res.Bytes()

}

//outputs a pointer to the block

func deserial(data []byte) *Block{
	var block Block

	dec := gob.NewDecoder(bytes.NewReader(data))

	err := dec.Decode(&block)

	Handle(err)
	return &block
}

func Handle(err error){
	if err != nil{
		log.Panic(err)
	}
}
