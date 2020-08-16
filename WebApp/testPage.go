package WebApp

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	Bchain "../blockchain/blockchain"
	// "../network"
)

type BlockData struct { //will be converted to byte[] when added to actual blockchain
	ItemName        string
	ItemDescription string
	ItemPrice       int
	Username        string
	BlockType       int //0: user, 1: listing, 2: purchase
	PurchasedHash   []byte
	PurchasedIndex  int //used for test network, not actually used in the blockchain
}

//var nodes []*network.Node
// var chain *Bchain.BlockChain
// var UserAddress string

func (n *Node) Test() {
	fmt.Printf("Node: %d is testing\n", n.me)
}

func (n *Node) InitWebApp(address string, port int) { //should make this a function of Node
	//need to set up Nodes here
	//maybe take in number of nodes from command line arg
	//numNodes := 5
	//var chain blockchain.Blockchain
	fmt.Println("server is starting")
	if Bchain.DBexists() {
		n.chain = Bchain.ContinueBlockChain(address) //does continue work if it's a new address for existing DB?
	} else {
		n.chain = Bchain.InitBlockchain(address)
	}
	n.userAddress = address

	http.HandleFunc("/", n.HelloServer)
	http.HandleFunc("/add-listing", n.ListingPage)
	http.HandleFunc("/view-listings", n.CurrentListings)
	http.HandleFunc("/submit-listing", n.AddListing)
	http.HandleFunc("/purchase", n.PurchasePage)
	http.HandleFunc("/bought", n.BuyListing)
	http.HandleFunc("/add-user", n.AddUser)
	http.HandleFunc("/users", n.UsersPage)

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

func (n *Node) HelloServer(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	//title := r.URL.Path[1:] //r.URL.Path[]
	// testBlock := Block{}
	// testBlock.Hash = 3

	t, _ := template.ParseFiles("testPage.html")
	t.Execute(w, nil)
}

//These functions put data into the blockchain in response to actions in the web app
func (n *Node) AddListing(w http.ResponseWriter, r *http.Request) {
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
	n.chain.AddBlock(transactions)

	//show confirmation page
	t, _ := template.ParseFiles("confirmation.html")
	t.Execute(w, nil)
}

func (n *Node) BuyListing(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Buying function called")
	//buyingBlock := InitBlockStruct()
}

func (n *Node) AddUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	fmt.Println(name)
	userBlock := InitBlockStruct(0, "", "", 0, name)
	encoded, _ := json.Marshal(userBlock)
	transactions := []*Bchain.Transaction{&Bchain.Transaction{Inputs: []Bchain.TXInput{Bchain.TXInput{Sig: string(encoded)}}}}
	n.chain.AddBlock(transactions)
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

func (n *Node) GetBlocks(blockType int) []*BlockData {
	var results []*BlockData
	iterator := n.chain.ChainIter()

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

	return results
}

//These functions retrieve data from the blockchain in order to display in the web app
func (n *Node) ListingPage(w http.ResponseWriter, r *http.Request) {
	//pass these to the testlistings page after adjusting it
	t, _ := template.ParseFiles("addListing.html")
	t.Execute(w, nil)
}

func (n *Node) CurrentListings(w http.ResponseWriter, r *http.Request) {
	//needs to be able to pull some data from the blockchain and display it here
	//traverse blockchain, decode it, and grab an array of the listing blocks
	listings := n.GetBlocks(1)
	fmt.Printf("Retrieved %d listings\n", len(listings))
	t, _ := template.ParseFiles("viewListings.html")
	t.Execute(w, listings)
	//Printchain()
}

func (n *Node) PurchasePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("testPurchase.html")
	t.Execute(w, nil)
}

func (n *Node) UsersPage(w http.ResponseWriter, r *http.Request) {
	users := n.GetBlocks(0)
	fmt.Println(len(users))
	t, _ := template.ParseFiles("testUsers.html")
	t.Execute(w, nil)
}

func (n *Node) Printchain() {
	//chain := Blockchain.ContinueBlockChain("")
	//defer chain.Database.Close()
	iter := n.chain.ChainIter()

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
