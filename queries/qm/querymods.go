package qm

import (
	"strings"

	"github.com/nl2repo/sqlboiler/queries"
)

type QueryMod interface {
	Apply(q *queries.Query)
}

type QueryModFunc func(q *queries.Query)

func (f QueryModFunc) Apply(q *queries.Query) {
	f(q)
}

func Apply(q *queries.Query, mods ...QueryMod) {
	for _, m := range mods {
		m.Apply(q)
	}
}

func SQL(sql string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.SetSQL(q, sql)
		queries.SetArgs(q, args...)
	})
}

func Load(relationship string, mods ...QueryMod) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendLoad(q, relationship)
		if len(mods) > 0 {
			queries.SetLoadMods(q, relationship, queryMods(mods))
		}
	})
}

type queryMods []QueryMod

func (qms queryMods) Apply(q *queries.Query) {
	for _, m := range qms {
		m.Apply(q)
	}
}

func Select(columns ...string) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendSelect(q, columns...)
	})
}

func From(table string) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendFrom(q, table)
	})
}

func Where(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendWhere(q, clause, args...)
	})
}

func And(clause string, args ...any) QueryMod {
	return Where(clause, args...)
}

func Or(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendWhere(q, clause, args...)
		queries.SetLastWhereAsOr(q)
	})
}

func Or2(mod QueryMod) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		mod.Apply(q)
		queries.SetLastWhereAsOr(q)
	})
}

func Expr(mods ...QueryMod) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendWhereLeftParen(q)
		for _, m := range mods {
			m.Apply(q)
		}
		queries.AppendWhereRightParen(q)
	})
}

func WhereIn(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendIn(q, clause, args...)
	})
}

func AndIn(clause string, args ...any) QueryMod {
	return WhereIn(clause, args...)
}

func OrIn(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendIn(q, clause, args...)
		queries.SetLastInAsOr(q)
	})
}

func WhereNotIn(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendNotIn(q, clause, args...)
	})
}

func AndNotIn(clause string, args ...any) QueryMod {
	return WhereNotIn(clause, args...)
}

func OrNotIn(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendNotIn(q, clause, args...)
		queries.SetLastInAsOr(q)
	})
}

func InnerJoin(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendInnerJoin(q, clause, args...)
	})
}

func LeftOuterJoin(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendLeftOuterJoin(q, clause, args...)
	})
}

func RightOuterJoin(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendRightOuterJoin(q, clause, args...)
	})
}

func FullOuterJoin(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendFullOuterJoin(q, clause, args...)
	})
}

func Distinct(clause string) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.SetDistinct(q, clause)
	})
}

func With(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendWith(q, clause, args...)
	})
}

func GroupBy(clause string) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendGroupBy(q, clause)
	})
}

func OrderBy(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendOrderBy(q, clause, args...)
	})
}

func Having(clause string, args ...any) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.AppendHaving(q, clause, args...)
	})
}

func Limit(limit int) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.SetLimit(q, limit)
	})
}

func Offset(offset int) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.SetOffset(q, offset)
	})
}

func For(clause string) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.SetFor(q, clause)
	})
}

func Comment(comment string) QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		queries.SetComment(q, comment)
	})
}

func WithDeleted() QueryMod {
	return QueryModFunc(func(q *queries.Query) {
		// placeholder for soft delete removal
	})
}

func Rels(r ...string) string {
	return strings.Join(r, ".")
}
