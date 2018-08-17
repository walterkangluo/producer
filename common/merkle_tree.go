package common

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"github.com/DSiSc/txpool/common"
	"io"
)

type merkleTreeNode struct {
	Hash  common.Hash
	Left  *merkleTreeNode
	Right *merkleTreeNode
}

type merkleTree struct {
	Depth uint
	Root  *merkleTreeNode
}

func Serialize(w io.Writer, u *common.Hash) error {
	_, err := w.Write(u[:])
	return err
}

func ComputeMerkleRoot(hashes []common.Hash) common.Hash {
	if len(hashes) == 0 {
		return common.Hash{}
	}
	if len(hashes) == 1 {
		return hashes[0]
	}
	tree, _ := newMerkleTree(hashes)
	return tree.Root.Hash
}

//Generate the leaves nodes
func generateLeaves(hashes []common.Hash) []*merkleTreeNode {
	var leaves []*merkleTreeNode
	for _, d := range hashes {
		node := &merkleTreeNode{
			Hash: d,
		}
		leaves = append(leaves, node)
	}
	return leaves
}

//calc the next level's hash use double sha256
func levelUp(nodes []*merkleTreeNode) []*merkleTreeNode {
	var nextLevel []*merkleTreeNode
	for i := 0; i < len(nodes)/2; i++ {
		var data []common.Hash
		data = append(data, nodes[i*2].Hash)
		data = append(data, nodes[i*2+1].Hash)
		hash := doubleSha256(data)
		node := &merkleTreeNode{
			Hash:  hash,
			Left:  nodes[i*2],
			Right: nodes[i*2+1],
		}
		nextLevel = append(nextLevel, node)
	}
	if len(nodes)%2 == 1 {
		var data []common.Hash
		data = append(data, nodes[len(nodes)-1].Hash)
		data = append(data, nodes[len(nodes)-1].Hash)
		hash := doubleSha256(data)
		node := &merkleTreeNode{
			Hash:  hash,
			Left:  nodes[len(nodes)-1],
			Right: nodes[len(nodes)-1],
		}
		nextLevel = append(nextLevel, node)
	}
	return nextLevel
}

func doubleSha256(s []common.Hash) common.Hash {
	b := new(bytes.Buffer)
	for _, d := range s {
		Serialize(b, &d)
	}
	temp := sha256.Sum256(b.Bytes())
	f := sha256.Sum256(temp[:])
	return common.Hash(f)
}

func newMerkleTree(hashes []common.Hash) (*merkleTree, error) {
	if len(hashes) == 0 {
		return nil, errors.New("NewMerkleTree input no item error.")
	}
	var height uint

	height = 1
	nodes := generateLeaves(hashes)
	for len(nodes) > 1 {
		nodes = levelUp(nodes)
		height += 1
	}
	mt := &merkleTree{
		Root:  nodes[0],
		Depth: height,
	}
	return mt, nil
}
