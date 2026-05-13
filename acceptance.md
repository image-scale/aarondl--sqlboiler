# Acceptance Criteria

## Task 1: Runtime support library (orm package) — DONE

## Task 2: Query building system

### Acceptance Criteria
- [ ] A Query struct accumulates SQL query state: select columns, from tables, joins, where clauses, group by, order by, having, limit, offset, distinct, for lock, comment
- [ ] A Dialect struct configures quoting characters (LQ, RQ), placeholder style (indexed vs ?), and SQL flavor flags (UseTopClause, UseOutputClause, etc.)
- [ ] QueryMod interface defines Apply(*Query) for composable modifiers
- [ ] Query mods exist for: SQL (raw), Select, From, Where, Or, InnerJoin, LeftOuterJoin, RightOuterJoin, FullOuterJoin, GroupBy, OrderBy, Having, Limit, Offset, Distinct, For, Comment, Load, With, WhereIn, WhereNotIn, OrIn, OrNotIn
- [ ] BuildQuery generates correct SELECT SQL with proper column quoting and placeholder conversion
- [ ] BuildQuery generates correct DELETE SQL with WHERE clauses
- [ ] BuildQuery generates correct UPDATE SQL with SET clauses and WHERE conditions
- [ ] Placeholder conversion handles indexed placeholders ($1, $2) for Postgres-style dialects
- [ ] WHERE clause generation supports AND/OR connectors, parentheses grouping, IN/NOT IN with argument expansion
- [ ] Raw SQL queries bypass the builder and pass through directly
- [ ] LIMIT/OFFSET handled correctly including MSSQL-style TOP/OFFSET FETCH syntax
- [ ] Join clauses generate correct INNER/LEFT/RIGHT/FULL OUTER JOIN SQL
- [ ] CTE (WITH) clauses are prepended correctly
- [ ] SQL comments are prepended as -- prefix
