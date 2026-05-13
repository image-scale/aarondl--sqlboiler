package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type HStore map[string]*string

func (h HStore) Value() (driver.Value, error) {
	if h == nil {
		return nil, nil
	}
	var parts []string
	for k, v := range h {
		if v == nil {
			parts = append(parts, fmt.Sprintf(`"%s"=>NULL`, hstoreEscape(k)))
		} else {
			parts = append(parts, fmt.Sprintf(`"%s"=>"%s"`, hstoreEscape(k), hstoreEscape(*v)))
		}
	}
	return strings.Join(parts, ", "), nil
}

func (h *HStore) Scan(src any) error {
	if src == nil {
		*h = nil
		return nil
	}
	s, err := asString(src)
	if err != nil {
		return fmt.Errorf("types.HStore: %w", err)
	}
	*h = parseHStore(s)
	return nil
}

func hstoreEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

func parseHStore(s string) HStore {
	result := make(HStore)
	s = strings.TrimSpace(s)
	if s == "" {
		return result
	}

	i := 0
	for i < len(s) {
		for i < len(s) && (s[i] == ' ' || s[i] == ',') {
			i++
		}
		if i >= len(s) {
			break
		}

		key := parseHStoreValue(&i, s)

		for i < len(s) && s[i] != '>' {
			i++
		}
		if i < len(s) {
			i++
		}

		for i < len(s) && s[i] == ' ' {
			i++
		}
		if i >= len(s) {
			break
		}

		if i+3 < len(s) && strings.ToUpper(s[i:i+4]) == "NULL" {
			result[key] = nil
			i += 4
		} else {
			val := parseHStoreValue(&i, s)
			result[key] = &val
		}
	}
	return result
}

func parseHStoreValue(i *int, s string) string {
	if *i >= len(s) {
		return ""
	}
	if s[*i] == '"' {
		*i++
		var buf strings.Builder
		for *i < len(s) && s[*i] != '"' {
			if s[*i] == '\\' && *i+1 < len(s) {
				*i++
				buf.WriteByte(s[*i])
			} else {
				buf.WriteByte(s[*i])
			}
			*i++
		}
		if *i < len(s) {
			*i++
		}
		return buf.String()
	}

	start := *i
	for *i < len(s) && s[*i] != ',' && s[*i] != '=' && s[*i] != ' ' {
		*i++
	}
	return s[start:*i]
}
