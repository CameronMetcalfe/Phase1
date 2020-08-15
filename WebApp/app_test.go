package WebApp

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func ClearBlockChain() {
	err := os.RemoveAll("./tmp/blocks")
	if err != nil {
		log.Fatal(err)
	}
}

//this test just starts up the web app without doing any other actions
func TestStartFresh(t *testing.T) {
	n := Node{me: 0, peers: make([]*Node, 1), chainLength: 0}
	ClearBlockChain()
	n.InitWebApp("Cameron", 8080)
}

//starts up the web app and populates the blockchain
func TestCreate(t *testing.T) {

}

func TestMultipleNodes(t *testing.T) {
	n0 := Node{me: 0, chainLength: 0}
	n1 := Node{me: 1, chainLength: 0}
	nodes := []*Node{&n0, &n1}
	n0.peers = nodes
	n1.peers = nodes
	ClearBlockChain()

	n0.InitWebApp("Cameron", 8080)
	fmt.Println("Got here")
	//n1.InitWebApp("Test", 9090)
}
