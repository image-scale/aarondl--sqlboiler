package gen

import "github.com/nl2repo/sqlboiler/importers"

type TagCase string

const (
	TagCaseCamel TagCase = "camel"
	TagCaseSnake TagCase = "snake"
	TagCaseTitle TagCase = "title"
	TagCaseAlias TagCase = "alias"
)

type Config struct {
	DriverName  string
	DriverConfig map[string]any

	PkgName   string
	OutFolder string

	TemplateDirs []string
	Tags         []string
	TagIgnore    []string

	Debug bool
	Wipe  bool

	AddGlobal            bool
	AddPanic             bool
	AddSoftDeletes       bool
	AddEnumTypes         bool
	NoContext            bool
	NoTests              bool
	NoHooks              bool
	NoAutoTimestamps     bool
	NoRowsAffected       bool
	NoDriverTemplates    bool
	NoBackReferencing    bool
	AlwaysWrapErrors     bool

	StructTagCases StructTagCases
	RelationTag    string

	Imports     importers.Collection
	Aliases     Aliases
	TypeReplaces []TypeReplace
	AutoColumns  AutoColumns
	Inflections  Inflections
}

type StructTagCases struct {
	Json TagCase
	Yaml TagCase
	Toml TagCase
	Boil TagCase
}

type TypeReplace struct {
	Tables  []string
	Match   ColumnMatch
	Replace ColumnReplace
	Imports importers.Set
}

type ColumnMatch struct {
	Name     string
	Type     string
	DBType   string
	Nullable *bool
}

type ColumnReplace struct {
	Type string
}

type AutoColumns struct {
	Created string
	Updated string
	Deleted string
}

type Inflections struct {
	Plural       map[string]string
	PluralExact  map[string]string
	Singular     map[string]string
	SingularExact map[string]string
	Irregular    map[string]string
}
