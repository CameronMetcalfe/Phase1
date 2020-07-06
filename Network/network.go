package network

//for RPC and network stuff
//currently no simulation of network lag/loss

type Node struct {
	//need array of blocks for the blockchain
	//need similar stuff to raft nodes (peers, etc)
	me    int
	peers []*Node
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
