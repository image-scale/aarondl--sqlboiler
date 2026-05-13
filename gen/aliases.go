package gen

import (
	"strings"
	"unicode"

	"github.com/nl2repo/sqlboiler/schema"
)

type Aliases struct {
	Tables map[string]TableAlias
}

type TableAlias struct {
	UpPlural      string
	UpSingular    string
	DownPlural    string
	DownSingular  string
	Columns       map[string]string
	Relationships map[string]RelationshipAlias
}

type RelationshipAlias struct {
	Local   string
	Foreign string
}

func (a Aliases) Table(table string) TableAlias {
	ta, ok := a.Tables[table]
	if !ok {
		panic("gen: alias not found for table: " + table)
	}
	return ta
}

func (t TableAlias) Column(column string) string {
	c, ok := t.Columns[column]
	if !ok {
		panic("gen: alias not found for column: " + column)
	}
	return c
}

func (t TableAlias) Relationship(fkey string) RelationshipAlias {
	r, ok := t.Relationships[fkey]
	if !ok {
		panic("gen: alias not found for relationship: " + fkey)
	}
	return r
}

func FillAliases(a *Aliases, tables []schema.Table) {
	if a.Tables == nil {
		a.Tables = make(map[string]TableAlias)
	}

	for _, t := range tables {
		if t.IsJoinTable {
			continue
		}

		ta, exists := a.Tables[t.Name]
		if !exists {
			ta = TableAlias{
				Columns:       make(map[string]string),
				Relationships: make(map[string]RelationshipAlias),
			}
		}

		if ta.UpSingular == "" {
			singular := singularize(t.Name)
			ta.UpSingular = titleCase(singular)
			ta.UpPlural = titleCase(pluralize(singular))
			ta.DownSingular = camelCase(singular)
			ta.DownPlural = camelCase(pluralize(singular))
		}

		if ta.Columns == nil {
			ta.Columns = make(map[string]string)
		}
		for _, col := range t.Columns {
			if _, exists := ta.Columns[col.Name]; !exists {
				alias := titleCase(col.Name)
				if len(alias) > 0 && alias[0] >= '0' && alias[0] <= '9' {
					alias = "C" + alias
				}
				ta.Columns[col.Name] = alias
			}
		}

		if ta.Relationships == nil {
			ta.Relationships = make(map[string]RelationshipAlias)
		}
		for _, rel := range t.ToOneRelationships {
			if _, exists := ta.Relationships[rel.Name]; !exists {
				local, foreign := NameToOne(rel.Column, rel.ForeignTable, rel.ForeignColumn, rel.Unique)
				ta.Relationships[rel.Name] = RelationshipAlias{Local: local, Foreign: foreign}
			}
		}
		for _, rel := range t.ToManyRelationships {
			if _, exists := ta.Relationships[rel.Name]; !exists {
				if rel.ToJoinTable {
					local, foreign := NameToMany(rel.JoinLocalColumn, rel.ForeignTable, rel.JoinForeignColumn, t.Name)
					ta.Relationships[rel.Name] = RelationshipAlias{Local: local, Foreign: foreign}
				} else {
					local, foreign := NameToOne(rel.ForeignColumn, rel.ForeignTable, rel.Column, false)
					ta.Relationships[rel.Name] = RelationshipAlias{Local: local, Foreign: foreign}
				}
			}
		}

		a.Tables[t.Name] = ta
	}
}

func titleCase(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

func camelCase(s string) string {
	tc := titleCase(s)
	if tc == "" {
		return ""
	}
	i := 0
	for i < len(tc) && unicode.IsUpper(rune(tc[i])) {
		i++
	}
	if i > 1 {
		i--
	}
	return strings.ToLower(tc[:i]) + tc[i:]
}

func splitWords(s string) []string {
	var words []string
	var current strings.Builder
	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(unicode.ToLower(r))
		}
	}
	if current.Len() > 0 {
		words = append(words, current.String())
	}
	return words
}

func singularize(s string) string {
	if strings.HasSuffix(s, "ies") && len(s) > 3 {
		return s[:len(s)-3] + "y"
	}
	if strings.HasSuffix(s, "sses") {
		return s[:len(s)-2]
	}
	if strings.HasSuffix(s, "ses") && !strings.HasSuffix(s, "sses") {
		base := s[:len(s)-2]
		if strings.HasSuffix(base, "s") || strings.HasSuffix(base, "x") ||
			strings.HasSuffix(base, "z") || strings.HasSuffix(base, "ch") ||
			strings.HasSuffix(base, "sh") {
			return base
		}
		return s[:len(s)-1]
	}
	if strings.HasSuffix(s, "xes") || strings.HasSuffix(s, "zes") ||
		strings.HasSuffix(s, "ches") || strings.HasSuffix(s, "shes") {
		return s[:len(s)-2]
	}
	if strings.HasSuffix(s, "s") && !strings.HasSuffix(s, "ss") &&
		!strings.HasSuffix(s, "us") && !strings.HasSuffix(s, "is") {
		return s[:len(s)-1]
	}
	return s
}

func pluralize(s string) string {
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") ||
		strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") && len(s) > 1 {
		ch := s[len(s)-2]
		if ch != 'a' && ch != 'e' && ch != 'i' && ch != 'o' && ch != 'u' {
			return s[:len(s)-1] + "ies"
		}
	}
	return s + "s"
}
