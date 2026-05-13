package qm

import (
	"testing"

	"github.com/nl2repo/sqlboiler/queries"
)

func TestSelectMod(t *testing.T) {
	q := &queries.Query{}
	Select("id", "name").Apply(q)
	cols := queries.GetSelect(q)
	if len(cols) != 2 || cols[0] != "id" || cols[1] != "name" {
		t.Errorf("Select mod failed, got %v", cols)
	}
}

func TestFromMod(t *testing.T) {
	q := &queries.Query{}
	From("users").Apply(q)
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("From mod should produce query")
	}
}

func TestWhereMod(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	Where("id = ?", 1).Apply(q)
	sql, args := queries.BuildQuery(q)
	if sql == "" {
		t.Error("Where mod should produce query")
	}
	if len(args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(args))
	}
}

func TestOrMod(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	Where("a = ?", 1).Apply(q)
	Or("b = ?", 2).Apply(q)
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("Or mod should produce query")
	}
}

func TestJoinMods(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	InnerJoin("orders ON users.id = orders.uid").Apply(q)
	LeftOuterJoin("carts ON users.id = carts.uid").Apply(q)
	RightOuterJoin("reviews ON users.id = reviews.uid").Apply(q)
	FullOuterJoin("logs ON users.id = logs.uid").Apply(q)
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("Join mods should produce query")
	}
}

func TestGroupByOrderByHavingMods(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	GroupBy("status").Apply(q)
	Having("count(*) > ?", 1).Apply(q)
	OrderBy("name ASC").Apply(q)
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("GroupBy/Having/OrderBy mods should produce query")
	}
}

func TestLimitOffsetMods(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	Limit(10).Apply(q)
	Offset(5).Apply(q)
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("Limit/Offset mods should produce query")
	}
}

func TestDistinctMod(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	Distinct("on").Apply(q)
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("Distinct mod should produce query")
	}
}

func TestForMod(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	For("UPDATE").Apply(q)
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("For mod should produce query")
	}
}

func TestCommentMod(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	Comment("trace-id-123").Apply(q)
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("Comment mod should produce query")
	}
}

func TestLoadMod(t *testing.T) {
	q := &queries.Query{}
	Load("Posts").Apply(q)
	Load("Posts.Comments", Where("visible = ?", true)).Apply(q)
	// verify it doesn't panic
}

func TestWithMod(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	With("cte AS (SELECT 1)").Apply(q)
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("With mod should produce query")
	}
}

func TestWhereInMod(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	WhereIn("id", 1, 2, 3).Apply(q)
	sql, args := queries.BuildQuery(q)
	if sql == "" {
		t.Error("WhereIn mod should produce query")
	}
	if len(args) != 3 {
		t.Errorf("expected 3 args, got %d", len(args))
	}
}

func TestWhereNotInMod(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	WhereNotIn("id", 1, 2).Apply(q)
	sql, args := queries.BuildQuery(q)
	if sql == "" {
		t.Error("WhereNotIn mod should produce query")
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d", len(args))
	}
}

func TestOrInMod(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	WhereIn("id", 1).Apply(q)
	OrIn("name", "a").Apply(q)
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("OrIn mod should produce query")
	}
}

func TestOrNotInMod(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	WhereNotIn("id", 1).Apply(q)
	OrNotIn("name", "a").Apply(q)
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("OrNotIn mod should produce query")
	}
}

func TestSQLMod(t *testing.T) {
	q := &queries.Query{}
	SQL("SELECT 1", 42).Apply(q)
	sql, args := queries.BuildQuery(q)
	if sql != "SELECT 1" {
		t.Errorf("expected raw SQL, got %q", sql)
	}
	if len(args) != 1 || args[0] != 42 {
		t.Errorf("expected args [42], got %v", args)
	}
}

func TestRelsConcatenation(t *testing.T) {
	r := Rels("Posts", "Comments", "Author")
	if r != "Posts.Comments.Author" {
		t.Errorf("Rels should join with ., got %q", r)
	}
}

func TestApplyMultipleMods(t *testing.T) {
	q := &queries.Query{}
	Apply(q, From("users"), Select("id"), Limit(5))
	sql, _ := queries.BuildQuery(q)
	if sql == "" {
		t.Error("Apply should produce query")
	}
}

func TestExprMod(t *testing.T) {
	q := &queries.Query{}
	queries.AppendFrom(q, "users")
	Expr(Where("a = ?", 1), Or("b = ?", 2)).Apply(q)
	sql, args := queries.BuildQuery(q)
	if sql == "" {
		t.Error("Expr mod should produce query")
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d", len(args))
	}
}
