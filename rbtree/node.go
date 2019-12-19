package rbtree

// Color is a typedef for booleans
type Color bool

// Red is a basic Coloring constant equaling true
const Red = true

// Black is a basic Coloring constant equaling false
const Black = false

// Node are nodes of base BST
type Node struct {
	Key    string
	Value  interface{}
	Left   *Node
	Right  *Node
	Color  Color
	Parent *Node
}

// Compare compares the Keys of the two nodes
// returns 1 if n's is greater than other's
// returns 0 if n's Key is equal to other's
// returns -1 if n's Key is less than other's
func (n Node) Compare(other Node) int {
	if n.Key > other.Key {
		return 1
	} else if n.Key < other.Key {
		return -1
	}
	return 0
}

// Grandparent returns the grandparent of the current node
func (n *Node) Grandparent() *Node {
	if n.Parent == nil || n.Parent.Parent == nil {
		panic("There is no grandParent!")
	}
	return n.Parent.Parent
}

// LeftUncle goes to find the Right uncle and panics if not
func (n *Node) LeftUncle() *Node {
	if n.Parent.Parent == nil {
		panic("There is no Left uncle because there is no grandParent!")
	}
	return n.Parent.Parent.Left
}

// RightUncle goes to find the Right uncle and panics if not
func (n *Node) RightUncle() *Node {
	if n.Parent.Parent == nil {
		panic("There is no Right uncle because there is no grandParent!")
	}
	return n.Parent.Parent.Right
}
