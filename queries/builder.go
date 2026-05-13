package queries

import (
	"fmt"
	"sort"
	"strings"
)

func BuildQuery(q *Query) (string, []any) {
	if q.rawSQL.sql != "" {
		return q.rawSQL.sql, q.rawSQL.args
	}
	if q.delete {
		return buildDeleteQuery(q)
	}
	if q.update != nil {
		return buildUpdateQuery(q)
	}
	return buildSelectQuery(q)
}

func buildSelectQuery(q *Query) (string, []any) {
	var buf strings.Builder
	var args []any

	if q.comment != "" {
		fmt.Fprintf(&buf, "-- %s\n", q.comment)
	}

	if len(q.withs) > 0 {
		buf.WriteString("WITH ")
		for i, w := range q.withs {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(w.clause)
			args = append(args, w.args...)
		}
		buf.WriteString(" ")
	}

	buf.WriteString("SELECT ")

	if q.count {
		buf.WriteString("COUNT(")
	}

	if q.distinct != "" {
		buf.WriteString("DISTINCT ")
		if q.distinct != "on" {
			buf.WriteString(q.distinct)
			buf.WriteString(" ")
		}
	}

	if q.dialect != nil && q.dialect.UseTopClause && q.limit != nil && !q.count {
		fmt.Fprintf(&buf, "TOP (%d) ", *q.limit)
	}

	if len(q.selectCols) > 0 {
		buf.WriteString(strings.Join(q.selectCols, ", "))
	} else {
		buf.WriteString("*")
	}

	if q.count {
		buf.WriteString(")")
	}

	if len(q.from) > 0 {
		buf.WriteString(" FROM ")
		buf.WriteString(strings.Join(q.from, ", "))
	}

	if len(q.joins) > 0 {
		for _, j := range q.joins {
			switch j.kind {
			case JoinInner:
				buf.WriteString(" INNER JOIN ")
			case JoinOuterLeft:
				buf.WriteString(" LEFT OUTER JOIN ")
			case JoinOuterRight:
				buf.WriteString(" RIGHT OUTER JOIN ")
			case JoinOuterFull:
				buf.WriteString(" FULL OUTER JOIN ")
			case JoinNatural:
				buf.WriteString(" NATURAL JOIN ")
			}
			buf.WriteString(j.clause)
			args = append(args, j.args...)
		}
	}

	whereSQL, whereArgs := buildWhereClause(q, len(args)+1)
	if whereSQL != "" {
		buf.WriteString(" WHERE ")
		buf.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	if len(q.groupBy) > 0 {
		buf.WriteString(" GROUP BY ")
		buf.WriteString(strings.Join(q.groupBy, ", "))
	}

	if len(q.having) > 0 {
		buf.WriteString(" HAVING ")
		for i, h := range q.having {
			if i > 0 {
				buf.WriteString(" AND ")
			}
			buf.WriteString(h.clause)
			args = append(args, h.args...)
		}
	}

	if len(q.orderBy) > 0 {
		buf.WriteString(" ORDER BY ")
		for i, o := range q.orderBy {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(o.clause)
			args = append(args, o.args...)
		}
	}

	if q.dialect != nil && q.dialect.UseTopClause {
		if q.offset > 0 || (q.limit != nil && len(q.orderBy) > 0) {
			if q.offset > 0 {
				fmt.Fprintf(&buf, " OFFSET %d ROWS", q.offset)
				if q.limit != nil {
					fmt.Fprintf(&buf, " FETCH NEXT %d ROWS ONLY", *q.limit)
				}
			}
		}
	} else {
		if q.limit != nil {
			fmt.Fprintf(&buf, " LIMIT %d", *q.limit)
		}
		if q.offset > 0 {
			fmt.Fprintf(&buf, " OFFSET %d", q.offset)
		}
	}

	if q.forlock != "" {
		buf.WriteString(" FOR ")
		buf.WriteString(q.forlock)
	}

	sql := buf.String()
	if q.dialect != nil && q.dialect.UseIndexPlaceholders {
		sql = convertPlaceholders(sql, 1)
	}
	return sql, args
}

func buildDeleteQuery(q *Query) (string, []any) {
	var buf strings.Builder
	var args []any

	if q.comment != "" {
		fmt.Fprintf(&buf, "-- %s\n", q.comment)
	}

	buf.WriteString("DELETE FROM ")
	if len(q.from) > 0 {
		buf.WriteString(strings.Join(q.from, ", "))
	}

	whereSQL, whereArgs := buildWhereClause(q, 1)
	if whereSQL != "" {
		buf.WriteString(" WHERE ")
		buf.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	sql := buf.String()
	if q.dialect != nil && q.dialect.UseIndexPlaceholders {
		sql = convertPlaceholders(sql, 1)
	}
	return sql, args
}

func buildUpdateQuery(q *Query) (string, []any) {
	var buf strings.Builder
	var args []any

	if q.comment != "" {
		fmt.Fprintf(&buf, "-- %s\n", q.comment)
	}

	buf.WriteString("UPDATE ")
	if len(q.from) > 0 {
		buf.WriteString(strings.Join(q.from, ", "))
	}
	buf.WriteString(" SET ")

	keys := make([]string, 0, len(q.update))
	for k := range q.update {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(k)
		buf.WriteString(" = ?")
		args = append(args, q.update[k])
	}

	whereSQL, whereArgs := buildWhereClause(q, len(args)+1)
	if whereSQL != "" {
		buf.WriteString(" WHERE ")
		buf.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	sql := buf.String()
	if q.dialect != nil && q.dialect.UseIndexPlaceholders {
		sql = convertPlaceholders(sql, 1)
	}
	return sql, args
}

func buildWhereClause(q *Query, startAt int) (string, []any) {
	if len(q.where) == 0 {
		return "", nil
	}

	var buf strings.Builder
	var args []any
	hasManualParens := false
	for _, w := range q.where {
		if w.kind == whereKindLeftParen || w.kind == whereKindRightParen {
			hasManualParens = true
			break
		}
	}

	for i, w := range q.where {
		switch w.kind {
		case whereKindLeftParen:
			if i > 0 && q.where[i-1].kind != whereKindLeftParen {
				if w.orSep {
					buf.WriteString(" OR ")
				} else {
					buf.WriteString(" AND ")
				}
			}
			buf.WriteString("(")
		case whereKindRightParen:
			buf.WriteString(")")
		case whereKindNormal:
			if i > 0 && q.where[i-1].kind != whereKindLeftParen {
				if w.orSep {
					buf.WriteString(" OR ")
				} else {
					buf.WriteString(" AND ")
				}
			}
			if !hasManualParens {
				buf.WriteString("(")
			}
			buf.WriteString(w.clause)
			if !hasManualParens {
				buf.WriteString(")")
			}
			args = append(args, w.args...)
		case whereKindIn:
			if i > 0 && q.where[i-1].kind != whereKindLeftParen {
				if w.orSep {
					buf.WriteString(" OR ")
				} else {
					buf.WriteString(" AND ")
				}
			}
			if len(w.args) == 0 {
				buf.WriteString("(1=0)")
			} else {
				buf.WriteString("(")
				buf.WriteString(w.clause)
				buf.WriteString(" IN (")
				placeholders := make([]string, len(w.args))
				for j := range w.args {
					placeholders[j] = "?"
				}
				buf.WriteString(strings.Join(placeholders, ","))
				buf.WriteString("))")
				args = append(args, w.args...)
			}
		case whereKindNotIn:
			if i > 0 && q.where[i-1].kind != whereKindLeftParen {
				if w.orSep {
					buf.WriteString(" OR ")
				} else {
					buf.WriteString(" AND ")
				}
			}
			if len(w.args) == 0 {
				buf.WriteString("(1=1)")
			} else {
				buf.WriteString("(")
				buf.WriteString(w.clause)
				buf.WriteString(" NOT IN (")
				placeholders := make([]string, len(w.args))
				for j := range w.args {
					placeholders[j] = "?"
				}
				buf.WriteString(strings.Join(placeholders, ","))
				buf.WriteString("))")
				args = append(args, w.args...)
			}
		}
	}

	return buf.String(), args
}

func convertPlaceholders(sql string, startAt int) string {
	var buf strings.Builder
	idx := startAt
	for i := 0; i < len(sql); i++ {
		if sql[i] == '?' {
			fmt.Fprintf(&buf, "$%d", idx)
			idx++
		} else {
			buf.WriteByte(sql[i])
		}
	}
	return buf.String()
}
