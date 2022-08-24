package datastructure

import (
	"testing"

	"github.com/duke-git/lancet/v2/internal"
)

func TestHashMap_PutAndGet(t *testing.T) {
	assert := internal.NewAssert(t, "TestHashMap_PutAndGet")

	hm := NewHashMap()

	hm.Put("abc", 3)
	assert.Equal(3, hm.Get("abc"))
	assert.IsNil(hm.Get("abcd"))

	hm.Put("abc", 4)
	assert.Equal(4, hm.Get("abc"))
}

func TestHashMap_Delete(t *testing.T) {
	assert := internal.NewAssert(t, "TestHashMap_Delete")

	hm := NewHashMap()

	hm.Put("abc", 3)
	assert.Equal(3, hm.Get("abc"))

	hm.Delete("abc")
	assert.IsNil(hm.Get("abc"))
}

func TestHashMap_Contains(t *testing.T) {
	assert := internal.NewAssert(t, "TestHashMap_Contains")

	hm := NewHashMap()
	assert.Equal(false, hm.Contains("abc"))

	hm.Put("abc", 3)
	assert.Equal(true, hm.Contains("abc"))
}