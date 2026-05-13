package importers

import (
	"sort"
	"strings"
	"testing"
)

func TestListSorting(t *testing.T) {
	l := List{`"fmt"`, `"context"`, `_ "embed"`, `"bytes"`}
	result := make(List, len(l))
	copy(result, l)
	// Sort and verify ordering
	sorted := make(List, len(l))
	copy(sorted, l)
	sort.Sort(sorted)
	if sorted[0] != `"bytes"` {
		t.Errorf("first should be bytes after sort, got %q", sorted[0])
	}
}

func TestSetFormat(t *testing.T) {
	s := Set{
		Standard:   List{`"fmt"`, `"context"`},
		ThirdParty: List{`"github.com/pkg/errors"`},
	}
	out := s.Format()
	str := string(out)
	if !strings.Contains(str, "import (") {
		t.Error("expected import block")
	}
	if !strings.Contains(str, `"fmt"`) || !strings.Contains(str, `"context"`) {
		t.Error("expected standard imports")
	}
	if !strings.Contains(str, `"github.com/pkg/errors"`) {
		t.Error("expected third party import")
	}
}

func TestSetFormatEmpty(t *testing.T) {
	s := Set{}
	out := s.Format()
	if len(out) != 0 {
		t.Error("empty set should produce no output")
	}
}

func TestMergeCollections(t *testing.T) {
	a := Collection{
		All: Set{Standard: List{`"fmt"`}},
	}
	b := Collection{
		All: Set{Standard: List{`"context"`}, ThirdParty: List{`"github.com/pkg/errors"`}},
	}
	merged := Merge(a, b)
	if len(merged.All.Standard) != 2 {
		t.Errorf("expected 2 standard imports, got %d", len(merged.All.Standard))
	}
	if len(merged.All.ThirdParty) != 1 {
		t.Errorf("expected 1 third party import, got %d", len(merged.All.ThirdParty))
	}
}

func TestMergeDeduplicates(t *testing.T) {
	a := Collection{
		All: Set{Standard: List{`"fmt"`, `"context"`}},
	}
	b := Collection{
		All: Set{Standard: List{`"fmt"`, `"os"`}},
	}
	merged := Merge(a, b)
	if len(merged.All.Standard) != 3 {
		t.Errorf("expected 3 unique standard imports after dedup, got %d: %v", len(merged.All.Standard), merged.All.Standard)
	}
}

func TestAddTypeImports(t *testing.T) {
	base := Set{
		Standard: List{`"fmt"`},
	}
	typeMap := Map{
		"json.RawMessage": {ThirdParty: List{`"encoding/json"`}},
		"null.String":     {ThirdParty: List{`"github.com/aarondl/null/v8"`}},
	}
	columnTypes := []string{"json.RawMessage", "null.String", "string"}
	result := AddTypeImports(base, typeMap, columnTypes)
	if len(result.Standard) != 1 {
		t.Errorf("expected 1 standard import, got %d", len(result.Standard))
	}
	if len(result.ThirdParty) != 2 {
		t.Errorf("expected 2 third party imports, got %d: %v", len(result.ThirdParty), result.ThirdParty)
	}
}

func TestAddTypeImportsDeduplicatesTypes(t *testing.T) {
	base := Set{}
	typeMap := Map{
		"null.String": {ThirdParty: List{`"github.com/aarondl/null/v8"`}},
	}
	columnTypes := []string{"null.String", "null.String", "null.String"}
	result := AddTypeImports(base, typeMap, columnTypes)
	if len(result.ThirdParty) != 1 {
		t.Errorf("expected 1 import after dedup, got %d", len(result.ThirdParty))
	}
}

func TestNewSet(t *testing.T) {
	s := NewSet([]string{`"fmt"`}, []string{`"github.com/x/y"`})
	if len(s.Standard) != 1 || len(s.ThirdParty) != 1 {
		t.Errorf("unexpected set: %+v", s)
	}
}

func TestMergeMap(t *testing.T) {
	a := Map{
		"file1": {Standard: List{`"fmt"`}},
	}
	b := Map{
		"file1": {Standard: List{`"os"`}},
		"file2": {ThirdParty: List{`"github.com/x"`}},
	}
	result := mergeMap(a, b)
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
	if len(result["file1"].Standard) != 2 {
		t.Errorf("expected 2 standard imports for file1, got %d", len(result["file1"].Standard))
	}
}

func TestNewDefaultImports(t *testing.T) {
	c := NewDefaultImports()
	if len(c.All.Standard) == 0 {
		t.Error("default imports should have standard imports")
	}
	if len(c.All.ThirdParty) == 0 {
		t.Error("default imports should have third party imports")
	}
	if len(c.Test.Standard) == 0 {
		t.Error("default test imports should have standard imports")
	}
}

func TestCollectionNilMerge(t *testing.T) {
	a := Collection{}
	b := Collection{
		Singleton: Map{"x": {Standard: List{`"os"`}}},
	}
	result := Merge(a, b)
	if result.Singleton == nil || len(result.Singleton["x"].Standard) != 1 {
		t.Error("merge with nil maps should work")
	}
}
