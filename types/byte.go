package types

import (
	"database/sql/driver"
	"fmt"
)

type Byte byte

func (b Byte) String() string {
	return string(rune(b))
}

func (b Byte) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", string(rune(b)))), nil
}

func (b *Byte) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("types.Byte: data too short")
	}
	s := string(data[1 : len(data)-1])
	if len(s) > 1 {
		return fmt.Errorf("types.Byte: too many characters")
	}
	if len(s) == 0 {
		*b = 0
		return nil
	}
	*b = Byte(s[0])
	return nil
}

func (b Byte) Value() (driver.Value, error) {
	return []byte{byte(b)}, nil
}

func (b *Byte) Scan(src any) error {
	if src == nil {
		*b = 0
		return nil
	}
	switch v := src.(type) {
	case uint8:
		*b = Byte(v)
	case string:
		if len(v) > 0 {
			*b = Byte(v[0])
		}
	case []byte:
		if len(v) > 0 {
			*b = Byte(v[0])
		}
	default:
		return fmt.Errorf("types.Byte: cannot scan type %T", src)
	}
	return nil
}
