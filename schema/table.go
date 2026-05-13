package schema

type Table struct {
	Name                string
	SchemaName          string
	Columns             []Column
	PKey                *PrimaryKey
	FKeys               []ForeignKey
	IsJoinTable         bool
	ToOneRelationships  []ToOneRelationship
	ToManyRelationships []ToManyRelationship
	IsView              bool
	ViewCapabilities    ViewCapabilities
}

type ViewCapabilities struct {
	CanInsert bool
	CanUpsert bool
}

func GetTable(tables []Table, name string) Table {
	for _, t := range tables {
		if t.Name == name {
			return t
		}
	}
	panic("schema: table not found: " + name)
}

func (t Table) GetColumn(name string) Column {
	for _, c := range t.Columns {
		if c.Name == name {
			return c
		}
	}
	panic("schema: column not found: " + name)
}

func (t Table) CanLastInsertID() bool {
	if t.PKey == nil || len(t.PKey.Columns) != 1 {
		return false
	}
	col := t.GetColumn(t.PKey.Columns[0])
	return col.Default != "" && isIntegerType(col.Type)
}

func (t Table) CanSoftDelete(deleteColumn string) bool {
	for _, c := range t.Columns {
		if c.Name == deleteColumn && c.Nullable && c.Type == "null.Time" {
			return true
		}
	}
	return false
}

func isIntegerType(typ string) bool {
	switch typ {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return true
	}
	return false
}
