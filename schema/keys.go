package schema

type PrimaryKey struct {
	Name    string
	Columns []string
}

type ForeignKey struct {
	Table                 string
	Name                  string
	Column                string
	Nullable              bool
	Unique                bool
	ForeignTable          string
	ForeignColumn         string
	ForeignColumnNullable bool
	ForeignColumnUnique   bool
}

type SQLColumnDef struct {
	Name string
	Type string
}

type SQLColumnDefs []SQLColumnDef

func (defs SQLColumnDefs) Names() []string {
	names := make([]string, len(defs))
	for i, d := range defs {
		names[i] = d.Name
	}
	return names
}

func (defs SQLColumnDefs) Types() []string {
	types := make([]string, len(defs))
	for i, d := range defs {
		types[i] = d.Type
	}
	return types
}

func MakeSQLColumnDefs(cols []Column, names []string) SQLColumnDefs {
	defs := make(SQLColumnDefs, len(names))
	colMap := make(map[string]Column, len(cols))
	for _, c := range cols {
		colMap[c.Name] = c
	}
	for i, n := range names {
		c := colMap[n]
		defs[i] = SQLColumnDef{Name: n, Type: c.Type}
	}
	return defs
}
