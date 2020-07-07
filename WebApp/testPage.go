package main

import (
	"fmt"
	"html/template"
	"net/http"
	"encoding/json"

	"../blockchain/blockchain"
	"../network"
)

type BlockData struct { //will be converted to byte[] when added to actual blockchain
	ItemName string
	ItemDescription string
	ItemPrice int
	Username        string
	BlockType       int
}

var nodes []*network.Node
chain := blockchain.InitBlockChain()

func main() {
	//need to set up Nodes here
	//maybe take in number of nodes from command line arg
	//numNodes := 5
	//var chain blockchain.Blockchain
	

	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/list", ListingPage)
	http.HandleFunc("/add-listing", AddListing)
	http.HandleFunc("/purchase", PurchasePage)
	http.HandleFunc("/bought", BuyListing)
	http.HandleFunc("/add-user", AddUser)

	http.ListenAndServe(":8080", nil)
}

func InitBlockStruct(blockType int, itemName string, desc string, price int, username string) *BlockData {
	result := *BlockData{}
	result.BlockType = blockType
	result.ItemName = itemName
	result.ItemDescription = desc
	result.ItemPrice = price
	result.Username = username
	return result
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	//title := r.URL.Path[1:] //r.URL.Path[]
	testBlock := Block{}
	testBlock.Hash = 3

	t, _ := template.ParseFiles("testPage.html")
	t.Execute(w, testBlock)
}

//These functions put data into the blockchain in response to actions in the web app
func AddListing(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add listing function called")
	listingBlock := InitBlockStruct(1, r.FormValue("item-name"), r.FormValue("item-description"), r.FormValue("price"), "")
	chain.AddBlock(json.Marshal(listingBlock))
}

func BuyListing(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Buying function called")
	//buyingBlock := InitBlockStruct()
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	fmt.Println(name)
	userBlock := InitBlockStruct(0, "", "", name)
	chain.AddBlock(json.Marshal(userBlock))
	fmt.Println("ended")
}

func ConvertBlock(block *blockchain.Block) *BlockData {
	data := *BlockData{}
	json.Unmarshal(block.Data, data)
	return data
}

func GetBlocks(blockType int) []BlockData {
	var results []blockchain.Block
	iterator := chain.Iterator()
	for current := iterator.Next() != nil {
		currentStruct := ConvertBlock(current)
		if currentStruct.BlockType == blockType {
			results = append(results, currentStruct)
		}
	}
}

//These functions retrieve data from the blockchain in order to display in the web app
func ListingPage(w http.ResponseWriter, r *http.Request) {
	//needs to be able to pull some data from the blockchain and display it here
	//traverse blockchain, decode it, and grab an array of the listing blocks
	listings := GetBlocks(1)
	//pass these to the testlistings page after adjusting it
	t, _ := template.ParseFiles("testListing.html")
	t.Execute(w, nil)
}

func PurchasePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("testPurchase.html")
	t.Execute(w, nil)
}
