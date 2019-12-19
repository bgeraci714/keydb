package rbtree

import (
	"bytes"
	"fmt"
)

// RBTree implementation
type RBTree struct {
	Root    *Node
	Compare func(a, b interface{}) int
}

// Get returns the value found for a given key, using an a boolean indicator if found
func (t *RBTree) Get(key string) (interface{}, bool) {
	node, found := t.GetNode(key)
	if found {
		return node.Value, true
	}
	return nil, false
}

// GetNode returns the node found for a given key, using a boolean indicator if found
func (t *RBTree) GetNode(key string) (*Node, bool) {
	return getNodeRec(t.Root, key, t.Compare)
}

func getNodeRec(n *Node, key string, compare func(a, b interface{}) int) (*Node, bool) {
	// node not found
	if n == nil {
		return nil, false
	}

	cmp := compare(key, n.Key)
	switch {
	case cmp > 0: // key > n.Key
		return getNodeRec(n.Right, key, compare)
	case cmp < 0: // key < n.Key
		return getNodeRec(n.Left, key, compare)
	default: // node was found
		return n, true
	}
}

func recolor(nodes []*Node) {
	for _, n := range nodes {
		n.Color = !n.Color
	}
}

// Translated from https://www.cs.auckland.ac.nz/software/AlgAnim/red_black.html
func leftRotate(t *RBTree, x *Node) {
	y := x.Right

	// turn y's left sub-tree into x's right sub-tree
	x.Right = y.Left
	if y.Left != nil {
		y.Left.Parent = x
	}

	// y's new parent was x's parent
	y.Parent = x.Parent

	// set the parent to point to y instead of x
	// need to check first if we're at the root
	if x.Parent == nil {
		t.Root = y
	} else if x == x.Parent.Left {
		// x is on left of its parent
		x.Parent.Left = y
	} else { // otherwise x must have been on the right
		x.Parent.Right = y
	}

	// put x on y's left
	y.Left = x
	x.Parent = y
}

func rightRotate(t *RBTree, x *Node) {
	y := x.Left

	// turn y's left sub-tree into x's left sub-tree
	x.Left = y.Right
	if y.Right != nil {
		y.Right.Parent = x
	}

	// y's new parent was x's parent
	y.Parent = x.Parent

	// set the parent to point to y instead of x
	// need to check first if we're at the root
	if x.Parent == nil {
		t.Root = y
	} else if x == x.Parent.Right {
		// x is on right of its parent
		x.Parent.Right = y
	} else { // otherwise x must have been on the left
		x.Parent.Left = y
	}

	// put x on y's right
	y.Right = x
	x.Parent = y
}

// Insert inserts with balanced algorithm involved
// Translated from https://www.cs.auckland.ac.nz/software/AlgAnim/red_black.html
func (t *RBTree) Insert(key string, val interface{}) {
	// Perform tree insert for tree T and node n
	t.insert(key, val)
	n, _ := t.GetNode(key) // will be optimized later

	n.Color = Red
	for n != t.Root && n.Parent.Color == Red {
		if n.Parent == n.LeftUncle() {
			// if n's parent is a left, the uncle is x's right uncle
			uncle := n.RightUncle()
			if uncle != nil && uncle.Color == Red {
				// Do case 1: change the colors
				n.Parent.Color = Black
				uncle.Color = Black
				n.Grandparent().Color = Red

				// move n up the tree
				n = n.Grandparent()
			} else {
				// uncle is a black node
				if n == n.Parent.Right {
					// and n is to the right
					// case 2 - move x up and rotate
					n = n.Parent
					leftRotate(t, n)
				}
				// case 3
				n.Parent.Color = Black
				n.Grandparent().Color = Red
				rightRotate(t, n.Grandparent())
			}
		} else {
			// if n's parent is a right, the uncle is x's left uncle
			uncle := n.LeftUncle()
			if uncle != nil && uncle.Color == Red {
				// Do case 1: change the colors
				n.Parent.Color = Black
				uncle.Color = Black
				n.Grandparent().Color = Red

				// move n up the tree
				n = n.Grandparent()
			} else {
				// uncle is a black node
				if n == n.Parent.Left {
					// and n is to the left
					// case 2 - move x up and rotate
					n = n.Parent
					rightRotate(t, n)
				}
				// case 3
				n.Parent.Color = Black
				n.Grandparent().Color = Red
				leftRotate(t, n.Grandparent())
			}
		}
	}
	t.Root.Color = Black
}

