package btree

import (
	"fmt"
	"testing"
)

func Test_BTree(t *testing.T) {
	root := NewBTree([]int{41, 44, 96, 46, 42, 20, 43})
	fmt.Println(root)
	//fmt.Println(root.Search(23))
	//fmt.Println(root.Search(45))
	//fmt.Println(root.Search(56))
}
