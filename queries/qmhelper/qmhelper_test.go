package qmhelper

import (
	"testing"

	"github.com/nl2repo/sqlboiler/queries"
)

type nullableString struct {
	Val   string
	Valid bool
}

func (n nullableString) IsZero() bool {
	return !n.Valid
}

func TestWhereBasicOperators(t *testing.T) {
	tests := []struct {
		op   Operator
		want string
	}{
		{EQ, "name = ?"},
		{NEQ, "name != ?"},
		{LT, "age < ?"},
		{LTE, "age <= ?"},
		{GT, "age > ?"},
		{GTE, "age >= ?"},
	}

	for _, tt := range tests {
		col := "name"
		if tt.op == LT || tt.op == LTE || tt.op == GT || tt.op == GTE {
			col = "age"
		}
		mod := Where(col, tt.op, "value")
		if mod.Clause != tt.want {
			t.Errorf("Where(%q, %q) clause = %q, want %q", col, tt.op, mod.Clause, tt.want)
		}
		if len(mod.Args) != 1 {
			t.Errorf("expected 1 arg, got %d", len(mod.Args))
		}
	}
}

func TestWhereNullEQWithNull(t *testing.T) {
	mod := WhereNullEQ("email", false, nullableString{Val: "", Valid: false})
	if mod.Clause != "email is null" {
		t.Errorf("expected 'email is null', got %q", mod.Clause)
	}
	if len(mod.Args) != 0 {
		t.Errorf("expected no args for null check, got %d", len(mod.Args))
	}
}

func TestWhereNullEQWithNullNegated(t *testing.T) {
	mod := WhereNullEQ("email", true, nullableString{Val: "", Valid: false})
	if mod.Clause != "email is not null" {
		t.Errorf("expected 'email is not null', got %q", mod.Clause)
	}
}

func TestWhereNullEQWithValue(t *testing.T) {
	mod := WhereNullEQ("email", false, nullableString{Val: "test@test.com", Valid: true})
	if mod.Clause != "email = ?" {
		t.Errorf("expected 'email = ?', got %q", mod.Clause)
	}
	if len(mod.Args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(mod.Args))
	}
}

func TestWhereNullEQWithValueNegated(t *testing.T) {
	mod := WhereNullEQ("email", true, nullableString{Val: "test@test.com", Valid: true})
	if mod.Clause != "email != ?" {
		t.Errorf("expected 'email != ?', got %q", mod.Clause)
	}
}

func TestWhereIsNull(t *testing.T) {
	mod := WhereIsNull("deleted_at")
	if mod.Clause != "deleted_at is null" {
		t.Errorf("expected 'deleted_at is null', got %q", mod.Clause)
	}
}

func TestWhereIsNotNull(t *testing.T) {
	mod := WhereIsNotNull("email")
	if mod.Clause != "email is not null" {
		t.Errorf("expected 'email is not null', got %q", mod.Clause)
	}
}

func TestWhereQueryModApply(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	mod := Where("age", GT, 18)
	mod.Apply(q)
	sql, args := queries.BuildQuery(q)
	if sql == "" {
		t.Error("Apply should produce a query")
	}
	if len(args) != 1 || args[0] != 18 {
		t.Errorf("expected args [18], got %v", args)
	}
}

type testModel struct {
	ID        int    `orm:"id"`
	Name      string `orm:"name"`
	CreatedAt string `orm:"created_at"`
}

func TestNonZeroDefaultSet(t *testing.T) {
	m := testModel{ID: 0, Name: "alice", CreatedAt: ""}
	defaults := []string{"id", "created_at"}
	result := NonZeroDefaultSet(defaults, &m)
	if len(result) != 0 {
		t.Errorf("expected empty result for zero-value fields, got %v", result)
	}

	m2 := testModel{ID: 5, Name: "bob", CreatedAt: "2024-01-01"}
	result2 := NonZeroDefaultSet(defaults, &m2)
	if len(result2) != 2 {
		t.Errorf("expected 2 non-zero defaults, got %v", result2)
	}
}

func TestNonZeroDefaultSetPartial(t *testing.T) {
	m := testModel{ID: 10, Name: "", CreatedAt: ""}
	defaults := []string{"id", "created_at"}
	result := NonZeroDefaultSet(defaults, &m)
	if len(result) != 1 || result[0] != "id" {
		t.Errorf("expected [id], got %v", result)
	}
}

func TestNonZeroDefaultSetPanicsOnMissing(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing column")
		}
	}()
	m := testModel{}
	NonZeroDefaultSet([]string{"nonexistent"}, &m)
}
