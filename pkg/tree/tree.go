package tree

import (
	"errors"
	"fmt"
)

type Input any

var (
	ErrNotFound = errors.New("no match")
)

type Node[T Input] interface {
	fmt.Stringer
	Fingerprint() string
	Query(T) (Node[T], error)
	Insert(any, Node[T])
	Merge(Node[T])
}
