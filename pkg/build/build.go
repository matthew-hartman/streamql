package build

import (
	"fmt"
	"strings"

	"code.arista.io/lib/streamql/pkg/tree"
)

type Builder[T tree.Input] struct {
	root  map[string]tree.Node[T]
	build []buildStep[T]
}

type buildStep[T tree.Input] struct {
	node tree.Node[T]
	arg  any
}

func NewBuilder[T tree.Input]() *Builder[T] {
	return &Builder[T]{
		root: map[string]tree.Node[T]{},
	}
}

func (b *Builder[T]) AddBoolExp(query string) *Builder[T] {
	n, err := tree.NewBooleanNode[T](query)
	if err != nil {
		panic(err)
	}
	b.build = append(b.build, buildStep[T]{
		node: n,
	})
	return b
}

func (b *Builder[T]) AddVarExp(query string, arg any) *Builder[T] {
	n, err := tree.NewRadixNode[T](query)
	if err != nil {
		panic(err)
	}
	b.build = append(b.build, buildStep[T]{
		node: n,
		arg:  arg,
	})
	return b
}

type BuildOpts[T tree.Input] struct {
	result tree.Node[T]
}

type BuildOpt[T tree.Input] func(*BuildOpts[T])

func WithCallback[T tree.Input](fn func(T)) BuildOpt[T] {
	return func(opt *BuildOpts[T]) {
		opt.result = tree.NewCallback[T](fn)
	}
}

func (b *Builder[T]) Build(options ...BuildOpt[T]) (bi *Builder[T], err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	opts := BuildOpts[T]{
		result: tree.NewResult[T](true),
	}
	for _, o := range options {
		o(&opts)
	}
	bi = b
	fp := []string{}
	for _, v := range b.build {
		fp = append(fp, v.node.Fingerprint())
	}
	p := strings.Join(fp, "::")

	var root tree.Node[T] = opts.result
	for _, v := range b.build {
		v.node.Insert(v.arg, root)
		root = v.node
	}

	if _, ok := b.root[p]; !ok {
		b.root[p] = root
	} else {
		b.root[p].Merge(root)
	}
	b.build = []buildStep[T]{}

	return
}

func (b *Builder[T]) Query(input T) (tree.Node[T], error) {
	for _, v := range b.root {
		return v.Query(input)
	}
	return nil, fmt.Errorf("not found")
}
