package awesomeProject1

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
)

type Block struct{
	Hash []byte
	Data []byte
	PrevHash []byte
	Nonce int
}

func (pow *ProofOfWork) Validate() bool{
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(pow.Target) == -1
}

func (b *Block) Serialize() []byte{
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err:= encoder.Encode(b)

	Handle(err)
	return res.Bytes()
}

func Deserialize(data []byte) *Block{
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	Handle(err)
	return &block

}

func Handle(err error){
	if err != nil{
		log.Panic(err)
	}
}

func createBlock(data string, prevHash []byte) *Block{
	block := &Block{[]byte{},[]byte(data),prevHash,0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis() *Block{
	return createBlock("Genesis",[]byte{})
}