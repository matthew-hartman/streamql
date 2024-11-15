package build

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	var input map[string]any

	err := json.NewDecoder(strings.NewReader(testData)).Decode(&input)
	if err != nil {
		t.Error(err)
	}

	b, err := NewBuilder[map[string]any]().
		AddVarExp(`.after.Value.owner.string`, "mhartman").
		AddBoolExp(`.after.Value.parent.string=="eos-trunk"`).
		Build()
	if err != nil {
		t.Error(err)
	}

	t.Log(b.Query(input))

	b, err = b.
		AddVarExp(`.after.Value.owner.string`, "raviv").
		AddBoolExp(`.after.Value.parent.string=="eos-trunk"`).
		Build()
	if err != nil {
		t.Error(err)
	}
	t.Log(b.Query(input))
}

const testData = `{"after": {"Value": {"parent": {"string": "eos-trunk"}, "owner": {"string": "raviv"}}}}`
