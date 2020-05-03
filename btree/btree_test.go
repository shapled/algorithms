package btree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BTree(t *testing.T) {
	root := NewBTree(41, 44, 96, 46, 42, 20, 43)
	fmt.Println(root)
	var ok bool
	ok, _, _ = root.Search(23)
	assert.False(t, ok)
	ok, _, _ = root.Search(45)
	assert.False(t, ok)
	ok, _, _ = root.Search(20)
	assert.True(t, ok)
	ok, _, _ = root.Search(56)
	assert.False(t, ok)
}
