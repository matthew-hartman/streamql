package tree

import "fmt"

type ResultNode[T, R Input] struct{ result R }

func NewResult[T, R Input](res R) Node[T]               { return &ResultNode[T, R]{res} }
func (node *ResultNode[T, R]) Insert(any, Node[T])      {}
func (node *ResultNode[T, R]) Fingerprint() string      { return "" }
func (node *ResultNode[T, R]) Query(T) (Node[T], error) { return node, nil }
func (node *ResultNode[T, R]) Result() R                { return node.result }
func (node *ResultNode[T, R]) String() string           { return fmt.Sprintf("%v", node.result) }
func (node *ResultNode[T, R]) Merge(val Node[T]) {
	other, ok := val.(*ResultNode[T, R])
	if !ok {
		panic("Incompatible merge")
	}
	if node.String() != other.String() {
		panic("Incompatible merge")
	}
}

type CallbackNode[T Input] struct{ callback func(T) }

func NewCallback[T Input](callback func(T)) Node[T]          { return &CallbackNode[T]{callback} }
func (node *CallbackNode[T]) Insert(any, Node[T])            {}
func (node *CallbackNode[T]) Query(input T) (Node[T], error) { node.callback(input); return nil, nil }
func (node *CallbackNode[T]) Fingerprint() string            { return fmt.Sprintf("%p", node) }
func (node *CallbackNode[T]) String() string                 { return fmt.Sprintf("%#v", node) }

func (node *CallbackNode[T]) Merge(val Node[T]) {
	other, ok := val.(*CallbackNode[T])
	if !ok {
		panic("Incompatible merge")
	}
	if node.String() != other.String() {
		panic("Incompatible merge")
	}
}
