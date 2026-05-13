package types

import (
	"encoding/json"
	"testing"
)

func TestJSONMarshalUnmarshal(t *testing.T) {
	var j JSON
	err := j.Marshal(map[string]int{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if j.String() == "" {
		t.Error("expected non-empty JSON string")
	}

	var m map[string]int
	err = j.Unmarshal(&m)
	if err != nil {
		t.Fatal(err)
	}
	if m["a"] != 1 {
		t.Errorf("expected a=1, got %v", m)
	}
}

func TestJSONScanValue(t *testing.T) {
	var j JSON
	err := j.Scan(`{"hello":"world"}`)
	if err != nil {
		t.Fatal(err)
	}
	val, err := j.Value()
	if err != nil {
		t.Fatal(err)
	}
	if val == nil {
		t.Error("expected non-nil value")
	}
}

func TestJSONScanNil(t *testing.T) {
	j := JSON(`{"test":1}`)
	err := j.Scan(nil)
	if err != nil {
		t.Fatal(err)
	}
	if j != nil {
		t.Errorf("expected nil after scanning nil, got %v", j)
	}
}

func TestJSONScanBytes(t *testing.T) {
	var j JSON
	err := j.Scan([]byte(`[1,2,3]`))
	if err != nil {
		t.Fatal(err)
	}
	if string(j) != "[1,2,3]" {
		t.Errorf("unexpected JSON: %s", j)
	}
}

func TestJSONValueEmpty(t *testing.T) {
	var j JSON
	val, err := j.Value()
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Error("empty JSON should return nil value")
	}
}

func TestJSONMarshalJSON(t *testing.T) {
	j := JSON(`{"key":"val"}`)
	data, err := json.Marshal(j)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"key":"val"}` {
		t.Errorf("unexpected marshaled: %s", data)
	}
}

func TestJSONUnmarshalJSON(t *testing.T) {
	var j JSON
	err := json.Unmarshal([]byte(`[1,2]`), &j)
	if err != nil {
		t.Fatal(err)
	}
	if string(j) != "[1,2]" {
		t.Errorf("unexpected unmarshaled: %s", j)
	}
}

func TestByteScanValue(t *testing.T) {
	var b Byte
	err := b.Scan(uint8('A'))
	if err != nil {
		t.Fatal(err)
	}
	if b != Byte('A') {
		t.Errorf("expected A, got %c", b)
	}
	val, err := b.Value()
	if err != nil {
		t.Fatal(err)
	}
	bs := val.([]byte)
	if bs[0] != 'A' {
		t.Errorf("expected [A], got %v", bs)
	}
}

func TestByteScanString(t *testing.T) {
	var b Byte
	err := b.Scan("X")
	if err != nil {
		t.Fatal(err)
	}
	if b != Byte('X') {
		t.Errorf("expected X, got %c", b)
	}
}

func TestByteScanBytes(t *testing.T) {
	var b Byte
	err := b.Scan([]byte("Z"))
	if err != nil {
		t.Fatal(err)
	}
	if b != Byte('Z') {
		t.Errorf("expected Z, got %c", b)
	}
}

func TestByteJSON(t *testing.T) {
	b := Byte('M')
	data, err := b.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var b2 Byte
	err = b2.UnmarshalJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if b2 != Byte('M') {
		t.Errorf("expected M, got %c", b2)
	}
}

func TestByteUnmarshalJSONTooLong(t *testing.T) {
	var b Byte
	err := b.UnmarshalJSON([]byte(`"AB"`))
	if err == nil {
		t.Error("expected error for too many characters")
	}
}

func TestInt64ArrayScanValue(t *testing.T) {
	var a Int64Array
	err := a.Scan("{1,2,3}")
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 3 || a[0] != 1 || a[1] != 2 || a[2] != 3 {
		t.Errorf("unexpected int64 array: %v", a)
	}
	val, err := a.Value()
	if err != nil {
		t.Fatal(err)
	}
	if val != "{1,2,3}" {
		t.Errorf("unexpected value: %v", val)
	}
}

func TestInt64ArrayNil(t *testing.T) {
	var a Int64Array
	err := a.Scan(nil)
	if err != nil {
		t.Fatal(err)
	}
	if a != nil {
		t.Error("expected nil after scanning nil")
	}
	val, err := a.Value()
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Error("nil array should produce nil value")
	}
}

func TestFloat64ArrayScanValue(t *testing.T) {
	var a Float64Array
	err := a.Scan("{1.5,2.7,3.14}")
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(a))
	}
	if a[0] != 1.5 || a[2] != 3.14 {
		t.Errorf("unexpected float64 array: %v", a)
	}
	val, err := a.Value()
	if err != nil {
		t.Fatal(err)
	}
	s := val.(string)
	if s == "" {
		t.Error("expected non-empty value")
	}
}

func TestBoolArrayScanValue(t *testing.T) {
	var a BoolArray
	err := a.Scan("{t,f,true,false}")
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 4 {
		t.Fatalf("expected 4 elements, got %d", len(a))
	}
	if !a[0] || a[1] || !a[2] || a[3] {
		t.Errorf("unexpected bool array: %v", a)
	}
	val, err := a.Value()
	if err != nil {
		t.Fatal(err)
	}
	if val != "{t,f,t,f}" {
		t.Errorf("unexpected value: %v", val)
	}
}

func TestStringArrayScanValue(t *testing.T) {
	var a StringArray
	err := a.Scan(`{"hello","world"}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 2 || a[0] != "hello" || a[1] != "world" {
		t.Errorf("unexpected string array: %v", a)
	}
	val, err := a.Value()
	if err != nil {
		t.Fatal(err)
	}
	if val == nil {
		t.Error("expected non-nil value")
	}
}

