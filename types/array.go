package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Int64Array []int64

func (a Int64Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	parts := make([]string, len(a))
	for i, v := range a {
		parts[i] = strconv.FormatInt(v, 10)
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (a *Int64Array) Scan(src any) error {
	if src == nil {
		*a = nil
		return nil
	}
	s, err := asString(src)
	if err != nil {
		return fmt.Errorf("types.Int64Array: %w", err)
	}
	elems := parseArrayElements(s)
	result := make(Int64Array, len(elems))
	for i, e := range elems {
		v, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			return fmt.Errorf("types.Int64Array: parsing %q: %w", e, err)
		}
		result[i] = v
	}
	*a = result
	return nil
}

type Float64Array []float64

func (a Float64Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	parts := make([]string, len(a))
	for i, v := range a {
		parts[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (a *Float64Array) Scan(src any) error {
	if src == nil {
		*a = nil
		return nil
	}
	s, err := asString(src)
	if err != nil {
		return fmt.Errorf("types.Float64Array: %w", err)
	}
	elems := parseArrayElements(s)
	result := make(Float64Array, len(elems))
	for i, e := range elems {
		v, err := strconv.ParseFloat(e, 64)
		if err != nil {
			return fmt.Errorf("types.Float64Array: parsing %q: %w", e, err)
		}
		result[i] = v
	}
	*a = result
	return nil
}

type BoolArray []bool

func (a BoolArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	parts := make([]string, len(a))
	for i, v := range a {
		if v {
			parts[i] = "t"
		} else {
			parts[i] = "f"
		}
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (a *BoolArray) Scan(src any) error {
	if src == nil {
		*a = nil
		return nil
	}
	s, err := asString(src)
	if err != nil {
		return fmt.Errorf("types.BoolArray: %w", err)
	}
	elems := parseArrayElements(s)
	result := make(BoolArray, len(elems))
	for i, e := range elems {
		switch strings.ToLower(e) {
		case "t", "true", "1":
			result[i] = true
		case "f", "false", "0":
			result[i] = false
		default:
			return fmt.Errorf("types.BoolArray: unknown bool %q", e)
		}
	}
	*a = result
	return nil
}

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	parts := make([]string, len(a))
	for i, v := range a {
		escaped := strings.ReplaceAll(v, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		parts[i] = `"` + escaped + `"`
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (a *StringArray) Scan(src any) error {
	if src == nil {
		*a = nil
		return nil
	}
	s, err := asString(src)
	if err != nil {
		return fmt.Errorf("types.StringArray: %w", err)
	}
	*a = parseStringArrayElements(s)
	return nil
}

type BytesArray [][]byte

func (a BytesArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	parts := make([]string, len(a))
	for i, v := range a {
		parts[i] = fmt.Sprintf(`"\\x%x"`, v)
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (a *BytesArray) Scan(src any) error {
	if src == nil {
		*a = nil
		return nil
	}
	s, err := asString(src)
	if err != nil {
		return fmt.Errorf("types.BytesArray: %w", err)
	}
	elems := parseStringArrayElements(s)
	result := make(BytesArray, len(elems))
	for i, e := range elems {
		result[i] = []byte(e)
	}
	*a = result
	return nil
}

func asString(src any) (string, error) {
	switch v := src.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return "", fmt.Errorf("cannot convert %T to string", src)
	}
}

func parseArrayElements(s string) []string {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return nil
	}
	s = s[1 : len(s)-1]
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func parseStringArrayElements(s string) []string {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return nil
	}
	s = s[1 : len(s)-1]
	if s == "" {
		return nil
	}

	var result []string
	i := 0
	for i < len(s) {
		if s[i] == '"' {
			i++
			var elem strings.Builder
			for i < len(s) && s[i] != '"' {
				if s[i] == '\\' && i+1 < len(s) {
					i++
					elem.WriteByte(s[i])
				} else {
					elem.WriteByte(s[i])
				}
				i++
			}
			result = append(result, elem.String())
			i++
			if i < len(s) && s[i] == ',' {
				i++
			}
		} else {
			end := strings.IndexByte(s[i:], ',')
			if end == -1 {
				result = append(result, s[i:])
				break
			}
			result = append(result, s[i:i+end])
			i += end + 1
		}
	}
	return result
}
