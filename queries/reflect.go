package queries

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode"
)

const sentinel = uint64(255)

type bindKind int

const (
	kindStruct         bindKind = iota
	kindSliceStruct
	kindPtrSliceStruct
)

func Bind(rows *sql.Rows, obj any) error {
	structType, sliceType, bk, err := bindChecks(obj)
	if err != nil {
		return err
	}
	return bind(rows, obj, structType, sliceType, bk)
}

func bindChecks(obj any) (reflect.Type, reflect.Type, bindKind, error) {
	typ := reflect.TypeOf(obj)
	if typ.Kind() != reflect.Ptr {
		return nil, nil, 0, fmt.Errorf("bind: obj must be a pointer, got %s", typ.Kind())
	}

	elem := typ.Elem()
	switch elem.Kind() {
	case reflect.Struct:
		return elem, nil, kindStruct, nil
	case reflect.Slice:
		sliceElem := elem.Elem()
		if sliceElem.Kind() == reflect.Ptr {
			innerStruct := sliceElem.Elem()
			if innerStruct.Kind() != reflect.Struct {
				return nil, nil, 0, fmt.Errorf("bind: *[]*T requires T to be a struct, got %s", innerStruct.Kind())
			}
			return innerStruct, elem, kindPtrSliceStruct, nil
		}
		if sliceElem.Kind() == reflect.Struct {
			return sliceElem, elem, kindSliceStruct, nil
		}
		return nil, nil, 0, fmt.Errorf("bind: slice element must be struct or *struct, got %s", sliceElem.Kind())
	default:
		return nil, nil, 0, fmt.Errorf("bind: unsupported type %s", elem.Kind())
	}
}

func bind(rows *sql.Rows, obj any, structType, sliceType reflect.Type, bk bindKind) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	mapping := MakeStructMapping(structType)
	bindMap, err := BindMapping(structType, mapping, cols)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(obj).Elem()
	found := false

	for rows.Next() {
		found = true
		var rowVal reflect.Value

		switch bk {
		case kindStruct:
			rowVal = val
		case kindSliceStruct:
			rowVal = reflect.New(structType).Elem()
		case kindPtrSliceStruct:
			rowVal = reflect.New(structType).Elem()
		}

		ptrs := PtrsFromMapping(rowVal, bindMap)
		if err := rows.Scan(ptrs...); err != nil {
			return err
		}

		switch bk {
		case kindSliceStruct:
			val.Set(reflect.Append(val, rowVal))
		case kindPtrSliceStruct:
			ptrVal := reflect.New(structType)
			ptrVal.Elem().Set(rowVal)
			val.Set(reflect.Append(val, ptrVal))
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if !found && bk == kindStruct {
		return sql.ErrNoRows
	}
	return nil
}

var mappingCacheMu sync.RWMutex
var mappingCacheMap = make(map[reflect.Type]map[string]uint64)

func MakeStructMapping(typ reflect.Type) map[string]uint64 {
	mappingCacheMu.RLock()
	if m, ok := mappingCacheMap[typ]; ok {
		mappingCacheMu.RUnlock()
		return m
	}
	mappingCacheMu.RUnlock()

	m := make(map[string]uint64)
	buildMapping(typ, m, 0, 0)

	mappingCacheMu.Lock()
	mappingCacheMap[typ] = m
	mappingCacheMu.Unlock()
	return m
}

func buildMapping(typ reflect.Type, mapping map[string]uint64, depth int, prefix uint64) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}

		name, recurse := getOrmTag(field)
		if name == "-" {
			continue
		}

		encoded := prefix | (uint64(i) << (uint(depth) * 8))
		finalEncoded := encoded | (sentinel << (uint(depth+1) * 8))

		if recurse && field.Type.Kind() == reflect.Struct {
			buildMapping(field.Type, mapping, depth+1, encoded)
			continue
		}

		if name == "" {
			name = TitleToSnake(field.Name)
		}
		mapping[name] = finalEncoded
	}
}

func getOrmTag(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get("orm")
	if tag == "" {
		tag = field.Tag.Get("boil")
	}
	if tag == "" {
		return "", false
	}
	if tag == "-" {
		return "-", false
	}

	parts := strings.Split(tag, ",")
	name := parts[0]
	recurse := false
	for _, p := range parts[1:] {
		if p == "bind" {
			recurse = true
		}
	}
	return name, recurse
}

func BindMapping(typ reflect.Type, mapping map[string]uint64, cols []string) ([]uint64, error) {
	result := make([]uint64, len(cols))
	for i, col := range cols {
		enc, ok := mapping[col]
		if !ok {
			return nil, fmt.Errorf("bind: could not find mapping for column %q in type %s", col, typ.Name())
		}
		result[i] = enc
	}
	return result, nil
}

func PtrsFromMapping(val reflect.Value, mapping []uint64) []any {
	ptrs := make([]any, len(mapping))
	for i, enc := range mapping {
		ptrs[i] = ptrFromMapping(val, enc, true).Interface()
	}
	return ptrs
}

func ValuesFromMapping(val reflect.Value, mapping []uint64) []any {
	vals := make([]any, len(mapping))
	for i, enc := range mapping {
		vals[i] = ptrFromMapping(val, enc, false).Interface()
	}
	return vals
}

func ptrFromMapping(val reflect.Value, enc uint64, addressOf bool) reflect.Value {
	current := val
	for depth := 0; depth < 8; depth++ {
		idx := int((enc >> (uint(depth) * 8)) & 0xFF)
		nextByte := (enc >> (uint(depth+1) * 8)) & 0xFF
		if nextByte == sentinel {
			field := current.Field(idx)
			if addressOf {
				return field.Addr()
			}
			return field
		}
		current = current.Field(idx)
	}
	return val
}

func TitleToSnake(name string) string {
	var buf strings.Builder
	runes := []rune(name)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]
				if unicode.IsLower(prev) || unicode.IsDigit(prev) {
					buf.WriteByte('_')
				} else if unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					buf.WriteByte('_')
				}
			}
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func Equal(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	aVal := extractValue(a)
	bVal := extractValue(b)

	if aVal == nil && bVal == nil {
		return true
	}
	if aVal == nil || bVal == nil {
		return false
	}

	switch av := aVal.(type) {
	case []byte:
		bv, ok := bVal.([]byte)
		if !ok {
			return false
		}
		if len(av) != len(bv) {
			return false
		}
		for i := range av {
			if av[i] != bv[i] {
				return false
			}
		}
		return true
	case time.Time:
		bv, ok := bVal.(time.Time)
		if !ok {
			return false
		}
		return av.Equal(bv)
	}

	return reflect.DeepEqual(aVal, bVal)
}

func extractValue(v any) any {
	if valuer, ok := v.(driver.Valuer); ok {
		val, err := valuer.Value()
		if err != nil {
			return v
		}
		return val
	}
	return v
}

func Assign(dst, src any) {
	if scanner, ok := dst.(sql.Scanner); ok {
		srcVal := extractValue(src)
		scanner.Scan(srcVal)
		return
	}
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr {
		return
	}
	srcVal := reflect.ValueOf(src)
	if srcVal.IsValid() && srcVal.Type().AssignableTo(dstVal.Elem().Type()) {
		dstVal.Elem().Set(srcVal)
	}
}

func IsNil(val any) bool {
	if val == nil {
		return true
	}
	if valuer, ok := val.(driver.Valuer); ok {
		v, err := valuer.Value()
		if err != nil {
			return false
		}
		return v == nil
	}
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}
