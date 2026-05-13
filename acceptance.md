# Acceptance Criteria

## Task 1: Runtime support library (boil package)

### Acceptance Criteria
- [ ] Executor interface defines Exec, Query, QueryRow methods for basic SQL execution
- [ ] ContextExecutor interface adds ExecContext, QueryContext, QueryRowContext methods
- [ ] Transactor interface adds Commit and Rollback to Executor
- [ ] Beginner interface provides Begin() method, ContextBeginner provides BeginTx()
- [ ] SetDB stores a global executor; GetDB retrieves it; if the executor also implements ContextExecutor, GetContextDB returns it
- [ ] SetLocation stores a timezone for timestamps; GetLocation retrieves it; defaults to UTC
- [ ] Column selection None() returns empty set; Infer() means auto-detect; Whitelist specifies exact columns; Blacklist excludes columns; Greylist adds to inferred
- [ ] InsertColumnSet correctly computes insert columns and return columns for each column strategy (None, Infer, Whitelist, Blacklist, Greylist)
- [ ] UpdateColumnSet correctly computes update columns for each strategy, excluding primary key columns
- [ ] HookPoint constants exist for BeforeInsert, BeforeUpdate, BeforeDelete, BeforeUpsert, AfterInsert, AfterSelect, AfterUpdate, AfterDelete, AfterUpsert
- [ ] SkipHooks returns a context that causes HooksAreSkipped to return true
- [ ] SkipTimestamps returns a context that causes TimestampsAreSkipped to return true
- [ ] DebugMode global flag defaults to false; DebugWriter defaults to os.Stdout
- [ ] WithDebug/IsDebug work per-context, falling back to global DebugMode
- [ ] WithDebugWriter/DebugWriterFrom work per-context, falling back to global DebugWriter
- [ ] WrapErr wraps an error; IsBoilErr correctly identifies wrapped errors
- [ ] Begin() panics if global DB doesn't implement Beginner; works if it does
- [ ] BeginTx() panics if global DB doesn't implement ContextBeginner; works if it does
