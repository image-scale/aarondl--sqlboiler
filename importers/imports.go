package importers

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type List []string

func (l List) Len() int      { return len(l) }
func (l List) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l List) Less(i, j int) bool {
	return normalizeImport(l[i]) < normalizeImport(l[j])
}

func normalizeImport(s string) string {
	s = strings.TrimLeft(s, "_ ")
	return s
}

type Set struct {
	Standard   List
	ThirdParty List
}

func (s Set) Format() []byte {
	if len(s.Standard) == 0 && len(s.ThirdParty) == 0 {
		return nil
	}

	sort.Sort(s.Standard)
	sort.Sort(s.ThirdParty)

	var buf bytes.Buffer
	buf.WriteString("import (\n")
	for _, imp := range s.Standard {
		fmt.Fprintf(&buf, "\t%s\n", imp)
	}
	if len(s.Standard) > 0 && len(s.ThirdParty) > 0 {
		buf.WriteString("\n")
	}
	for _, imp := range s.ThirdParty {
		fmt.Fprintf(&buf, "\t%s\n", imp)
	}
	buf.WriteString(")\n")
	return buf.Bytes()
}

type Map map[string]Set

type Collection struct {
	All           Set
	Test          Set
	Singleton     Map
	TestSingleton Map
	BasedOnType   Map
}

func NewSet(standard, thirdParty []string) Set {
	return Set{
		Standard:   List(standard),
		ThirdParty: List(thirdParty),
	}
}

func Merge(a, b Collection) Collection {
	return Collection{
		All:           mergeSet(a.All, b.All),
		Test:          mergeSet(a.Test, b.Test),
		Singleton:     mergeMap(a.Singleton, b.Singleton),
		TestSingleton: mergeMap(a.TestSingleton, b.TestSingleton),
		BasedOnType:   mergeMap(a.BasedOnType, b.BasedOnType),
	}
}

func mergeSet(a, b Set) Set {
	return Set{
		Standard:   deduplicateList(append(a.Standard, b.Standard...)),
		ThirdParty: deduplicateList(append(a.ThirdParty, b.ThirdParty...)),
	}
}

func mergeMap(a, b Map) Map {
	if a == nil && b == nil {
		return nil
	}
	result := make(Map)
	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		if existing, ok := result[k]; ok {
			result[k] = mergeSet(existing, v)
		} else {
			result[k] = v
		}
	}
	return result
}

func deduplicateList(list List) List {
	seen := make(map[string]struct{}, len(list))
	result := make(List, 0, len(list))
	for _, item := range list {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func AddTypeImports(base Set, typeMap Map, columnTypes []string) Set {
	result := Set{
		Standard:   make(List, len(base.Standard)),
		ThirdParty: make(List, len(base.ThirdParty)),
	}
	copy(result.Standard, base.Standard)
	copy(result.ThirdParty, base.ThirdParty)

	seen := make(map[string]struct{})
	for _, ct := range columnTypes {
		if _, exists := seen[ct]; exists {
			continue
		}
		seen[ct] = struct{}{}
		if typeSet, ok := typeMap[ct]; ok {
			result = mergeSet(result, typeSet)
		}
	}
	return result
}

func NewDefaultImports() Collection {
	return Collection{
		All: Set{
			Standard: List{
				`"database/sql"`,
				`"fmt"`,
				`"reflect"`,
				`"strings"`,
				`"sync"`,
				`"time"`,
			},
			ThirdParty: List{
				`"github.com/friendsofgo/errors"`,
				`"github.com/nl2repo/sqlboiler/orm"`,
				`"github.com/nl2repo/sqlboiler/queries"`,
				`"github.com/nl2repo/sqlboiler/queries/qm"`,
			},
		},
		Test: Set{
			Standard: List{
				`"bytes"`,
				`"reflect"`,
				`"testing"`,
			},
			ThirdParty: List{
				`"github.com/nl2repo/sqlboiler/orm"`,
				`"github.com/nl2repo/sqlboiler/queries"`,
			},
		},
	}
}
