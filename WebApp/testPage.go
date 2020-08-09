package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	Bchain "../blockchain"
)

type BlockData struct { //will be converted to byte[] when added to actual blockchain
	ItemName        string
	ItemDescription string
	ItemPrice       int
	Username        string
	BlockType       int
}

//var nodes []*network.Node
var chain *Bchain.BlockChain

func main() {
	//need to set up Nodes here
	//maybe take in number of nodes from command line arg
	//numNodes := 5
	//var chain blockchain.Blockchain
	fmt.Println("server is starting")
	chain = Bchain.InitBlockchain()

	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/list", ListingPage)
	http.HandleFunc("/add-listing", AddListing)
	http.HandleFunc("/purchase", PurchasePage)
	http.HandleFunc("/bought", BuyListing)
	http.HandleFunc("/add-user", AddUser)
	http.HandleFunc("/users", UsersPage)

	http.ListenAndServe(":8080", nil)
}

func InitBlockStruct(blockType int, itemName string, desc string, price int, username string) *BlockData {
	result := &BlockData{}
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
	// testBlock := Block{}
	// testBlock.Hash = 3

	t, _ := template.ParseFiles("testPage.html")
	t.Execute(w, nil)
}

//These functions put data into the blockchain in response to actions in the web app
func AddListing(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add listing function called")
	var price string
	price = r.FormValue("price")
	var priceInt int
	priceInt, _ = strconv.Atoi(price)
	listingBlock := InitBlockStruct(1, r.FormValue("item-name"), r.FormValue("item-description"), priceInt, "")
	encoded, _ := json.Marshal(listingBlock)
	chain.AddBlock(string(encoded))

	//show confirmation page
	t, _ := template.ParseFiles("confirmation.html")
	t.Execute(w, nil)
}

func BuyListing(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Buying function called")
	//buyingBlock := InitBlockStruct()
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	fmt.Println(name)
	userBlock := InitBlockStruct(0, "", "", 0, name)
	encoded, _ := json.Marshal(userBlock)
	chain.AddBlock(string(encoded))
	fmt.Println("ended")

	//show confirmation page
	t, _ := template.ParseFiles("confirmation.html")
	t.Execute(w, nil)
}

func ConvertBlock(block *Bchain.Block) BlockData {
	data := BlockData{}
	json.Unmarshal(block.Data, data)
	return data
}

func GetBlocks(blockType int) []BlockData {
	var results []BlockData
	iterator := chain.ChainIter()
	fmt.Println(iterator)
	current := iterator.Next()
	for current != nil {
		currentStruct := ConvertBlock(current)
		if currentStruct.BlockType == blockType {
			results = append(results, currentStruct)
		}
		current = iterator.Next()
	}
	return results
}

//These functions retrieve data from the blockchain in order to display in the web app
func ListingPage(w http.ResponseWriter, r *http.Request) {
	//needs to be able to pull some data from the blockchain and display it here
	//traverse blockchain, decode it, and grab an array of the listing blocks
	//listings := GetBlocks(1)
	//pass these to the testlistings page after adjusting it
	t, _ := template.ParseFiles("testListing.html")
	t.Execute(w, nil)
}

func PurchasePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("testPurchase.html")
	t.Execute(w, nil)
}

func UsersPage(w http.ResponseWriter, r *http.Request) {
	users := GetBlocks(0)
	fmt.Println(len(users))
	t, _ := template.ParseFiles("testUsers.html")
	t.Execute(w, nil)
}
