package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSON json.RawMessage

func (j JSON) String() string {
	return string(j)
}

func (j JSON) Unmarshal(dest any) error {
	return json.Unmarshal([]byte(j), dest)
}

func (j *JSON) Marshal(src any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	*j = JSON(data)
	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return []byte(j), nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return fmt.Errorf("types.JSON: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(j), &raw); err != nil {
		return nil, err
	}
	return string(j), nil
}

func (j *JSON) Scan(src any) error {
	if src == nil {
		*j = nil
		return nil
	}
	switch v := src.(type) {
	case string:
		*j = JSON(v)
	case []byte:
		cp := make([]byte, len(v))
		copy(cp, v)
		*j = JSON(cp)
	default:
		return fmt.Errorf("types.JSON: cannot scan type %T", src)
	}
	return nil
}
