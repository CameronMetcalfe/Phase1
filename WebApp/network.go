package WebApp

//for RPC and network stuff
//currently no simulation of network lag/loss
import (
	Bchain "../go_blockchain/blockchain"
)

type Network struct {
	nodes     []*Node
	connected []bool
	nodeNames map[string]int
}

type Node struct {
	//need array of blocks for the blockchain
	//need similar stuff to raft nodes (peers, etc)
	me          int
	peers       []*Node //should this be a map and use useraddress?
	chain       *Bchain.BlockChain
	chainSeen   []*BlockData
	network     *Network
	chainLength int
	userAddress string
}

type MessageArgs struct {
	InteractionType    int
	CurrentChainLength int
	block              *BlockData
}

type MessageReply struct {
	Accepted        bool
	PeerChainLength int
}

//checks that values held by nodes match up
func (net *Network) CheckChains() {

}

func (n *Node) SendMessage(to int, args *MessageArgs, reply *MessageReply) {

}

func (n *Node) VerifyMessage(args *MessageArgs, reply *MessageReply) {
	reply.PeerChainLength = n.chainLength
	//reject the message if it's behind in the chain
	if args.CurrentChainLength < n.chainLength {
		reply.Accepted = false
		return
	}

	reply.Accepted = true
	// if args.CurrentChainLength > n.chainLength {

	// }
}

func (n *Node) CreateListing(itemName string, itemPrice int, itemDescription string) {
	block := InitBlockStruct(1, itemName, itemDescription, itemPrice, n.userAddress)
	numVotes := 0
	//maxReceivedChainLength := 0
	for i, _ := range n.network.nodes {
		if i != n.me && n.network.connected[i] {
			args := &MessageArgs{InteractionType: 1, CurrentChainLength: n.chainLength, block: block}
			reply := &MessageReply{}
			n.SendMessage(i, args, reply)

			if reply.Accepted {
				numVotes++
			} else {
				//maxReceivedChainLength
			}
		}
	}
}

func (n *Node) VerifyNewUser(username string) bool { //should take in
	//should check the blockchain for a user with the same username, if none is found return true
	return true
}

func (n *Node) VerifyNewListing() bool { //not sure if this is necessary
	return true
}

func (n *Node) VerifyNewPurchase() bool { //definitely necessary, must ensure product is still available and buyer is paying the same price as listed by the seller
	return true
}
