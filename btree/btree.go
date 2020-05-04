package btree

import (
	"strconv"
	"strings"
)

/*
M 阶 b-tree
1. 每个结点最多有 m-1 个关键字。
2. 根结点最少可以只有1个关键字。
3. 非根结点至少有 ceil(m/2)-1 个关键字。
4. 每个结点中的关键字都按照从小到大的顺序排列，每个关键字的左子树中的所有关键字都小于它，而右子树中的所有关键字都大于它。
5. 所有叶子结点都位于同一层，或者说根结点到每个叶子结点的长度都相同。

代码参考： https://blog.csdn.net/qq_36183935/article/details/80382490

感受：
1. 之前学的树算法都是想办法操作子节点，这个倒好，操作起父节点了，优秀！
*/

const M = 3

type (
	Node struct {
		IsLeaf		bool
		KeyCount	int
		Keys		[M]int			//	不得超过 M-1，最后一个位置是为了方便切分前写出
		Children	[M+1]*Node		//	不得超过 M，理由同上
		Parent		*Node
	}
)

func NewBTree(values ...int) *Node {
	tree := &Node{IsLeaf:true}
	for _, value := range values {
		tree = tree.Insert(value)
	}
	return tree
}

func (node *Node) Search(value int) (bool, *Node, int) {
	for i:=0; i<node.KeyCount; i++ {
		if value == node.Keys[i] {
			return true, node, i
		} else if value < node.Keys[i] {
			if node.IsLeaf {
				return false, node, i
			} else {
				return node.Children[i].Search(value)
			}
		}
	}
	if node.IsLeaf {
		return false, node, node.KeyCount
	} else {
		return node.Children[node.KeyCount].Search(value)
	}
}

func (node *Node) Root() *Node {
	for node.Parent != nil {
		node = node.Parent
	}
	return node
}

func (root *Node) Insert(value int) *Node {
	found, n, index := root.Search(value)
	if found {
		return root
	}
	for i:=n.KeyCount-1; i>=index; i-- {
		n.Keys[i+1] = n.Keys[i]
	}
	n.Keys[index] = value
	n.KeyCount += 1
	if n.KeyCount > M - 1 {
		root = n.Split()
	}
	return root
}

func (node *Node) Split() *Node {
	pos := (M - 1) / 2
	parent := node.Parent
	if parent == nil {
		parent = &Node{
			IsLeaf:   false,
			KeyCount: 0,
			Keys:     [M]int{},
			Children: [M+1]*Node{node},
			Parent:   nil,
		}
		node.Parent = parent
	}
	sibling := &Node{
		IsLeaf:   node.IsLeaf,
		KeyCount: 0,
		Keys:     [M]int{},
		Children: [M+1]*Node{},
		Parent:   parent,
	}
	posParent := 0
	for ; posParent<parent.KeyCount; posParent++ {
		if node.Keys[pos] < parent.Keys[posParent] {
			break
		}
	}
	for i:=parent.KeyCount-1; i>=posParent; i-- {
		parent.Keys[i+1] = parent.Keys[i]
		parent.Children[i+2] = parent.Children[i+1]
	}
	parent.Keys[posParent] = node.Keys[pos]
	parent.Children[posParent+1] = sibling
	parent.KeyCount += 1
	for i:=pos+1; i<=node.KeyCount; i++ {
		if !sibling.IsLeaf {
			sibling.Children[sibling.KeyCount] = node.Children[i]
			sibling.Children[sibling.KeyCount].Parent = sibling
		}
		if i != node.KeyCount {
			sibling.Keys[sibling.KeyCount] = node.Keys[i]
			sibling.KeyCount += 1
		}
	}
	node.KeyCount = pos
	if parent.KeyCount > M - 1 {
		return parent.Split()
	}
	for parent.Parent != nil {
		parent = parent.Parent
	}
	return parent
}

func (root *Node) Delete(value int) (bool, *Node) {
	found, n, index := root.Search(value)
	if !found {
		return false, root
	}
	leaf := n
	if !n.IsLeaf {
		leaf = n.Children[index+1]
		for !leaf.IsLeaf {
			leaf = leaf.Children[0]
		}
		n.Keys[index] = leaf.Keys[0]
		leaf.Keys[0] = value
		index = 0
	}
	for i:=index; i<leaf.KeyCount; i++ {
		leaf.Keys[i] = leaf.Keys[i+1]
	}
	leaf.KeyCount -= 1
	if leaf.KeyCount < (M - 1) / 2 {
		root = leaf.Merge(value)
	}
	return true, root
}

