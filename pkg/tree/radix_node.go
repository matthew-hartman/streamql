package tree

import (
	"crypto/sha256"
	"fmt"

	"github.com/itchyny/gojq"
)

const radixSize = 8

type RadixNode[T Input] struct {
	query    *gojq.Code
	queryStr string
	children [radixSize]*RadixNode[T]
	next     Node[T]
}

func NewRadixNode[T Input](query string) (*RadixNode[T], error) {
	q, err := gojq.Parse(query)
	if err != nil {
		return nil, err
	}
	check, err := gojq.Compile(q)
	if err != nil {
		return nil, err
	}
	return &RadixNode[T]{
		queryStr: query,
		query:    check,
	}, nil
}

func (node *RadixNode[T]) Merge(val Node[T]) {
	other, ok := val.(*RadixNode[T])
	if !ok {
		panic("Incompatible merge")
	}
	if node.queryStr != other.queryStr {
		panic("Incompatible merge")
	}
	for i, b := range other.children {
		if b == nil {
			continue
		}
		// If b[j] has a value and a[i] is nil use b[i]
		if b != nil && node.children[i] == nil {
			node.children[i] = b
			continue
		}
		// if both have values merge their children
		node.children[i].Merge(b)
	}
}

func (node *RadixNode[T]) Query(i T) (Node[T], error) {
	v, ok := node.query.Run(i).Next()
	if !ok {
		return nil, ErrNotFound
	}
	q := convert(v, len(node.children))
	j := node
	for _, a := range q {
		j = j.children[a]
		if j == nil {
			return nil, ErrNotFound
		}
	}
	if j.next != nil {
		return j.next.Query(i)
	}
	return nil, ErrNotFound
}

func (node *RadixNode[T]) Insert(input any, next Node[T]) {
	node.insert(convert(input, len(node.children)), next)
}

func (node *RadixNode[T]) Fingerprint() string {
	h := sha256.New()
	h.Write([]byte(node.queryStr))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (node *RadixNode[T]) String() string {
	return fmt.Sprintf("%v :: %v", node.next, node.children)
}

func (node *RadixNode[T]) insert(input []byte, next Node[T]) {
	j := input[0]
	if len(input) > 0 {
		input = input[1:]
	}
	if node.children[j] == nil {
		node.children[j] = &RadixNode[T]{}
	}
	if len(input) == 0 {
		node.children[j].next = next
		return
	}
	node.children[j].insert(input, next)
}

func convert(input any, base int) []byte {
	switch v1 := input.(type) {
	case float64:
		return convertInt(int(v1), base)
	case int:
		return convertInt(v1, base)
	case string:
		return convertString(v1, base)
	default:
		panic(fmt.Sprintf("unsupported type: %T", v1))
	}
}

func convertInt(input int, base int) []byte {
	var ret []byte
	for {
		a := input % base
		input = input / base
		ret = append(ret, byte(a))
		if input == 0 {
			break
		}
	}
	return ret
}

func convertString(input string, base int) []byte {
	var ret []byte
	for _, c := range input {
		i := c
		for {
			a := i % rune(base)
			i = i / rune(base)
			ret = append(ret, byte(a))
			if i == 0 {
				break
			}
		}
	}
	return ret
}
