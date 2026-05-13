package queries

import (
	"strings"
	"testing"
)

func TestBuildSelectBasic(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	AppendSelect(q, "id", "name")

	sql, args := BuildQuery(q)
	if !strings.Contains(sql, "SELECT id, name") {
		t.Errorf("expected SELECT id, name, got %q", sql)
	}
	if !strings.Contains(sql, "FROM users") {
		t.Errorf("expected FROM users, got %q", sql)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

func TestBuildSelectStar(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "SELECT *") {
		t.Errorf("expected SELECT *, got %q", sql)
	}
}

func TestBuildSelectWithWhere(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	AppendWhere(q, "age > ?", 18)
	AppendWhere(q, "name = ?", "alice")

	sql, args := BuildQuery(q)
	if !strings.Contains(sql, "WHERE") {
		t.Errorf("expected WHERE clause, got %q", sql)
	}
	if !strings.Contains(sql, "age > ?") {
		t.Errorf("expected age > ?, got %q", sql)
	}
	if !strings.Contains(sql, "AND") {
		t.Errorf("expected AND between where clauses, got %q", sql)
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d: %v", len(args), args)
	}
}

func TestBuildSelectOrWhere(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	AppendWhere(q, "age > ?", 18)
	AppendWhere(q, "name = ?", "bob")
	SetLastWhereAsOr(q)

	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "OR") {
		t.Errorf("expected OR separator, got %q", sql)
	}
}

func TestBuildSelectWithJoins(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	AppendInnerJoin(q, "orders ON users.id = orders.user_id")
	AppendLeftOuterJoin(q, "payments ON orders.id = payments.order_id")

	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "INNER JOIN orders ON users.id = orders.user_id") {
		t.Errorf("expected INNER JOIN, got %q", sql)
	}
	if !strings.Contains(sql, "LEFT OUTER JOIN payments ON orders.id = payments.order_id") {
		t.Errorf("expected LEFT OUTER JOIN, got %q", sql)
	}
}

func TestBuildSelectRightAndFullJoin(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "a")
	AppendRightOuterJoin(q, "b ON a.id = b.a_id")
	AppendFullOuterJoin(q, "c ON a.id = c.a_id")

	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "RIGHT OUTER JOIN b ON a.id = b.a_id") {
		t.Errorf("expected RIGHT OUTER JOIN, got %q", sql)
	}
	if !strings.Contains(sql, "FULL OUTER JOIN c ON a.id = c.a_id") {
		t.Errorf("expected FULL OUTER JOIN, got %q", sql)
	}
}

func TestBuildSelectGroupByHavingOrderBy(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "orders")
	AppendSelect(q, "user_id", "count(*)")
	AppendGroupBy(q, "user_id")
	AppendHaving(q, "count(*) > ?", 5)
	AppendOrderBy(q, "user_id ASC")

	sql, args := BuildQuery(q)
	if !strings.Contains(sql, "GROUP BY user_id") {
		t.Errorf("expected GROUP BY, got %q", sql)
	}
	if !strings.Contains(sql, "HAVING count(*) > ?") {
		t.Errorf("expected HAVING, got %q", sql)
	}
	if !strings.Contains(sql, "ORDER BY user_id ASC") {
		t.Errorf("expected ORDER BY, got %q", sql)
	}
	if len(args) != 1 || args[0] != 5 {
		t.Errorf("expected args [5], got %v", args)
	}
}

func TestBuildSelectLimitOffset(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	SetLimit(q, 10)
	SetOffset(q, 20)

	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "LIMIT 10") {
		t.Errorf("expected LIMIT 10, got %q", sql)
	}
	if !strings.Contains(sql, "OFFSET 20") {
		t.Errorf("expected OFFSET 20, got %q", sql)
	}
}

func TestBuildSelectDistinct(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	SetDistinct(q, "on")
	AppendSelect(q, "name")

	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "DISTINCT") {
		t.Errorf("expected DISTINCT, got %q", sql)
	}
}

func TestBuildSelectForUpdate(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	SetFor(q, "UPDATE")

	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "FOR UPDATE") {
		t.Errorf("expected FOR UPDATE, got %q", sql)
	}
}

func TestBuildSelectComment(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	SetComment(q, "user query")

	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "-- user query") {
		t.Errorf("expected comment, got %q", sql)
	}
}

func TestBuildSelectCTE(t *testing.T) {
	q := &Query{}
	AppendWith(q, "active_users AS (SELECT * FROM users WHERE active = ?)", true)
	AppendFrom(q, "active_users")

	sql, args := BuildQuery(q)
	if !strings.Contains(sql, "WITH active_users AS (SELECT * FROM users WHERE active = ?)") {
		t.Errorf("expected WITH clause, got %q", sql)
	}
	if len(args) != 1 || args[0] != true {
		t.Errorf("expected args [true], got %v", args)
	}
}

