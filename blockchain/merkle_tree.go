package blockchain

// MerkleTree represents a Merkle tree
type MerkleTree struct {
	RootNode *MerkleNode
}

// MerkleNode represents a Merkle tree node
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}