// Insert adds a new key value pair to the tree
func (t *RBTree) insert(key string, val interface{}) {
	t.Root = insertRec(t.Root, key, val, t.Compare, nil)
}

// Size returns the size of the tree
func (t RBTree) Size() int {
	return sizeRec(t.Root)
}

// Height returns the height of the tree
func (t RBTree) Height() int {
	return height(t.Root)
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func height(n *Node) int {
	if n == nil {
		return -1
	}
	return 1 + max(height(n.Left), height(n.Right))
}

func sizeRec(n *Node) int {
	if n == nil {
		return 0
	}
	return 1 + sizeRec(n.Left) + sizeRec(n.Right)
}

func insertRec(n *Node, key string, val interface{}, compare func(a, b interface{}) int, parent *Node) *Node {
	if n == nil {
		return &Node{key, val, nil, nil, false, parent} // might be an issue if this is allocated on function's stack, be mindful of this
	}

	cmp := compare(key, n.Key)
	switch {
	case cmp > 0: // key > n.Key
		n.Right = insertRec(n.Right, key, val, compare, n)
	case cmp < 0: // key < n.Key
		n.Left = insertRec(n.Left, key, val, compare, n)
	default: // key == n.Key
		n.Value = val // overwrite old value if the keys match
	}
	return n
}

// PrintInorder prints out all the nodes of a subtree inorder
func PrintInorder(n *Node) {
	if n != nil {
		PrintInorder(n.Left)
		fmt.Print(n.Key)
		PrintInorder(n.Right)
	}
}

// Delete deletes item with the matching key
func (t *RBTree) Delete(key string) {
	t.Root = delete(t.Root, key, t.Compare)
}

func delete(n *Node, key string, compare func(a, b interface{}) int) *Node {
	if n == nil {
		return nil
	}

	if cmp := compare(key, n.Key); cmp > 0 { // key > n.Key
		n.Right = delete(n.Right, key, compare)
	} else if cmp < 0 { // key < n.Key
		n.Left = delete(n.Left, key, compare)
	} else { // key == n.Key
		if n.Right == nil { // if no right child
			return n.Left
		} else if n.Left == nil { // if there's a right but no left
			return n.Right
		}

		// both right and left child
		tmp := n                       // copy over node n
		n = min(tmp.Right)             // swap n for its min on the right
		n.Right = deleteMin(tmp.Right) // replace the min that was just copied
		n.Left = tmp.Left              // copy over original left node
	}
	return n
}

func min(n *Node) *Node {
	if n.Left == nil {
		return n
	}
	return min(n.Left)
}

func deleteMin(n *Node) *Node {
	if n.Left == nil {
		return n.Right
	}
	n.Left = deleteMin(n.Left)
	return n
}

// ToString prints out a dash spaced version of the tree
func (t RBTree) ToString() string {
	return printSubtree(t.Root, 0)
}

func printSubtree(n *Node, h int) string {
	if n == nil {
		return ""
	}
	s := ""
	for i := 0; i < h; i++ {
		s += "-"
	}
	s += n.Key + "\n"
	s += printSubtree(n.Left, h+1)
	s += printSubtree(n.Right, h+1)
	return s
}

// ToMap converts tree to a map
func (t RBTree) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	addToMap(t.Root, &m)
	return m
}

func addToMap(n *Node, m *map[string]interface{}) {
	if n == nil {
		return
	}
	addToMap(n.Left, m)
	addToMap(n.Right, m)
	(*m)[n.Key] = n.Value
}

// MarshalBinary marshals the tree into a byte format
func (t RBTree) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	// if t.Root.MarshalBinary()
	// fmt.Fprintf(&b, )
	return b.Bytes(), nil
}

// func marshalInorder(n *Node) bytes.Buffer {
// 	var b bytes.Buffer
// 	if n == nil {
// 		return b
// 	}
// 	if bytes, err := b.WriteString(n.Key); err != nil {

// 	}

// }

// PrintInorder prints out all the nodes of a subtree inorder
// func PrintInorder(n *Node) {
// 	if n != nil {
// 		PrintInorder(n.Left)
// 		fmt.Print(n.Key)
// 		PrintInorder(n.Right)
// 	}
// }
