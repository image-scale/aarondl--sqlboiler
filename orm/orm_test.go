package orm

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"testing"
	"time"
)

type mockExecutor struct{}

func (m *mockExecutor) Exec(query string, args ...any) (sql.Result, error) {
	return nil, nil
}
func (m *mockExecutor) Query(query string, args ...any) (*sql.Rows, error) {
	return nil, nil
}
func (m *mockExecutor) QueryRow(query string, args ...any) *sql.Row {
	return nil
}

type mockContextExecutor struct {
	mockExecutor
}

func (m *mockContextExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return nil, nil
}
func (m *mockContextExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return nil, nil
}
func (m *mockContextExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return nil
}

func TestSetDBAndGetDB(t *testing.T) {
	exec := &mockExecutor{}
	SetDB(exec)
	if got := GetDB(); got != exec {
		t.Errorf("GetDB() = %v, want %v", got, exec)
	}
	if got := GetContextDB(); got != nil {
		t.Errorf("GetContextDB() should be nil for non-context executor, got %v", got)
	}
}

func TestSetDBWithContextExecutor(t *testing.T) {
	exec := &mockContextExecutor{}
	SetDB(exec)
	if got := GetDB(); got != exec {
		t.Errorf("GetDB() = %v, want %v", got, exec)
	}
	if got := GetContextDB(); got != exec {
		t.Errorf("GetContextDB() = %v, want %v", got, exec)
	}
}

func TestLocationDefaultsToUTC(t *testing.T) {
	loc := GetLocation()
	if loc != time.UTC {
		t.Errorf("default location = %v, want UTC", loc)
	}
}

func TestSetAndGetLocation(t *testing.T) {
	eastern, _ := time.LoadLocation("America/New_York")
	SetLocation(eastern)
	if got := GetLocation(); got != eastern {
		t.Errorf("GetLocation() = %v, want %v", got, eastern)
	}
	SetLocation(time.UTC)
}

func TestHookPoints(t *testing.T) {
	points := []HookPoint{
		BeforeInsertHook, BeforeUpdateHook, BeforeDeleteHook, BeforeUpsertHook,
		AfterInsertHook, AfterSelectHook, AfterUpdateHook, AfterDeleteHook, AfterUpsertHook,
	}
	for i, p := range points {
		if int(p) != i+1 {
			t.Errorf("HookPoint %d has value %d, want %d", i, p, i+1)
		}
	}
}

func TestSkipHooksContext(t *testing.T) {
	ctx := context.Background()
	if HooksAreSkipped(ctx) {
		t.Error("hooks should not be skipped on fresh context")
	}
	ctx = SkipHooks(ctx)
	if !HooksAreSkipped(ctx) {
		t.Error("hooks should be skipped after SkipHooks")
	}
}

func TestSkipTimestampsContext(t *testing.T) {
	ctx := context.Background()
	if TimestampsAreSkipped(ctx) {
		t.Error("timestamps should not be skipped on fresh context")
	}
	ctx = SkipTimestamps(ctx)
	if !TimestampsAreSkipped(ctx) {
		t.Error("timestamps should be skipped after SkipTimestamps")
	}
}

func TestDebugDefaults(t *testing.T) {
	if DebugMode {
		t.Error("DebugMode should default to false")
	}
	if DebugWriter == nil {
		t.Error("DebugWriter should not be nil")
	}
}

func TestWithDebugContext(t *testing.T) {
	ctx := context.Background()
	DebugMode = false
	if IsDebug(ctx) {
		t.Error("IsDebug should return false when global DebugMode is false")
	}
	ctx = WithDebug(ctx, true)
	if !IsDebug(ctx) {
		t.Error("IsDebug should return true after WithDebug(ctx, true)")
	}
	DebugMode = true
	ctx2 := context.Background()
	if !IsDebug(ctx2) {
		t.Error("IsDebug should return true when global DebugMode is true")
	}
	ctx2 = WithDebug(ctx2, false)
	if IsDebug(ctx2) {
		t.Error("IsDebug should return false when per-context override is false")
	}
	DebugMode = false
}

func TestWithDebugWriterContext(t *testing.T) {
	ctx := context.Background()
	defaultWriter := DebugWriterFrom(ctx)
	if defaultWriter == nil {
		t.Error("DebugWriterFrom should return non-nil default")
	}
	var buf bytes.Buffer
	ctx = WithDebugWriter(ctx, &buf)
	if got := DebugWriterFrom(ctx); got != io.Writer(&buf) {
		t.Error("DebugWriterFrom should return the per-context writer")
	}
}

func TestWrapErrAndIsBoilErr(t *testing.T) {
	original := errors.New("something went wrong")
	wrapped := WrapErr(original)
	if !IsBoilErr(wrapped) {
		t.Error("IsBoilErr should return true for wrapped error")
	}
	if wrapped.Error() != "something went wrong" {
		t.Errorf("wrapped error message = %q, want %q", wrapped.Error(), "something went wrong")
	}
	if IsBoilErr(original) {
		t.Error("IsBoilErr should return false for non-wrapped error")
	}
}

func TestColumnsNone(t *testing.T) {
	c := None()
	if !c.IsNone() {
		t.Error("None().IsNone() should be true")
	}
	insert, ret := c.InsertColumnSet(nil, nil, nil, nil)
	if insert != nil || ret != nil {
		t.Errorf("None InsertColumnSet should return nil, nil; got %v, %v", insert, ret)
	}
	update := c.UpdateColumnSet([]string{"a", "b"}, []string{"a"})
	if update != nil {
		t.Errorf("None UpdateColumnSet should return nil; got %v", update)
	}
}

