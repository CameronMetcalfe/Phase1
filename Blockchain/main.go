package main


import (
	"../go_blockchain/Blockchain"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

/*Adding a key-value pair using database BadgerDB, byte key-value pairs being stored
in folders

 */

type CommandLine struct{
	blockchain *Blockchain.BlockChain
}

//Usage commands
func (cl *CommandLine) printUsage(){
	fmt.Println("Usage: ")
	fmt.Println("To add a block: add -block <Block_Data>")
	fmt.Println("To print the chain: print")
}

func (cl *CommandLine) validateArgs(){
	if len(os.Args)<2{
		cl.printUsage()
		runtime.Goexit() //exits the application by shutting down the routines, useful for badger db; prevents it from corrupting
	}
}

func (cl *CommandLine) addBlock(data string){
	cl.blockchain.AddBlock(data)
	fmt.Println("Block Added")
}

func (cl *CommandLine) Printchain(){
	iter := cl.blockchain.ChainIter()

	for{
		block := iter.Next()
		fmt.Printf("Block Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Previous Hash: %x\n",block.PrevHash)
		proof := Blockchain.NewProof(block)
		fmt.Printf("Proof of Work: %s\n",strconv.FormatBool(proof.Validate()))
		fmt.Println()

		//break if the chain has come to an end
		if len(block.PrevHash) == 0{
			break
		}
	}
}

func (cl *CommandLine) run(){
	cl.validateArgs()

	addcommand := flag.NewFlagSet("add",flag.ExitOnError)
	printcommand := flag.NewFlagSet("print",flag.ExitOnError)
	adddata := addcommand.String("block","","Block Data")

	switch os.Args[1]{
	case "add":
		err := addcommand.Parse(os.Args[2:])
		Blockchain.Handle(err)
	case "print":
		err := printcommand.Parse(os.Args[2:])
		Blockchain.Handle(err)
	default:
		cl.printUsage()
		runtime.Goexit()
	}
	if addcommand.Parsed(){
		if *adddata == ""{
			addcommand.Usage()
			runtime.Goexit()
		}
		cl.addBlock(*adddata)
	}
	if printcommand.Parsed(){
		cl.Printchain()
	}
}

func main() {
	defer os.Exit(0) //only executes if go channel exits properly
	chain := Blockchain.InitBlockchain()
	defer chain.Database.Close()

	cl := CommandLine{chain}
	cl.run()
}
