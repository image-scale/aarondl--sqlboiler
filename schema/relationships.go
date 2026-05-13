package schema

type ToOneRelationship struct {
	Name                  string
	Table                 string
	Column                string
	Nullable              bool
	Unique                bool
	ForeignTable          string
	ForeignColumn         string
	ForeignColumnNullable bool
	ForeignColumnUnique   bool
}

type ToManyRelationship struct {
	Name                       string
	Table                      string
	Column                     string
	Nullable                   bool
	Unique                     bool
	ForeignTable               string
	ForeignColumn              string
	ForeignColumnNullable      bool
	ForeignColumnUnique        bool
	ToJoinTable                bool
	JoinTable                  string
	JoinLocalFKeyName          string
	JoinLocalColumn            string
	JoinLocalColumnNullable    bool
	JoinLocalColumnUnique      bool
	JoinForeignFKeyName        string
	JoinForeignColumn          string
	JoinForeignColumnNullable  bool
	JoinForeignColumnUnique    bool
}

func DetectToOneRelationships(table string, tables []Table) []ToOneRelationship {
	var rels []ToOneRelationship
	for _, t := range tables {
		if t.IsJoinTable || t.Name == table {
			continue
		}
		for _, fk := range t.FKeys {
			if fk.ForeignTable != table {
				continue
			}
			if !fk.Unique {
				continue
			}
			rels = append(rels, ToOneRelationship{
				Name:                  fk.Name,
				Table:                 table,
				Column:                fk.ForeignColumn,
				Nullable:              fk.ForeignColumnNullable,
				Unique:                fk.ForeignColumnUnique,
				ForeignTable:          t.Name,
				ForeignColumn:         fk.Column,
				ForeignColumnNullable: fk.Nullable,
				ForeignColumnUnique:   fk.Unique,
			})
		}
	}
	return rels
}

func DetectToManyRelationships(table string, tables []Table) []ToManyRelationship {
	var rels []ToManyRelationship
	for _, t := range tables {
		if t.Name == table {
			continue
		}

		if t.IsJoinTable {
			rels = append(rels, detectManyToMany(table, t, tables)...)
			continue
		}

		for _, fk := range t.FKeys {
			if fk.ForeignTable != table {
				continue
			}
			if fk.Unique {
				continue
			}
			rels = append(rels, ToManyRelationship{
				Name:                  fk.Name,
				Table:                 table,
				Column:                fk.ForeignColumn,
				Nullable:              fk.ForeignColumnNullable,
				Unique:                fk.ForeignColumnUnique,
				ForeignTable:          t.Name,
				ForeignColumn:         fk.Column,
				ForeignColumnNullable: fk.Nullable,
				ForeignColumnUnique:   fk.Unique,
			})
		}
	}
	return rels
}

func detectManyToMany(table string, joinTable Table, tables []Table) []ToManyRelationship {
	var rels []ToManyRelationship
	var localFK, foreignFK *ForeignKey

	for i, fk := range joinTable.FKeys {
		if fk.ForeignTable == table {
			localFK = &joinTable.FKeys[i]
		}
	}
	if localFK == nil {
		return nil
	}

	for i, fk := range joinTable.FKeys {
		if fk.ForeignTable != table || fk.Column != localFK.Column {
			foreignFK = &joinTable.FKeys[i]
			break
		}
	}
	if foreignFK == nil {
		return nil
	}

	rels = append(rels, ToManyRelationship{
		Name:                      localFK.Name,
		Table:                     table,
		Column:                    localFK.ForeignColumn,
		Nullable:                  localFK.ForeignColumnNullable,
		Unique:                    localFK.ForeignColumnUnique,
		ForeignTable:              foreignFK.ForeignTable,
		ForeignColumn:             foreignFK.ForeignColumn,
		ForeignColumnNullable:     foreignFK.ForeignColumnNullable,
		ForeignColumnUnique:       foreignFK.ForeignColumnUnique,
		ToJoinTable:               true,
		JoinTable:                 joinTable.Name,
		JoinLocalFKeyName:         localFK.Name,
		JoinLocalColumn:           localFK.Column,
		JoinLocalColumnNullable:   localFK.Nullable,
		JoinLocalColumnUnique:     localFK.Unique,
		JoinForeignFKeyName:       foreignFK.Name,
		JoinForeignColumn:         foreignFK.Column,
		JoinForeignColumnNullable: foreignFK.Nullable,
		JoinForeignColumnUnique:   foreignFK.Unique,
	})
	return rels
}