func TestColumnsInfer(t *testing.T) {
	c := Infer()
	if !c.IsInfer() {
		t.Error("Infer().IsInfer() should be true")
	}
	allCols := []string{"id", "name", "age", "created_at"}
	defaults := []string{"id", "created_at"}
	noDefaults := []string{"name", "age"}
	nonZeroDefaults := []string{"created_at"}

	insert, ret := c.InsertColumnSet(allCols, defaults, noDefaults, nonZeroDefaults)
	if len(insert) != 3 {
		t.Errorf("Infer insert should have 3 cols, got %d: %v", len(insert), insert)
	}
	if !sliceContains(insert, "name") || !sliceContains(insert, "age") || !sliceContains(insert, "created_at") {
		t.Errorf("Infer insert should contain name, age, created_at; got %v", insert)
	}
	if len(ret) != 1 || ret[0] != "id" {
		t.Errorf("Infer return should be [id], got %v", ret)
	}
}

func TestColumnsWhitelist(t *testing.T) {
	c := Whitelist("name", "age")
	if !c.IsWhitelist() {
		t.Error("Whitelist().IsWhitelist() should be true")
	}
	defaults := []string{"id", "created_at"}
	insert, ret := c.InsertColumnSet(nil, defaults, nil, nil)
	if len(insert) != 2 || !sliceContains(insert, "name") || !sliceContains(insert, "age") {
		t.Errorf("Whitelist insert should be [age, name], got %v", insert)
	}
	if len(ret) != 2 || !sliceContains(ret, "id") || !sliceContains(ret, "created_at") {
		t.Errorf("Whitelist return should be [id, created_at], got %v", ret)
	}
}

func TestColumnsBlacklist(t *testing.T) {
	c := Blacklist("created_at")
	if !c.IsBlacklist() {
		t.Error("Blacklist().IsBlacklist() should be true")
	}
	defaults := []string{"id", "created_at"}
	noDefaults := []string{"name", "age"}
	nonZeroDefaults := []string{"created_at"}

	insert, ret := c.InsertColumnSet(nil, defaults, noDefaults, nonZeroDefaults)
	if sliceContains(insert, "created_at") {
		t.Errorf("Blacklist should exclude created_at from insert, got %v", insert)
	}
	if !sliceContains(insert, "name") || !sliceContains(insert, "age") {
		t.Errorf("Blacklist should include name and age, got %v", insert)
	}
	if !sliceContains(ret, "created_at") {
		t.Errorf("Blacklist return should include created_at since it was excluded from insert, got %v", ret)
	}
}

func TestColumnsGreylist(t *testing.T) {
	c := Greylist("created_at")
	if !c.IsGreylist() {
		t.Error("Greylist().IsGreylist() should be true")
	}
	defaults := []string{"id", "created_at"}
	noDefaults := []string{"name", "age"}
	nonZeroDefaults := []string{}

	insert, ret := c.InsertColumnSet(nil, defaults, noDefaults, nonZeroDefaults)
	if !sliceContains(insert, "created_at") {
		t.Errorf("Greylist should add created_at to insert, got %v", insert)
	}
	if !sliceContains(insert, "name") || !sliceContains(insert, "age") {
		t.Errorf("Greylist should include name and age, got %v", insert)
	}
	if sliceContains(ret, "created_at") {
		t.Errorf("Greylist return should NOT include created_at since it was inserted, got %v", ret)
	}
}

func TestUpdateColumnSetInfer(t *testing.T) {
	c := Infer()
	allCols := []string{"id", "name", "age", "created_at"}
	pkeys := []string{"id"}
	result := c.UpdateColumnSet(allCols, pkeys)
	if sliceContains(result, "id") {
		t.Error("Infer update should exclude primary keys")
	}
	if !sliceContains(result, "name") || !sliceContains(result, "age") || !sliceContains(result, "created_at") {
		t.Errorf("Infer update should include non-PK columns, got %v", result)
	}
}

func TestUpdateColumnSetWhitelist(t *testing.T) {
	c := Whitelist("name")
	result := c.UpdateColumnSet([]string{"id", "name", "age"}, []string{"id"})
	if len(result) != 1 || result[0] != "name" {
		t.Errorf("Whitelist update should be exactly [name], got %v", result)
	}
}

func TestUpdateColumnSetBlacklist(t *testing.T) {
	c := Blacklist("age")
	result := c.UpdateColumnSet([]string{"id", "name", "age", "email"}, []string{"id"})
	if sliceContains(result, "id") || sliceContains(result, "age") {
		t.Errorf("Blacklist update should exclude id and age, got %v", result)
	}
	if !sliceContains(result, "name") || !sliceContains(result, "email") {
		t.Errorf("Blacklist update should include name and email, got %v", result)
	}
}

func TestBeginPanicsWithNonBeginner(t *testing.T) {
	SetDB(&mockExecutor{})
	defer func() {
		if r := recover(); r == nil {
			t.Error("Begin should panic when DB doesn't implement Beginner")
		}
	}()
	Begin()
}

func TestBeginTxPanicsWithNonContextBeginner(t *testing.T) {
	SetDB(&mockExecutor{})
	defer func() {
		if r := recover(); r == nil {
			t.Error("BeginTx should panic when DB doesn't implement ContextBeginner")
		}
	}()
	BeginTx(context.Background(), nil)
}

func sliceContains(s []string, val string) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}
	return false
}