func TestStringArrayWithEscapes(t *testing.T) {
	var a StringArray
	err := a.Scan(`{"he\"llo","wor\\ld"}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(a))
	}
	if a[0] != `he"llo` {
		t.Errorf("expected he\"llo, got %q", a[0])
	}
	if a[1] != `wor\ld` {
		t.Errorf("expected wor\\ld, got %q", a[1])
	}
}

func TestBytesArrayScanValue(t *testing.T) {
	var a BytesArray
	err := a.Scan(`{"abc","def"}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(a))
	}
	if string(a[0]) != "abc" || string(a[1]) != "def" {
		t.Errorf("unexpected bytes array: %v", a)
	}
}

func TestHStoreScanValue(t *testing.T) {
	var h HStore
	err := h.Scan(`"key1"=>"val1", "key2"=>NULL`)
	if err != nil {
		t.Fatal(err)
	}
	if len(h) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(h))
	}
	if h["key1"] == nil || *h["key1"] != "val1" {
		t.Errorf("expected key1=val1, got %v", h["key1"])
	}
	if h["key2"] != nil {
		t.Errorf("expected key2=NULL, got %v", h["key2"])
	}
}

func TestHStoreValue(t *testing.T) {
	v1 := "world"
	h := HStore{"hello": &v1, "nil_key": nil}
	val, err := h.Value()
	if err != nil {
		t.Fatal(err)
	}
	if val == nil {
		t.Error("expected non-nil value")
	}
}

func TestHStoreNil(t *testing.T) {
	var h HStore
	err := h.Scan(nil)
	if err != nil {
		t.Fatal(err)
	}
	if h != nil {
		t.Error("expected nil after scanning nil")
	}
	val, err := h.Value()
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Error("nil hstore should produce nil value")
	}
}

func TestHStoreEscapedValues(t *testing.T) {
	var h HStore
	err := h.Scan(`"k\"ey"=>"v\\al"`)
	if err != nil {
		t.Fatal(err)
	}
	if len(h) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h))
	}
	if h[`k"ey`] == nil || *h[`k"ey`] != `v\al` {
		t.Errorf("unexpected hstore: %v", h)
	}
}

func TestStringArrayEmpty(t *testing.T) {
	var a StringArray
	err := a.Scan("{}")
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 0 {
		t.Errorf("expected empty array, got %v", a)
	}
}

func TestInt64ArrayScanBytes(t *testing.T) {
	var a Int64Array
	err := a.Scan([]byte("{10,20}"))
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 2 || a[0] != 10 || a[1] != 20 {
		t.Errorf("unexpected: %v", a)
	}
}
