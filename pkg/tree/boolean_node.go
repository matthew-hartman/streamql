package tree

import (
	"crypto/sha256"
	"fmt"

	"github.com/itchyny/gojq"
)

type BooleanNode[T Input] struct {
	// false -> 0, true -> 1
	children [2]Node[T]
	check    *gojq.Code
	queryStr string
}

func NewBooleanNode[T Input](query string) (*BooleanNode[T], error) {
	q, err := gojq.Parse(query)
	if err != nil {
		return nil, err
	}
	check, err := gojq.Compile(q)
	if err != nil {
		return nil, err
	}
	return &BooleanNode[T]{
		check:    check,
		queryStr: query,
	}, nil
}

func (node *BooleanNode[T]) Merge(val Node[T]) {
	other, ok := val.(*BooleanNode[T])
	if !ok {
		panic("Incompatible merge")
	}
	if node.queryStr != other.queryStr {
		panic("Incompatible merge")
	}
	node.children[1].Merge(other.children[1])
}

func (node *BooleanNode[T]) Insert(a any, n Node[T]) {
	node.children[1] = n
}

func (node *BooleanNode[T]) Fingerprint() string {
	h := sha256.New()
	h.Write([]byte(node.queryStr))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (node *BooleanNode[T]) Query(i T) (Node[T], error) {
	v, ok := node.check.Run(i).Next()
	if !ok {
		return nil, ErrNotFound
	}

	switch v := v.(type) {
	case bool:
		if v {
			return node.children[1].Query(i)
		}
		return nil, ErrNotFound
	case error:
		return nil, v
	default:
		return nil, ErrNotFound
	}
}

func (node *BooleanNode[T]) String() string {
	return fmt.Sprintf("%v", node.children)
}