func TestBuildRawQuery(t *testing.T) {
	q := Raw("SELECT * FROM users WHERE id = ?", 42)
	sql, args := BuildQuery(q)
	if sql != "SELECT * FROM users WHERE id = ?" {
		t.Errorf("expected raw SQL, got %q", sql)
	}
	if len(args) != 1 || args[0] != 42 {
		t.Errorf("expected args [42], got %v", args)
	}
}

func TestBuildDeleteQuery(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	AppendWhere(q, "id = ?", 1)
	SetDelete(q)

	sql, args := BuildQuery(q)
	if !strings.Contains(sql, "DELETE FROM users") {
		t.Errorf("expected DELETE FROM, got %q", sql)
	}
	if !strings.Contains(sql, "WHERE") {
		t.Errorf("expected WHERE clause, got %q", sql)
	}
	if len(args) != 1 || args[0] != 1 {
		t.Errorf("expected args [1], got %v", args)
	}
}

func TestBuildUpdateQuery(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	SetUpdate(q, map[string]any{"name": "bob", "age": 30})
	AppendWhere(q, "id = ?", 1)

	sql, args := BuildQuery(q)
	if !strings.Contains(sql, "UPDATE users SET") {
		t.Errorf("expected UPDATE users SET, got %q", sql)
	}
	if !strings.Contains(sql, "age = ?") || !strings.Contains(sql, "name = ?") {
		t.Errorf("expected column assignments, got %q", sql)
	}
	if !strings.Contains(sql, "WHERE") {
		t.Errorf("expected WHERE clause, got %q", sql)
	}
	if len(args) != 3 {
		t.Errorf("expected 3 args, got %d: %v", len(args), args)
	}
}

func TestPlaceholderConversion(t *testing.T) {
	d := &Dialect{UseIndexPlaceholders: true, LQ: '"', RQ: '"'}
	q := &Query{}
	SetDialect(q, d)
	AppendFrom(q, "users")
	AppendWhere(q, "age > ?", 21)
	AppendWhere(q, "name = ?", "x")

	sql, _ := BuildQuery(q)
	if strings.Contains(sql, "?") {
		t.Errorf("should not have ? placeholders with indexed mode, got %q", sql)
	}
	if !strings.Contains(sql, "$1") || !strings.Contains(sql, "$2") {
		t.Errorf("expected $1 and $2 placeholders, got %q", sql)
	}
}

func TestWhereInClause(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	AppendIn(q, "id", 1, 2, 3)

	sql, args := BuildQuery(q)
	if !strings.Contains(sql, "IN") {
		t.Errorf("expected IN clause, got %q", sql)
	}
	if len(args) != 3 {
		t.Errorf("expected 3 args, got %d", len(args))
	}
}

func TestWhereNotInClause(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	AppendNotIn(q, "id", 4, 5)

	sql, args := BuildQuery(q)
	if !strings.Contains(sql, "NOT IN") {
		t.Errorf("expected NOT IN clause, got %q", sql)
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d", len(args))
	}
}

func TestWhereInEmptyArgs(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	AppendIn(q, "id")

	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "1=0") {
		t.Errorf("empty IN should produce 1=0, got %q", sql)
	}
}

func TestWhereNotInEmptyArgs(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	AppendNotIn(q, "id")

	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "1=1") {
		t.Errorf("empty NOT IN should produce 1=1, got %q", sql)
	}
}

func TestMSSQLTopClause(t *testing.T) {
	d := &Dialect{UseTopClause: true, LQ: '[', RQ: ']'}
	q := &Query{}
	SetDialect(q, d)
	AppendFrom(q, "users")
	SetLimit(q, 5)

	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "TOP (5)") {
		t.Errorf("expected TOP (5) for MSSQL, got %q", sql)
	}
}

func TestMSSQLOffsetFetch(t *testing.T) {
	d := &Dialect{UseTopClause: true, LQ: '[', RQ: ']'}
	q := &Query{}
	SetDialect(q, d)
	AppendFrom(q, "users")
	AppendOrderBy(q, "id")
	SetLimit(q, 10)
	SetOffset(q, 20)

	sql, _ := BuildQuery(q)
	if !strings.Contains(sql, "OFFSET 20 ROWS") {
		t.Errorf("expected OFFSET ROWS for MSSQL, got %q", sql)
	}
	if !strings.Contains(sql, "FETCH NEXT 10 ROWS ONLY") {
		t.Errorf("expected FETCH NEXT for MSSQL, got %q", sql)
	}
}

func TestWhereParenGrouping(t *testing.T) {
	q := &Query{}
	AppendFrom(q, "users")
	AppendWhereLeftParen(q)
	AppendWhere(q, "age > ?", 18)
	AppendWhere(q, "age < ?", 65)
	SetLastWhereAsOr(q)
	AppendWhereRightParen(q)

	sql, args := BuildQuery(q)
	if !strings.Contains(sql, "(") || !strings.Contains(sql, ")") {
		t.Errorf("expected parentheses in WHERE, got %q", sql)
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d: %v", len(args), args)
	}
}
