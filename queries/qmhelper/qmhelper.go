package qmhelper

import (
	"fmt"
	"reflect"

	"github.com/nl2repo/sqlboiler/queries"
)

type Nullable interface {
	IsZero() bool
}

type Operator string

const (
	EQ  Operator = "="
	NEQ Operator = "!="
	LT  Operator = "<"
	LTE Operator = "<="
	GT  Operator = ">"
	GTE Operator = ">="
)

type WhereQueryMod struct {
	Clause string
	Args   []any
}

func (w WhereQueryMod) Apply(q *queries.Query) {
	queries.AppendWhere(q, w.Clause, w.Args...)
}

func Where(name string, op Operator, value any) WhereQueryMod {
	return WhereQueryMod{
		Clause: fmt.Sprintf("%s %s ?", name, string(op)),
		Args:   []any{value},
	}
}

func WhereNullEQ(name string, negated bool, value any) WhereQueryMod {
	isNull := isZeroValue(value)
	if isNull {
		if negated {
			return WhereQueryMod{Clause: fmt.Sprintf("%s is not null", name)}
		}
		return WhereQueryMod{Clause: fmt.Sprintf("%s is null", name)}
	}
	if negated {
		return WhereQueryMod{
			Clause: fmt.Sprintf("%s != ?", name),
			Args:   []any{value},
		}
	}
	return WhereQueryMod{
		Clause: fmt.Sprintf("%s = ?", name),
		Args:   []any{value},
	}
}

func WhereIsNull(name string) WhereQueryMod {
	return WhereQueryMod{Clause: fmt.Sprintf("%s is null", name)}
}

func WhereIsNotNull(name string) WhereQueryMod {
	return WhereQueryMod{Clause: fmt.Sprintf("%s is not null", name)}
}

func isZeroValue(val any) bool {
	if val == nil {
		return true
	}
	if n, ok := val.(Nullable); ok {
		return n.IsZero()
	}
	rv := reflect.ValueOf(val)
	return rv.IsZero()
}

func NonZeroDefaultSet(defaults []string, obj any) []string {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	result := make([]string, 0, len(defaults))
	for _, d := range defaults {
		fieldIdx := findFieldByColumn(typ, d)
		if fieldIdx < 0 {
			panic(fmt.Sprintf("qmhelper: could not find field for column %q in type %s", d, typ.Name()))
		}
		field := val.Field(fieldIdx)
		if !field.IsZero() {
			result = append(result, d)
		}
	}
	return result
}

func findFieldByColumn(typ reflect.Type, col string) int {
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		tag := f.Tag.Get("orm")
		if tag == "" {
			tag = f.Tag.Get("boil")
		}
		if tag == col {
			return i
		}
		if tag == "" {
			snake := queries.TitleToSnake(f.Name)
			if snake == col {
				return i
			}
		}
	}
	return -1
}
