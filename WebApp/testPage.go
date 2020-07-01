package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/list", ListingPage)
	http.HandleFunc("/add-listing", AddListing)
	http.HandleFunc("/purchase", PurchasePage)
	http.HandleFunc("/bought", BuyListing)

	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	//title := r.URL.Path[1:] //r.URL.Path[]

	t, _ := template.ParseFiles("testPage.html")
	t.Execute(w, nil)
}

func ListingPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("testListing.html")
	t.Execute(w, nil)
}

func AddListing(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add listing function called")
}

func PurchasePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("testPurchase.html")
	t.Execute(w, nil)
}

func BuyListing(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Buying function called")
}
