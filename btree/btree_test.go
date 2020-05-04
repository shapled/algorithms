package btree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BTree(t *testing.T) {
	root := NewBTree(41, 44, 96, 46, 42, 20, 43, 3, 77, 99)
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
	var deleted bool
	deleted, root = root.Delete(41)
	assert.True(t, deleted)
	_, root = root.Delete(43)
	_, root = root.Delete(46)
	_, root = root.Delete(77)
	_, root = root.Delete(20)
	_, root = root.Delete(3)
	deleted, root = root.Delete(7)
	assert.False(t, deleted)
	_, root = root.Delete(44)
	_, root = root.Delete(99)
	_, root = root.Delete(42)
	_, root = root.Delete(96)
	assert.Equal(t, 0, root.KeyCount)
}
