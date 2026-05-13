package schema

import (
	"fmt"
	"sync"
)

type DBInfo struct {
	Schema  string
	Tables  []Table
	Dialect Dialect
}

type Dialect struct {
	LQ byte
	RQ byte

	UseIndexPlaceholders    bool
	UseLastInsertID         bool
	UseSchema               bool
	UseDefaultKeyword       bool
	UseTopClause            bool
	UseOutputClause         bool
	UseCaseWhenExistsClause bool
}

type Driver interface {
	Assemble(config Config) (*DBInfo, error)
	Templates() (map[string]string, error)
	Imports() (map[string]any, error)
}

type Constructor interface {
	TableNames(schema string, whitelist, blacklist []string) ([]string, error)
	Columns(schema, tableName string, whitelist, blacklist []string) ([]Column, error)
	PrimaryKeyInfo(schema, tableName string) (*PrimaryKey, error)
	ForeignKeyInfo(schema, tableName string) ([]ForeignKey, error)
	TranslateColumnType(Column) Column
}

var (
	registryMu sync.RWMutex
	registry   = make(map[string]Driver)
)

func RegisterDriver(name string, d Driver) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = d
}

func GetDriver(name string) Driver {
	registryMu.RLock()
	defer registryMu.RUnlock()
	d, ok := registry[name]
	if !ok {
		panic(fmt.Sprintf("schema: driver %q not registered", name))
	}
	return d
}

func TablesFromList(list []string) []string {
	var tables []string
	for _, item := range list {
		hasDot := false
		for _, ch := range item {
			if ch == '.' {
				hasDot = true
				break
			}
		}
		if !hasDot {
			tables = append(tables, item)
		}
	}
	return tables
}

func ColumnsFromList(list []string, tableName string) []string {
	var cols []string
	for _, item := range list {
		for i, ch := range item {
			if ch == '.' {
				tbl := item[:i]
				col := item[i+1:]
				if tbl == tableName || tbl == "*" {
					cols = append(cols, col)
				}
				break
			}
		}
	}
	return cols
}
