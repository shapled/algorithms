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
*/

const (
	M	= 10
)

type Node struct {
	Keys		[]int		// Max Length: M-1
	Children	[]*Node		// Max Length: M
	Parent		*Node
	Dup			[]*Node		// 重复节点，保证树本身无重复节点，可以最快的速度索引
}

func NewBTree(values []int) *Node {
	if len(values) == 0 {
		return nil
	}
	root := &Node{Keys:[]int{values[0]}}
	for _, value := range values[1:] {
		root.Insert(value)
	}
	return root
}

func (root *Node) Search(value int) (bool, *Node) {
	if root == nil {
		return false, &Node{}
	}
	for i, key := range root.Keys {
		if value == key {
			return true, root
		} else if value < key {
			if i < len(root.Children) {
				return root.Children[i].Search(value)
			} else {
				return false, root
			}
		}
	}
	if len(root.Children) > len(root.Keys) {
		return root.Children[len(root.Children)-1].Search(value)
	}
	return false, root
}

func (root *Node) Insert(value int) *Node{
	ok, node := root.Search(value)
	if ok {
		newNode := &Node{
			Keys:     []int{value},
			Children: nil,
			Parent:   nil,
			Dup:      nil,
		}
		node.Dup = append(node.Dup, newNode)
		return newNode
	}
	node.Keys = append(node.Keys, value)
	return node
}

func (root *Node) String() string {
	nodes			:= []*Node{root}
	nextLayerNodes	:= make([]*Node, 0, 0)
	outputTexts		:= make([]string, 0, 0)
	for len(nodes) != 0 {
		layerTexts := make([]string, 0, 0)
		for _, node := range nodes {
			nodeTexts := make([]string, 0, 0)
			for i:=0; i<len(node.Keys) || i<len(node.Children); i++ {
				if i < len(node.Children) {
					child := node.Children[i]
					nodeTexts = append(nodeTexts, "o")
					nextLayerNodes = append(nextLayerNodes, child)
				}
				if i < len(node.Keys) {
					nodeTexts = append(nodeTexts, strconv.Itoa(node.Keys[i]))
				}
			}
			layerTexts = append(layerTexts, "(" + strings.Join(nodeTexts, ", ") + ")")
		}
		nodes = nextLayerNodes
		outputTexts = append(outputTexts, strings.Join(layerTexts, " "))
	}
	return strings.Join(outputTexts, "\n")
}
