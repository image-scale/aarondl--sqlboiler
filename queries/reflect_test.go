package queries

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

type testUser struct {
	ID   int    `orm:"id"`
	Name string `orm:"name"`
	Age  int    `orm:"age"`
}

type testIgnored struct {
	ID      int    `orm:"id"`
	Secret  string `orm:"-"`
	Visible string `orm:"visible"`
}

type testNested struct {
	ID   int `orm:"id"`
	Info testNestedInfo `orm:",bind"`
}

type testNestedInfo struct {
	Email string `orm:"email"`
	Phone string `orm:"phone"`
}

type testAutoName struct {
	UserID      int
	FirstName   string
	HTTPSEnabled bool
}

func TestMakeStructMapping(t *testing.T) {
	m := MakeStructMapping(reflect.TypeOf(testUser{}))
	if _, ok := m["id"]; !ok {
		t.Error("mapping should contain 'id'")
	}
	if _, ok := m["name"]; !ok {
		t.Error("mapping should contain 'name'")
	}
	if _, ok := m["age"]; !ok {
		t.Error("mapping should contain 'age'")
	}
}

func TestMakeStructMappingIgnored(t *testing.T) {
	m := MakeStructMapping(reflect.TypeOf(testIgnored{}))
	if _, ok := m["secret"]; ok {
		t.Error("mapping should NOT contain '-' tagged field")
	}
	if _, ok := m["id"]; !ok {
		t.Error("mapping should contain 'id'")
	}
	if _, ok := m["visible"]; !ok {
		t.Error("mapping should contain 'visible'")
	}
}

func TestMakeStructMappingNested(t *testing.T) {
	m := MakeStructMapping(reflect.TypeOf(testNested{}))
	if _, ok := m["id"]; !ok {
		t.Error("mapping should contain 'id'")
	}
	if _, ok := m["email"]; !ok {
		t.Error("mapping should contain 'email' from nested struct")
	}
	if _, ok := m["phone"]; !ok {
		t.Error("mapping should contain 'phone' from nested struct")
	}
}

func TestTitleToSnake(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"UserID", "user_id"},
		{"FirstName", "first_name"},
		{"HTTPSEnabled", "https_enabled"},
		{"ID", "id"},
		{"Simple", "simple"},
		{"APIKey", "api_key"},
		{"HTMLParser", "html_parser"},
	}
	for _, tt := range tests {
		got := TitleToSnake(tt.input)
		if got != tt.want {
			t.Errorf("TitleToSnake(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMakeStructMappingAutoName(t *testing.T) {
	m := MakeStructMapping(reflect.TypeOf(testAutoName{}))
	if _, ok := m["user_id"]; !ok {
		t.Errorf("mapping should contain 'user_id', got keys: %v", mapKeys(m))
	}
	if _, ok := m["first_name"]; !ok {
		t.Errorf("mapping should contain 'first_name', got keys: %v", mapKeys(m))
	}
	if _, ok := m["https_enabled"]; !ok {
		t.Errorf("mapping should contain 'https_enabled', got keys: %v", mapKeys(m))
	}
}

func TestBindMapping(t *testing.T) {
	m := MakeStructMapping(reflect.TypeOf(testUser{}))
	cols := []string{"id", "name", "age"}
	result, err := BindMapping(reflect.TypeOf(testUser{}), m, cols)
	if err != nil {
		t.Fatalf("BindMapping error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 mappings, got %d", len(result))
	}
}

func TestBindMappingMissingColumn(t *testing.T) {
	m := MakeStructMapping(reflect.TypeOf(testUser{}))
	cols := []string{"id", "nonexistent"}
	_, err := BindMapping(reflect.TypeOf(testUser{}), m, cols)
	if err == nil {
		t.Error("expected error for missing column mapping")
	}
}

func TestPtrsFromMapping(t *testing.T) {
	m := MakeStructMapping(reflect.TypeOf(testUser{}))
	cols := []string{"id", "name", "age"}
	bindMap, _ := BindMapping(reflect.TypeOf(testUser{}), m, cols)

	u := testUser{}
	val := reflect.ValueOf(&u).Elem()
	ptrs := PtrsFromMapping(val, bindMap)
	if len(ptrs) != 3 {
		t.Fatalf("expected 3 ptrs, got %d", len(ptrs))
	}

	*(ptrs[0].(*int)) = 42
	*(ptrs[1].(*string)) = "alice"
	*(ptrs[2].(*int)) = 30

	if u.ID != 42 || u.Name != "alice" || u.Age != 30 {
		t.Errorf("PtrsFromMapping didn't point to correct fields: %+v", u)
	}
}

func TestValuesFromMapping(t *testing.T) {
	m := MakeStructMapping(reflect.TypeOf(testUser{}))
	cols := []string{"id", "name"}
	bindMap, _ := BindMapping(reflect.TypeOf(testUser{}), m, cols)

	u := testUser{ID: 1, Name: "bob"}
	val := reflect.ValueOf(&u).Elem()
	vals := ValuesFromMapping(val, bindMap)
	if len(vals) != 2 {
		t.Fatalf("expected 2 vals, got %d", len(vals))
	}
	if vals[0].(int) != 1 {
		t.Errorf("expected ID=1, got %v", vals[0])
	}
	if vals[1].(string) != "bob" {
		t.Errorf("expected Name=bob, got %v", vals[1])
	}
}

func TestBindSingleStruct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(1, "alice", 25)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		t.Fatal(err)
	}

	var u testUser
	err = Bind(sqlRows, &u)
	if err != nil {
		t.Fatalf("Bind error: %v", err)
	}
	if u.ID != 1 || u.Name != "alice" || u.Age != 25 {
		t.Errorf("unexpected bound result: %+v", u)
	}
}

func TestBindSingleStructNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "age"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		t.Fatal(err)
	}

	var u testUser
	err = Bind(sqlRows, &u)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestBindSliceStruct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(1, "alice", 25).
		AddRow(2, "bob", 30)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		t.Fatal(err)
	}

	var users []testUser
	err = Bind(sqlRows, &users)
	if err != nil {
		t.Fatalf("Bind error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Name != "alice" || users[1].Name != "bob" {
		t.Errorf("unexpected results: %+v", users)
	}
}

func TestBindPtrSliceStruct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(1, "alice", 25).
		AddRow(2, "bob", 30)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		t.Fatal(err)
	}

	var users []*testUser
	err = Bind(sqlRows, &users)
	if err != nil {
		t.Fatalf("Bind error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Name != "alice" || users[1].Name != "bob" {
		t.Errorf("unexpected results: %+v, %+v", users[0], users[1])
	}
}

func TestBindInvalidTarget(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	sqlRows, _ := db.Query("SELECT id")

	var x int
	err = Bind(sqlRows, &x)
	if err == nil {
		t.Error("expected error for non-struct target")
	}

	err = Bind(sqlRows, x)
	if err == nil {
		t.Error("expected error for non-pointer target")
	}
}

func TestEqual(t *testing.T) {
	if !Equal(nil, nil) {
		t.Error("nil == nil should be true")
	}
	if Equal(nil, 1) {
		t.Error("nil != 1")
	}
	if !Equal(1, 1) {
		t.Error("1 == 1")
	}
	if Equal(1, 2) {
		t.Error("1 != 2")
	}
	if !Equal([]byte{1, 2}, []byte{1, 2}) {
		t.Error("[]byte{1,2} == []byte{1,2}")
	}
	if Equal([]byte{1}, []byte{2}) {
		t.Error("[]byte{1} != []byte{2}")
	}

	now := time.Now()
	if !Equal(now, now) {
		t.Error("same time should be equal")
	}
	if Equal(now, now.Add(time.Second)) {
		t.Error("different times should not be equal")
	}
}

type testValuer struct {
	val any
}

func (tv testValuer) Value() (driver.Value, error) {
	return tv.val, nil
}

func TestEqualWithValuer(t *testing.T) {
	a := testValuer{val: "hello"}
	if !Equal(a, testValuer{val: "hello"}) {
		t.Error("valuers with same value should be equal")
	}
}

func TestIsNil(t *testing.T) {
	if !IsNil(nil) {
		t.Error("nil should be nil")
	}
	var p *int
	if !IsNil(p) {
		t.Error("nil pointer should be nil")
	}
	x := 5
	if IsNil(&x) {
		t.Error("non-nil pointer should not be nil")
	}
	if IsNil(42) {
		t.Error("non-pointer value should not be nil")
	}
}

type testScanner struct {
	Val any
}

func (ts *testScanner) Scan(src any) error {
	ts.Val = src
	return nil
}

func TestAssignWithScanner(t *testing.T) {
	s := &testScanner{}
	Assign(s, "hello")
	if s.Val != "hello" {
		t.Errorf("Assign with Scanner: expected 'hello', got %v", s.Val)
	}
}

func TestAssignWithReflect(t *testing.T) {
	var x int
	Assign(&x, 42)
	if x != 42 {
		t.Errorf("Assign with reflect: expected 42, got %d", x)
	}
}

func mapKeys(m map[string]uint64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