// base 用于当前节点没有 key 时确定此节点在父节点中的位置
func (node *Node) Merge(base int) *Node {
	parent := node.Parent
	if parent == nil {
		for node.KeyCount == 0 && !node.IsLeaf {
			node = node.Children[0]
			node.Parent = nil
		}
		return node
	}
	if node.KeyCount > 0 {
		base = node.Keys[0]
	}
	posParent := 0
	for ; posParent<parent.KeyCount; posParent++ {
		if base < parent.Keys[posParent] {
			break
		}
	}
	leftSibling := posParent - 1
	rightSibling := posParent + 1
	if leftSibling >= 0 && parent.Children[leftSibling].KeyCount <= (M - 1) / 2 {
		parent.mergeNode(leftSibling)
		if parent.KeyCount < (M - 1) / 2 {
			node = parent.Merge(base)
		}
	} else if rightSibling <= parent.KeyCount && parent.Children[rightSibling].KeyCount <= (M - 1) / 2 {
		parent.mergeNode(posParent)
		if parent.KeyCount < (M - 1) / 2 {
			node = parent.Merge(base)
		}
	} else if leftSibling >= 0 {
		left := parent.Children[leftSibling]
		for i:=node.KeyCount; i>0; i-- {
			node.Children[i+1] = node.Children[i]
			if i != node.KeyCount {
				node.Keys[i+1] = node.Keys[i]
			}
		}
		node.Keys[0] = parent.Keys[leftSibling]
		node.Children[0] = left.Children[left.KeyCount]
		if node.Children[0] != nil {
			node.Children[0].Parent = node
		}
		node.KeyCount += 1
		parent.Keys[leftSibling] = left.Keys[left.KeyCount-1]
		left.KeyCount -= 1
	} else if rightSibling <= parent.KeyCount {
		right := parent.Children[rightSibling]
		node.Keys[node.KeyCount] = parent.Keys[posParent]
		node.Children[node.KeyCount+1] = right.Children[0]
		if node.Children[node.KeyCount+1] != nil {
			node.Children[node.KeyCount+1].Parent = node
		}
		node.KeyCount += 1
		parent.Keys[posParent] = right.Keys[0]
		for i:=0; i<right.KeyCount; i++ {
			right.Children[i] = right.Children[i+1]
			if i != right.KeyCount - 1 {
				right.Keys[i] = right.Keys[i+1]
			}
		}
		right.KeyCount -= 1
	}
	return node.Root()
}

// merge left and left+1
func (node *Node) mergeNode(leftPos int) {
	left := node.Children[leftPos]
	right := node.Children[leftPos+1]
	left.Keys[left.KeyCount] = node.Keys[leftPos]
	for i:=0; i<=right.KeyCount; i++ {
		if i != right.KeyCount {
			left.Keys[left.KeyCount+1+i] = right.Keys[i]
		}
		if !right.IsLeaf {
			left.Children[left.KeyCount+1+i] = right.Children[i]
			left.Children[left.KeyCount+1+i].Parent = left
		}
	}
	left.KeyCount += 1 + right.KeyCount
	for i:=leftPos; i<node.KeyCount-1; i++ {
		node.Keys[i] = node.Keys[i+1]
		node.Children[i+1] = node.Children[i+2]
	}
	node.KeyCount -= 1
}

func (root *Node) String() string {
	nodes := []*Node{root}
	outputTexts := make([]string, 0, 0)
	for len(nodes) != 0 {
		nextLayerNodes := make([]*Node, 0, 0)
		layerTexts := make([]string, 0, 0)
		for _, node := range nodes {
			nodeTexts := make([]string, 0, 0)
			for i := 0; i <= node.KeyCount; i++ {
				if !node.IsLeaf {
					child := node.Children[i]
					if child != nil {
						nodeTexts = append(nodeTexts, "o")
						nextLayerNodes = append(nextLayerNodes, child)
					}
				}
				if i != node.KeyCount {
					nodeTexts = append(nodeTexts, strconv.Itoa(node.Keys[i]))
				}
			}
			if node.IsLeaf {
				layerTexts = append(layerTexts, "([leaf] "+strings.Join(nodeTexts, ", ")+")")
			} else {
				layerTexts = append(layerTexts, "("+strings.Join(nodeTexts, ", ")+")")
			}
		}
		nodes = nextLayerNodes
		outputTexts = append(outputTexts, strings.Join(layerTexts, " "))
	}
	return strings.Join(outputTexts, "\n")
}
