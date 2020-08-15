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

func (n *Node) SendMessage(to int, args *MessageArgs, reply *MessageReply) bool {
	if !n.network.connected[to] {
		return false
	}

	n.network.nodes[to].VerifyMessage(args, reply)
	return true
}

func (n *Node) VerifyMessage(args *MessageArgs, reply *MessageReply) {
	reply.PeerChainLength = n.chainLength
	//reject the message if it's behind in the chain
	if args.CurrentChainLength < n.chainLength {
		reply.Accepted = false
		return
	}

	reply.Accepted = true
}

func (n *Node) SuggestNewBlock(block *BlockData) bool {
	numVotes := 0
	maxReceivedChainLength := 0
	longestChainNode := 0
	failed := false
	for i, _ := range n.network.nodes {
		if i != n.me && n.network.connected[i] {
			args := &MessageArgs{InteractionType: block.BlockType, CurrentChainLength: n.chainLength, block: block}
			reply := &MessageReply{}
			n.SendMessage(i, args, reply)

			if reply.Accepted {
				numVotes++
			} else {
				failed = true
				if reply.PeerChainLength > maxReceivedChainLength {
					maxReceivedChainLength = reply.PeerChainLength
					longestChainNode = 1
				}
			}

			if reply.PeerChainLength < n.chainLength {
				n.SendMissingBlocks(reply.PeerChainLength, i)
			}
		}
	}

	if failed { //failed because chain is out of date, update it and try again
		n.network.nodes[longestChainNode].SendMissingBlocks(n.chainLength, n.me)
		return n.SuggestNewBlock(block)
	} else if numVotes < (len(n.network.nodes)) { //failed because of a network partition, return false. User should be alerted to try again later
		return false
	} else {
		return true
	}
}

func (n *Node) ConfirmBlock(block *BlockData) {
	for i, _ := range n.network.nodes {
		if i != n.me {
			n.SendMissingBlocks(n.chainLength-1, i)
		}
	}
}

func (n *Node) CreateListing(itemName string, itemPrice int, itemDescription string) bool {
	block := InitBlockStruct(1, itemName, itemDescription, itemPrice, n.userAddress)
	passed := n.SuggestNewBlock(block)

	if passed {
		//send out update messages to all the nodes
		n.chainLength++
		n.chainSeen = append(n.chainSeen, block)
		n.ConfirmBlock(block)
		return true
	}
	return false
}

func (n *Node) PurchaseListing(itemNum int) bool {
	block := &BlockData{BlockType: 2, PurchasedIndex: itemNum}
	passed := n.SuggestNewBlock(block)

	if passed {
		//handle stuff
		n.chainLength++
		n.chainSeen = append(n.chainSeen, block)
		n.ConfirmBlock(block)
		return true
	}
	return false
}

func (n *Node) SendMissingBlocks(start int, to int) bool {
	if !n.network.connected[to] {
		return false
	}

	blocksToSend := make([]*BlockData, n.chainLength-start+1)
	for i := start; i < n.chainLength; i++ {
		blocksToSend[i-start] = n.chainSeen[i]
	}
	if n.network.connected[to] {
		//call method for handling adding existing blocks
		n.network.nodes[to].AddSentBlocks(blocksToSend)
	}
	return true
}

func (n *Node) AddSentBlocks(sentBlocks []*BlockData) {
	n.chainSeen = append(n.chainSeen, sentBlocks...)
	n.chainLength += len(sentBlocks)
}

func InitNode(me int, username string, net *Network) *Node {
	createdNode := Node{me: me, userAddress: username, network: net}
	createdNode.chainSeen = make([]*BlockData, 0)
	return &createdNode
}
