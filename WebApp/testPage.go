package WebApp

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	Bchain "../go_blockchain/blockchain"
	// "../network"
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
var UserAddress string

func (n *Node) Test() {
	fmt.Printf("Node: %d is testing\n", n.me)
}

func InitWebApp(address string, port int) { //should make this a function of Node
	//need to set up Nodes here
	//maybe take in number of nodes from command line arg
	//numNodes := 5
	//var chain blockchain.Blockchain
	fmt.Println("server is starting")
	if Bchain.DBexists() {
		chain = Bchain.ContinueBlockChain(address) //does continue work if it's a new address for existing DB?
	} else {
		chain = Bchain.InitBlockchain(address)
	}
	UserAddress = address

	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/add-listing", ListingPage)
	http.HandleFunc("/view-listings", CurrentListings)
	http.HandleFunc("/submit-listing", AddListing)
	http.HandleFunc("/purchase", PurchasePage)
	http.HandleFunc("/bought", BuyListing)
	http.HandleFunc("/add-user", AddUser)
	http.HandleFunc("/users", UsersPage)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
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

	//block added for testing
	fmt.Printf("For the new item; name: %s, price: %d, desc: %s\n", listingBlock.ItemName, listingBlock.ItemPrice, listingBlock.ItemDescription)

	encoded, _ := json.Marshal(listingBlock)
	transactions := []*Bchain.Transaction{&Bchain.Transaction{Inputs: []Bchain.TXInput{Bchain.TXInput{Sig: string(encoded)}}}}
	chain.AddBlock(transactions)

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
	transactions := []*Bchain.Transaction{&Bchain.Transaction{Inputs: []Bchain.TXInput{Bchain.TXInput{Sig: string(encoded)}}}}
	chain.AddBlock(transactions)
	fmt.Println("ended")

	//show confirmation page
	t, _ := template.ParseFiles("confirmation.html")
	t.Execute(w, nil)
}

func ConvertBlock(block *Bchain.Block) *BlockData {
	data := &BlockData{}
	err := json.Unmarshal([]byte(block.Transactions[0].Inputs[0].Sig), &data)
	if err != nil {
		//this block was a transaction, so couldn't unmarshal
		return nil
	}
	fmt.Printf("For the retrieved item; name: %s, price: %d, desc: %s\n", data.ItemName, data.ItemPrice, data.ItemDescription)
	return data
}

func GetBlocks(blockType int) []*BlockData {
	var results []*BlockData
	iterator := chain.ChainIter()

	for {
		current := iterator.Next()
		currentStruct := ConvertBlock(current)
		if currentStruct != nil && currentStruct.BlockType == blockType {
			results = append(results, currentStruct)
		}

		if len(current.PrevHash) == 0 {
			break
		}
	}

	// for current != nil { //may have to change this to break on prevHash == 0
	// 	currentStruct := ConvertBlock(current)
	// 	if currentStruct.BlockType == blockType {
	// 		results = append(results, currentStruct)
	// 	}
	// 	current = iterator.Next()
	// }
	return results
}

//These functions retrieve data from the blockchain in order to display in the web app
func ListingPage(w http.ResponseWriter, r *http.Request) {
	//pass these to the testlistings page after adjusting it
	t, _ := template.ParseFiles("addListing.html")
	t.Execute(w, nil)
}

func CurrentListings(w http.ResponseWriter, r *http.Request) {
	//needs to be able to pull some data from the blockchain and display it here
	//traverse blockchain, decode it, and grab an array of the listing blocks
	listings := GetBlocks(1)
	fmt.Printf("Retrieved %d listings\n", len(listings))
	t, _ := template.ParseFiles("viewListings.html")
	t.Execute(w, listings)
	//Printchain()
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

func Printchain() {
	//chain := Blockchain.ContinueBlockChain("")
	//defer chain.Database.Close()
	iter := chain.ChainIter()

	for {
		block := iter.Next()
		fmt.Printf("previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := Bchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}
