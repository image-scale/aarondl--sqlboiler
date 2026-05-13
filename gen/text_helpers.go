package gen

import (
	"strings"

	"github.com/nl2repo/sqlboiler/schema"
)

var identifierSuffixes = []string{"_id", "_uuid", "_guid", "_oid"}

func NameToOne(fkColumn, foreignTable, foreignColumn string, unique bool) (localFn, foreignFn string) {
	localName := trimSuffixes(fkColumn)
	foreignSingular := singularize(foreignTable)

	if localName == foreignSingular {
		foreignFn = titleCase(foreignSingular)
	} else {
		foreignFn = titleCase(localName) + titleCase(foreignSingular)
	}

	if unique {
		localFn = titleCase(singularize(foreignTable))
	} else {
		localFn = titleCase(pluralize(singularize(foreignTable)))
	}

	return localFn, foreignFn
}

func NameToMany(localFKCol, foreignTable, foreignFKCol, localTable string) (lhsFn, rhsFn string) {
	localTrimmed := trimSuffixes(localFKCol)
	foreignTrimmed := trimSuffixes(foreignFKCol)

	foreignSingular := singularize(foreignTable)
	localSingular := singularize(localTable)

	if foreignTrimmed == foreignSingular {
		lhsFn = titleCase(pluralize(foreignSingular))
	} else {
		lhsFn = titleCase(foreignTrimmed) + titleCase(pluralize(foreignSingular))
	}

	if localTrimmed == localSingular {
		rhsFn = titleCase(pluralize(localSingular))
	} else {
		rhsFn = titleCase(localTrimmed) + titleCase(pluralize(localSingular))
	}

	return lhsFn, rhsFn
}

func trimSuffixes(s string) string {
	for _, suffix := range identifierSuffixes {
		if strings.HasSuffix(s, suffix) {
			return s[:len(s)-len(suffix)]
		}
	}
	return s
}

func IsPrimitive(typ string) bool {
	switch typ {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "byte", "rune", "string", "bool":
		return true
	}
	return false
}

func IsNullPrimitive(typ string) bool {
	if !strings.HasPrefix(typ, "null.") {
		return false
	}
	inner := typ[5:]
	switch strings.ToLower(inner) {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "byte", "string", "bool":
		return true
	}
	return false
}

func ConvertNullToPrimitive(typ string) string {
	if !strings.HasPrefix(typ, "null.") {
		return typ
	}
	return strings.ToLower(typ[5:])
}

func UsesPrimitives(tables []schema.Table, table, column, foreignTable, foreignColumn string) bool {
	var localType, foreignType string
	for _, t := range tables {
		if t.Name == table {
			for _, c := range t.Columns {
				if c.Name == column {
					localType = c.Type
					break
				}
			}
		}
		if t.Name == foreignTable {
			for _, c := range t.Columns {
				if c.Name == foreignColumn {
					foreignType = c.Type
					break
				}
			}
		}
	}
	return IsPrimitive(localType) && IsPrimitive(foreignType)
}
