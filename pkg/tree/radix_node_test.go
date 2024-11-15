package tree

import (
	"testing"
)

func TestRadix(t *testing.T) {
	a := &RadixNode[string]{}

	r := NewResult[string]("test")

	a.Insert("test", r)
	a.Insert("test123", r)
	a.Insert("123test", r)

	j := a
	for _, v := range convert("test123", radixSize) {
		t.Log(j, j.next)
		j = j.children[v]
	}

	j = a
	for _, v := range convert("123test", radixSize) {
		t.Log(j, j.next)
		j = j.children[v]
	}
}

func TestRadixA(t *testing.T) {
	a := &RadixNode[string]{}

	a.Insert(123, NewResult[string]("test1"))

	j := a
	for _, v := range convert(123, radixSize) {
		t.Log(v, j, j.next)
		j = j.children[v]
	}
	t.Log(j)
}

func TestRadixMerge(t *testing.T) {
	a := &RadixNode[string]{}
	b := &RadixNode[string]{}

	a.Insert("a", NewResult[string]("1"))
	a.Insert("b", NewResult[string]("2"))
	a.Insert("c", NewResult[string]("3"))
	b.Insert("x", NewResult[string]("4"))
	b.Insert("y", NewResult[string]("5"))
	b.Insert("z", NewResult[string]("6"))

	t.Log(a.children)

	a.Merge(b)

	t.Log(a.children)
}
