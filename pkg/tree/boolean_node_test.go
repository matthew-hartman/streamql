package tree

import (
	"fmt"
	"testing"
)

type TestNode[T Input] struct {
	ID int64
}

func (_ *TestNode[T]) Query(i T) (Node[T], error) {
	return nil, nil
}

func (t *TestNode[T]) String() string {
	return fmt.Sprint(t.ID)
}

func TestBooleanNode(t *testing.T) {

	tt := []struct {
		query string
		input any
		exp   bool
	}{
		{".test0==true", map[string]any{"test0": true}, true},
		{".test1==true", map[string]any{"test": true}, false},
		{".test2==true", map[string]any{"test2": false}, false},
	}

	for _, tc := range tt {
		t.Run(tc.query, func(t *testing.T) {
			node, err := NewBooleanNode[any](tc.query)
			if err != nil {
				t.Error(err)
			}

			n, err := node.Query(tc.input)
			if err != nil {
				t.Error(err)
			}
			result := false
			if i, ok := n.(*ResultNode[any, bool]); ok {
				result = i.Result()
			}
			if result != tc.exp {
				t.Errorf("incorrect return %v, %v", result, n)
			}
		})
	}
}
