package Blockchain

import (
	"fmt"

	"../github.com/dgraph-io/badger"
)

const dbpath = "./tmp/blocks"

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

//structure to iterate through the blockchain
type Iterator struct {
	currenth []byte
	DB       *badger.DB
}

func (chain *BlockChain) AddBlock(data string) {
	var lasthash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh")) //"lh" stands for lasthash, gives a pointer to item stored in the lh key and an error
		Handle(err)
		var hash []byte
		err = item.Value(func(val []byte) error {
			hash = append([]byte{}, val...)
			return nil
		})
		return err
	})
	Handle(err)

	fmt.Println("Printing block lasthash")
	fmt.Print(lasthash)
	newblock := CreateBlock(data, lasthash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newblock.Hash, newblock.serial())
		Handle(err)
		err = txn.Set([]byte("lh"), newblock.Hash)

		chain.LastHash = newblock.Hash
		return err
	})
	Handle(err)
}

/*Initializing Blockchain using BadgerDB.

 */

func InitBlockchain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbpath)
	opts.Dir = dbpath
	opts.ValueDir = dbpath

	db, err := badger.Open(opts)
	Handle(err)

	//two ways to access the database, view or update, update allows for writes
	err = db.Update(func(txn *badger.Txn) error {
		//Checking if there is already a blockchain structure initialized in the database
		//if txn.Get returns KeyNotFound error, means no structure available, so initialize
		//a genesis block
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No Blockchain exists")
			genesis := Genesis()
			fmt.Println("Genesis Block proved")
			err = txn.Set(genesis.Hash, genesis.serial())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			var Hash []byte
			err = item.Value(func(val []byte) error {
				Hash = append([]byte{}, val...)
				return nil
			})
			return err
		}

	})
	Handle(err)
	chain := BlockChain{lastHash, db}
	return &chain
}

//Manually implementing an iterator instead of pre-built functions, iterating backwards as we start from the last hash
func (chain *BlockChain) ChainIter() *Iterator {
	iter := &Iterator{chain.LastHash, chain.Database}
	return iter
}

//Function that returns a pointer to the next block
func (iter *Iterator) Next() *Block {
	var block *Block

	//in order to find the next block, we read through the chain and hence use the view function instead of update
	err := iter.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.currenth)
		Handle(err)

		var encblock []byte
		err = item.Value(func(val []byte) error {
			encblock = append([]byte{}, val...)
			return nil
		})
		block = deserial(encblock)

		return err
	})
	Handle(err)

	iter.currenth = block.PrevHash

	return block
}
