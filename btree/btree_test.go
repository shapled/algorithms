package btree

import (
	"fmt"
	"testing"
)

func Test_BTree(t *testing.T) {
	root := NewBTree([]int{23, 1, 44, 56})
	fmt.Println(root)
	fmt.Println(root.Search(23))
	fmt.Println(root.Search(45))
	fmt.Println(root.Search(56))
}
