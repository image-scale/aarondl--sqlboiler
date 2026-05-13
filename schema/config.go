package schema

import (
	"fmt"
	"strconv"
)

type Config map[string]any

func (c Config) MustString(key string) string {
	v, ok := c[key]
	if !ok {
		panic(fmt.Sprintf("schema.Config: missing key %q", key))
	}
	s, ok := v.(string)
	if !ok || s == "" {
		panic(fmt.Sprintf("schema.Config: key %q is not a non-empty string", key))
	}
	return s
}

func (c Config) String(key string) (string, bool) {
	v, ok := c[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok && s != ""
}

func (c Config) DefaultString(key, def string) string {
	s, ok := c.String(key)
	if !ok {
		return def
	}
	return s
}

func (c Config) MustInt(key string) int {
	v, ok := c[key]
	if !ok {
		panic(fmt.Sprintf("schema.Config: missing key %q", key))
	}
	return toInt(v)
}

func (c Config) Int(key string) (int, bool) {
	v, ok := c[key]
	if !ok {
		return 0, false
	}
	n := toInt(v)
	return n, n != 0
}

func (c Config) DefaultInt(key string, def int) int {
	n, ok := c.Int(key)
	if !ok {
		return def
	}
	return n
}

func (c Config) DefaultBool(key string, def bool) bool {
	v, ok := c[key]
	if !ok {
		return def
	}
	b, ok := v.(bool)
	if !ok {
		return def
	}
	return b
}

func (c Config) StringSlice(key string) ([]string, bool) {
	v, ok := c[key]
	if !ok {
		return nil, false
	}
	switch sv := v.(type) {
	case []string:
		return sv, true
	case []any:
		result := make([]string, 0, len(sv))
		for _, item := range sv {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result, true
	}
	return nil, false
}

func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case string:
		i, _ := strconv.Atoi(n)
		return i
	}
	return 0
}
