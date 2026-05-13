package gen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"text/template"

	"github.com/nl2repo/sqlboiler/importers"
	"github.com/nl2repo/sqlboiler/schema"
)

func TestNewState(t *testing.T) {
	cfg := &Config{PkgName: "models", DriverName: "postgres"}
	s, err := NewState(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if s.Config.PkgName != "models" {
		t.Error("expected PkgName=models")
	}
}

func TestExecuteTemplate(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse("Hello {{.Name}}"))
	data := struct{ Name string }{"World"}
	err := ExecuteTemplate(tmpl, "test", data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecuteTemplateToBytes(t *testing.T) {
	tmpl := template.Must(template.New("greeting").Parse("Hi {{.}}"))
	out, err := ExecuteTemplateToBytes(tmpl, "greeting", "Alice")
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "Hi Alice" {
		t.Errorf("expected 'Hi Alice', got %q", out)
	}
}

func TestFormatGoSource(t *testing.T) {
	src := []byte("package main\n\nfunc main(){fmt.Println(\"hello\")}\n")
	formatted, err := FormatGoSource(src)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(formatted), "package main") {
		t.Error("formatted source should contain package declaration")
	}
}

func TestFormatGoSourceError(t *testing.T) {
	src := []byte("this is not valid go code {{{")
	_, err := FormatGoSource(src)
	if err == nil {
		t.Error("expected format error for invalid source")
	}
}

func TestWriteOutput(t *testing.T) {
	dir := t.TempDir()
	content := []byte("package test\n\nvar x = 1\n")
	err := WriteOutput(dir, "test.go", content, true)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "test.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "package test") {
		t.Error("output file should contain package declaration")
	}
}

func TestWriteOutputCreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "sub", "dir")
	content := []byte("package nested\n")
	err := WriteOutput(dir, "file.go", content, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuildFileHeader(t *testing.T) {
	imps := importers.Set{
		Standard:   importers.List{`"fmt"`, `"os"`},
		ThirdParty: importers.List{`"github.com/pkg/errors"`},
	}
	header := BuildFileHeader("models", imps)
	s := string(header)
	if !strings.Contains(s, "DO NOT EDIT") {
		t.Error("expected no-edit disclaimer")
	}
	if !strings.Contains(s, "package models") {
		t.Error("expected package declaration")
	}
	if !strings.Contains(s, `"fmt"`) {
		t.Error("expected imports")
	}
}

func TestLoadTemplatesFromFS(t *testing.T) {
	fsys := fstest.MapFS{
		"templates/model.go.tpl": &fstest.MapFile{
			Data: []byte(`// {{.Table.Name}} model`),
		},
		"templates/query.go.tpl": &fstest.MapFile{
			Data: []byte(`// query for {{.Table.Name}}`),
		},
	}
	tl, err := LoadTemplatesFromFS(fsys, "templates", DefaultTemplateFuncs())
	if err != nil {
		t.Fatal(err)
	}
	names := tl.Names()
	if len(names) < 2 {
		t.Errorf("expected at least 2 templates, got %d: %v", len(names), names)
	}
}

func TestLoadTemplatesFromDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.go.tpl"), []byte(`// {{.}}`), 0644)
	tl, err := LoadTemplatesFromDir(dir, DefaultTemplateFuncs())
	if err != nil {
		t.Fatal(err)
	}
	names := tl.Names()
	if len(names) == 0 {
		t.Error("expected at least one template")
	}
}

func TestGetOutputFilename(t *testing.T) {
	tests := []struct {
		name   string
		isTest bool
		want   string
	}{
		{"users", false, "users.go"},
		{"users", true, "users_test.go"},
		{"_special", false, "und_special.go"},
	}
	for _, tt := range tests {
		got := GetOutputFilename(tt.name, tt.isTest)
		if got != tt.want {
			t.Errorf("GetOutputFilename(%q, %v) = %q, want %q", tt.name, tt.isTest, got, tt.want)
		}
	}
}

func TestStateRun(t *testing.T) {
	cfg := &Config{
		PkgName:    "models",
		DriverName: "test",
		Aliases:    Aliases{Tables: map[string]TableAlias{}},
	}
	s := &State{
		Config: cfg,
		Tables: []schema.Table{
			{Name: "users", Columns: []schema.Column{{Name: "id"}}},
		},
		Dialect: schema.Dialect{LQ: '"', RQ: '"'},
	}
	FillAliases(&s.Config.Aliases, s.Tables)

	tmpl := template.Must(template.New("model").Parse(`// Model: {{.Table.Name}}`))
	s.Templates = &TemplateList{Template: tmpl}

	err := s.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStateRunSkipsJoinTables(t *testing.T) {
	cfg := &Config{
		PkgName: "models",
		Aliases: Aliases{Tables: map[string]TableAlias{}},
	}
	s := &State{
		Config: cfg,
		Tables: []schema.Table{
			{Name: "users"},
			{Name: "user_roles", IsJoinTable: true},
		},
	}
	FillAliases(&s.Config.Aliases, s.Tables)

	executedTables := []string{}
	tmpl := template.Must(template.New("track").Funcs(template.FuncMap{
		"track": func(name string) string {
			executedTables = append(executedTables, name)
			return ""
		},
	}).Parse(`{{track .Table.Name}}`))
	s.Templates = &TemplateList{Template: tmpl}
	s.Run()
	if len(executedTables) != 1 || executedTables[0] != "users" {
		t.Errorf("expected only 'users', got %v", executedTables)
	}
}

func TestDefaultTemplateFuncs(t *testing.T) {
	funcs := DefaultTemplateFuncs()
	if funcs["titleCase"] == nil {
		t.Error("expected titleCase in default funcs")
	}
	if funcs["camelCase"] == nil {
		t.Error("expected camelCase in default funcs")
	}
	if funcs["singular"] == nil {
		t.Error("expected singular in default funcs")
	}
	if funcs["plural"] == nil {
		t.Error("expected plural in default funcs")
	}
	if funcs["isPrimitive"] == nil {
		t.Error("expected isPrimitive in default funcs")
	}
}

func TestTemplateListNames(t *testing.T) {
	tmpl := template.New("")
	template.Must(tmpl.New("a").Parse("a"))
	template.Must(tmpl.New("b").Parse("b"))
	tl := &TemplateList{Template: tmpl}
	names := tl.Names()
	if len(names) < 2 {
		t.Errorf("expected at least 2 names, got %v", names)
	}
}

func TestTemplateListNamesNil(t *testing.T) {
	var tl *TemplateList
	names := tl.Names()
	if names != nil {
		t.Error("nil template list should return nil names")
	}
}
