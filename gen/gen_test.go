package gen

import (
	"testing"

	"github.com/nl2repo/sqlboiler/schema"
)

func TestFillAliases(t *testing.T) {
	tables := []schema.Table{
		{
			Name:    "users",
			Columns: []schema.Column{{Name: "id"}, {Name: "first_name"}, {Name: "email"}},
		},
		{
			Name:    "blog_posts",
			Columns: []schema.Column{{Name: "id"}, {Name: "title"}, {Name: "user_id"}},
		},
	}
	a := Aliases{}
	FillAliases(&a, tables)

	if _, ok := a.Tables["users"]; !ok {
		t.Fatal("expected users alias")
	}
	if _, ok := a.Tables["blog_posts"]; !ok {
		t.Fatal("expected blog_posts alias")
	}

	ua := a.Tables["users"]
	if ua.UpSingular != "User" {
		t.Errorf("expected UpSingular=User, got %q", ua.UpSingular)
	}
	if ua.UpPlural != "Users" {
		t.Errorf("expected UpPlural=Users, got %q", ua.UpPlural)
	}
	if ua.DownSingular != "user" {
		t.Errorf("expected DownSingular=user, got %q", ua.DownSingular)
	}
	if ua.DownPlural != "users" {
		t.Errorf("expected DownPlural=users, got %q", ua.DownPlural)
	}
}

func TestFillAliasesColumns(t *testing.T) {
	tables := []schema.Table{
		{
			Name:    "users",
			Columns: []schema.Column{{Name: "id"}, {Name: "first_name"}, {Name: "email_address"}},
		},
	}
	a := Aliases{}
	FillAliases(&a, tables)

	ua := a.Tables["users"]
	if ua.Columns["id"] != "Id" {
		t.Errorf("expected Id, got %q", ua.Columns["id"])
	}
	if ua.Columns["first_name"] != "FirstName" {
		t.Errorf("expected FirstName, got %q", ua.Columns["first_name"])
	}
	if ua.Columns["email_address"] != "EmailAddress" {
		t.Errorf("expected EmailAddress, got %q", ua.Columns["email_address"])
	}
}

func TestFillAliasesSkipsJoinTables(t *testing.T) {
	tables := []schema.Table{
		{Name: "users"},
		{Name: "user_roles", IsJoinTable: true},
	}
	a := Aliases{}
	FillAliases(&a, tables)
	if _, ok := a.Tables["user_roles"]; ok {
		t.Error("join tables should be skipped")
	}
}

func TestAliasTablePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	a := Aliases{Tables: map[string]TableAlias{}}
	a.Table("nonexistent")
}

func TestAliasColumnPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	ta := TableAlias{Columns: map[string]string{}}
	ta.Column("nonexistent")
}

func TestTitleCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user_name", "UserName"},
		{"id", "Id"},
		{"blog_posts", "BlogPosts"},
		{"first-name", "FirstName"},
	}
	for _, tt := range tests {
		got := titleCase(tt.input)
		if got != tt.want {
			t.Errorf("titleCase(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user_name", "userName"},
		{"id", "id"},
		{"blog_posts", "blogPosts"},
	}
	for _, tt := range tests {
		got := camelCase(tt.input)
		if got != tt.want {
			t.Errorf("camelCase(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSingularize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"users", "user"},
		{"categories", "category"},
		{"boxes", "box"},
		{"addresses", "address"},
		{"bus", "bus"},
	}
	for _, tt := range tests {
		got := singularize(tt.input)
		if got != tt.want {
			t.Errorf("singularize(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user", "users"},
		{"category", "categories"},
		{"box", "boxes"},
		{"bus", "buses"},
	}
	for _, tt := range tests {
		got := pluralize(tt.input)
		if got != tt.want {
			t.Errorf("pluralize(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNameToOne(t *testing.T) {
	local, foreign := NameToOne("user_id", "users", "id", false)
	if local != "Users" {
		t.Errorf("expected local=Users, got %q", local)
	}
	if foreign != "User" {
		t.Errorf("expected foreign=User, got %q", foreign)
	}
}

func TestNameToOneUnique(t *testing.T) {
	local, foreign := NameToOne("user_id", "users", "id", true)
	if local != "User" {
		t.Errorf("expected local=User (singular because unique), got %q", local)
	}
	if foreign != "User" {
		t.Errorf("expected foreign=User, got %q", foreign)
	}
}

func TestNameToOneWithPrefix(t *testing.T) {
	local, foreign := NameToOne("producer_id", "users", "id", false)
	if local != "Users" {
		t.Errorf("expected local=Users, got %q", local)
	}
	if foreign != "ProducerUser" {
		t.Errorf("expected foreign=ProducerUser, got %q", foreign)
	}
}

func TestNameToMany(t *testing.T) {
	lhs, rhs := NameToMany("user_id", "tags", "tag_id", "users")
	if lhs != "Tags" {
		t.Errorf("expected lhs=Tags, got %q", lhs)
	}
	if rhs != "Users" {
		t.Errorf("expected rhs=Users, got %q", rhs)
	}
}

func TestIsPrimitive(t *testing.T) {
	if !IsPrimitive("int") {
		t.Error("int should be primitive")
	}
	if !IsPrimitive("string") {
		t.Error("string should be primitive")
	}
	if IsPrimitive("null.String") {
		t.Error("null.String should not be primitive")
	}
}

func TestIsNullPrimitive(t *testing.T) {
	if !IsNullPrimitive("null.String") {
		t.Error("null.String should be null primitive")
	}
	if !IsNullPrimitive("null.Int64") {
		t.Error("null.Int64 should be null primitive")
	}
	if IsNullPrimitive("string") {
		t.Error("string should not be null primitive")
	}
}

func TestConvertNullToPrimitive(t *testing.T) {
	if ConvertNullToPrimitive("null.Int64") != "int64" {
		t.Error("expected int64")
	}
	if ConvertNullToPrimitive("null.String") != "string" {
		t.Error("expected string")
	}
	if ConvertNullToPrimitive("int") != "int" {
		t.Error("non-null should pass through")
	}
}

func TestUsesPrimitives(t *testing.T) {
	tables := []schema.Table{
		{Name: "users", Columns: []schema.Column{{Name: "id", Type: "int"}}},
		{Name: "orders", Columns: []schema.Column{{Name: "id", Type: "int"}, {Name: "user_id", Type: "int"}}},
	}
	if !UsesPrimitives(tables, "orders", "user_id", "users", "id") {
		t.Error("both sides are primitives")
	}
}

func TestTrimSuffixes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user_id", "user"},
		{"session_uuid", "session"},
		{"owner_guid", "owner"},
		{"name", "name"},
	}
	for _, tt := range tests {
		got := trimSuffixes(tt.input)
		if got != tt.want {
			t.Errorf("trimSuffixes(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestConfigTagCaseConstants(t *testing.T) {
	if TagCaseCamel != "camel" {
		t.Error("TagCaseCamel should be 'camel'")
	}
	if TagCaseSnake != "snake" {
		t.Error("TagCaseSnake should be 'snake'")
	}
	if TagCaseTitle != "title" {
		t.Error("TagCaseTitle should be 'title'")
	}
	if TagCaseAlias != "alias" {
		t.Error("TagCaseAlias should be 'alias'")
	}
}
