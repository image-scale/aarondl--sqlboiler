package queries

type Dialect struct {
	LQ byte
	RQ byte

	UseIndexPlaceholders bool
	UseLastInsertID      bool
	UseSchema            bool
	UseDefaultKeyword    bool
	UseTopClause         bool
	UseOutputClause      bool
	UseCaseWhenExistsClause bool
}
