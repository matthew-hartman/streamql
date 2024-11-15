package tree

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTree(t *testing.T) {
	var input map[string]any

	err := json.NewDecoder(strings.NewReader(testData)).Decode(&input)
	if err != nil {
		t.Error(err)
	}

	node, err := NewBooleanNode[any](`.after.Value.parent.string=="eos-trunk"`)
	if err != nil {
		t.Error(err)
	}

	a, err := node.Query(input)
	if err != nil {
		t.Error(err)
	}
	if a == nil || a.String() != "true" {
		t.Error("unexpected result", a)
	}

	r, err := NewRadixNode[any](`.after.Value.owner.string`)
	if err != nil {
		t.Error(err)
	}

	r.Insert("mhartman", NewResult[any, bool](true))
	r.Insert("raviv", NewResult[any, bool](true))

	if a == nil || a.String() != "true" {
		t.Error("unexpected result", a)
	}

	a, err = r.Query(input)
	if err != nil {
		t.Error(err)
	}

	if a == nil || a.String() != "true" {
		t.Error("unexpected result", a)
	}
}

const testData = `{"after": {"Value": {"parent": {"string": "eos-trunk"}, "owner": {"string": "raviv"}}}}`
