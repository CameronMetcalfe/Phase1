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

func SetUpNetwork(num int) Network {
	net := Network{connected: make([]bool, num)}
	nodes := make([]*Node, num)
	for i := 0; i < num; i++ {
		nodes[i] = InitNode(i, fmt.Sprintf("test%d", i), &net)
	}
	net.nodes = nodes
	return net
}

func CheckNetworkChains(net Network, sameLength bool) bool {
	numNodes := len(net.nodes)
	shortestChain := net.nodes[0].chainLength
	for i := 1; i < numNodes; i++ {
		//make sure chains are the same length if sameLength was specified
		if sameLength && shortestChain != net.nodes[i].chainLength {
			return false
		}
		if net.nodes[i].chainLength < shortestChain {
			shortestChain = net.nodes[i].chainLength
		}
	}

	for i := 0; i < shortestChain; i++ {
		val := net.nodes[0].chainSeen[i]
		for j := 1; j < numNodes; j++ {
			if net.nodes[j].chainSeen[i] != val {
				return false
			}
		}
	}
	return true
}

func TestBasics(t *testing.T) {
	net := SetUpNetwork(3)
	net.nodes[0].CreateListing("testItem 1", 5, "just to test some stuff")
	net.nodes[1].CreateListing("testItem 2", 7, "just to test some stuff")
	net.nodes[2].CreateListing("testItem 3", 2, "just to test some stuff")
	t.Logf("3 listings added, one for each node\n")

	//check stuff
	createPassed := CheckNetworkChains(net, true)
	if !createPassed {
		t.Logf("Node chains aren't in sync after adding listings\n")
		t.Fail()
	}

	net.nodes[1].PurchaseListing(0)
	net.nodes[0].PurchaseListing(2)

	purchasePassed := CheckNetworkChains(net, true)
	if !purchasePassed {
		t.Logf("Node chains aren't in sync after adding purchases\n")
		t.Fail()
	}
}
