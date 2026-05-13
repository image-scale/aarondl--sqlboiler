package orm

import "context"

type HookPoint int

const (
	BeforeInsertHook HookPoint = iota + 1
	BeforeUpdateHook
	BeforeDeleteHook
	BeforeUpsertHook
	AfterInsertHook
	AfterSelectHook
	AfterUpdateHook
	AfterDeleteHook
	AfterUpsertHook
)

func SkipHooks(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeySkipHooks, true)
}

func HooksAreSkipped(ctx context.Context) bool {
	skip, ok := ctx.Value(ctxKeySkipHooks).(bool)
	return ok && skip
}

func SkipTimestamps(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeySkipTimestamps, true)
}

func TimestampsAreSkipped(ctx context.Context) bool {
	skip, ok := ctx.Value(ctxKeySkipTimestamps).(bool)
	return ok && skip
}
