package main

import (
	"fmt"
	"html/template"
	"net/http"

	"../network"
)

type Block struct {
	Hash            int
	ItemDescription string
	Username        string
	BlockType       int
}

var nodes []*network.Node

func main() {
	//need to set up Nodes here
	//maybe take in number of nodes from command line arg
	//numNodes := 5

	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/list", ListingPage)
	http.HandleFunc("/add-listing", AddListing)
	http.HandleFunc("/purchase", PurchasePage)
	http.HandleFunc("/bought", BuyListing)
	http.HandleFunc("/add-user", AddUser)

	http.ListenAndServe(":8080", nil)
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
}

func BuyListing(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Buying function called")
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	// for key, values := range r.Form {
	// 	for _, value := range values {
	// 		fmt.Println(key, value)
	// 	}
	// }
	name := r.FormValue("username")
	fmt.Println(name)
	fmt.Println("ended")
}

//These functions retrieve data from the blockchain in order to display in the web app
func ListingPage(w http.ResponseWriter, r *http.Request) {
	//needs to be able to pull some data from the blockchain and display it here
	t, _ := template.ParseFiles("testListing.html")
	t.Execute(w, nil)
}

func PurchasePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("testPurchase.html")
	t.Execute(w, nil)
}
