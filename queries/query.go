package queries

const (
	JoinInner     = iota
	JoinOuterLeft
	JoinOuterRight
	JoinOuterFull
	JoinNatural
)

const (
	whereKindNormal = iota
	whereKindLeftParen
	whereKindRightParen
	whereKindIn
	whereKindNotIn
)

type Query struct {
	dialect *Dialect
	rawSQL  rawSQL

	load     []string
	loadMods map[string]Applicator

	delete bool
	update map[string]any

	withs      []argClause
	selectCols []string
	count      bool
	from       []string
	joins      []join
	where      []where
	groupBy    []string
	orderBy    []argClause
	having     []argClause
	limit      *int
	offset     int
	forlock    string
	distinct   string
	comment    string
}

type Applicator interface {
	Apply(q *Query)
}

type rawSQL struct {
	sql  string
	args []any
}

type where struct {
	clause string
	args   []any
	kind   int
	orSep  bool
}

type argClause struct {
	clause string
	args   []any
}

type join struct {
	kind   int
	clause string
	args   []any
}

func Raw(query string, args ...any) *Query {
	return &Query{
		rawSQL: rawSQL{sql: query, args: args},
	}
}

func SetDialect(q *Query, d *Dialect) {
	q.dialect = d
}

func SetSQL(q *Query, sql string) {
	q.rawSQL.sql = sql
}

func SetArgs(q *Query, args ...any) {
	q.rawSQL.args = args
}

func SetLoad(q *Query, relationships ...string) {
	q.load = relationships
}

func AppendLoad(q *Query, rel string) {
	q.load = append(q.load, rel)
}

func SetLoadMods(q *Query, rel string, app Applicator) {
	if q.loadMods == nil {
		q.loadMods = make(map[string]Applicator)
	}
	q.loadMods[rel] = app
}

func SetSelect(q *Query, cols []string) {
	q.selectCols = cols
}

func GetSelect(q *Query) []string {
	return q.selectCols
}

func AppendSelect(q *Query, cols ...string) {
	q.selectCols = append(q.selectCols, cols...)
}

func SetDistinct(q *Query, clause string) {
	q.distinct = clause
}

func SetCount(q *Query) {
	q.count = true
}

func SetDelete(q *Query) {
	q.delete = true
}

func SetLimit(q *Query, limit int) {
	q.limit = &limit
}

func SetOffset(q *Query, offset int) {
	q.offset = offset
}

func SetFor(q *Query, clause string) {
	q.forlock = clause
}

func SetComment(q *Query, comment string) {
	q.comment = comment
}

func SetUpdate(q *Query, cols map[string]any) {
	q.update = cols
}

func AppendFrom(q *Query, from ...string) {
	q.from = append(q.from, from...)
}

func SetFrom(q *Query, from ...string) {
	q.from = from
}

func AppendInnerJoin(q *Query, clause string, args ...any) {
	q.joins = append(q.joins, join{kind: JoinInner, clause: clause, args: args})
}

func AppendLeftOuterJoin(q *Query, clause string, args ...any) {
	q.joins = append(q.joins, join{kind: JoinOuterLeft, clause: clause, args: args})
}

func AppendRightOuterJoin(q *Query, clause string, args ...any) {
	q.joins = append(q.joins, join{kind: JoinOuterRight, clause: clause, args: args})
}

func AppendFullOuterJoin(q *Query, clause string, args ...any) {
	q.joins = append(q.joins, join{kind: JoinOuterFull, clause: clause, args: args})
}

func AppendHaving(q *Query, clause string, args ...any) {
	q.having = append(q.having, argClause{clause: clause, args: args})
}

func AppendWhere(q *Query, clause string, args ...any) {
	q.where = append(q.where, where{clause: clause, args: args, kind: whereKindNormal})
}

func AppendIn(q *Query, clause string, args ...any) {
	q.where = append(q.where, where{clause: clause, args: args, kind: whereKindIn})
}

func AppendNotIn(q *Query, clause string, args ...any) {
	q.where = append(q.where, where{clause: clause, args: args, kind: whereKindNotIn})
}

func SetLastWhereAsOr(q *Query) {
	if len(q.where) > 0 {
		q.where[len(q.where)-1].orSep = true
	}
}

func SetLastInAsOr(q *Query) {
	if len(q.where) > 0 {
		q.where[len(q.where)-1].orSep = true
	}
}

func AppendWhereLeftParen(q *Query) {
	q.where = append(q.where, where{kind: whereKindLeftParen})
}

func AppendWhereRightParen(q *Query) {
	q.where = append(q.where, where{kind: whereKindRightParen})
}

func AppendGroupBy(q *Query, clause string) {
	q.groupBy = append(q.groupBy, clause)
}

func AppendOrderBy(q *Query, clause string, args ...any) {
	q.orderBy = append(q.orderBy, argClause{clause: clause, args: args})
}

func AppendWith(q *Query, clause string, args ...any) {
	q.withs = append(q.withs, argClause{clause: clause, args: args})
}

func GetDialect(q *Query) *Dialect {
	return q.dialect
}
