package orm

import "sort"

const (
	columnsNone      = 0
	columnsInfer     = 1
	columnsWhitelist = 2
	columnsGreylist  = 3
	columnsBlacklist = 4
)

type Columns struct {
	Kind int
	Cols []string
}

func None() Columns {
	return Columns{Kind: columnsNone}
}

func Infer() Columns {
	return Columns{Kind: columnsInfer}
}

func Whitelist(columns ...string) Columns {
	return Columns{Kind: columnsWhitelist, Cols: columns}
}

func Blacklist(columns ...string) Columns {
	return Columns{Kind: columnsBlacklist, Cols: columns}
}

func Greylist(columns ...string) Columns {
	return Columns{Kind: columnsGreylist, Cols: columns}
}

func (c Columns) IsNone() bool      { return c.Kind == columnsNone }
func (c Columns) IsInfer() bool     { return c.Kind == columnsInfer }
func (c Columns) IsWhitelist() bool { return c.Kind == columnsWhitelist }
func (c Columns) IsBlacklist() bool { return c.Kind == columnsBlacklist }
func (c Columns) IsGreylist() bool  { return c.Kind == columnsGreylist }

func (c Columns) InsertColumnSet(allCols, defaults, noDefaults, nonZeroDefaults []string) ([]string, []string) {
	switch c.Kind {
	case columnsNone:
		return nil, nil
	case columnsInfer:
		insert := make([]string, 0, len(noDefaults)+len(nonZeroDefaults))
		insert = append(insert, noDefaults...)
		insert = append(insert, nonZeroDefaults...)
		sort.Strings(insert)
		ret := stringSliceDiff(defaults, insert)
		return insert, ret
	case columnsWhitelist:
		insert := make([]string, len(c.Cols))
		copy(insert, c.Cols)
		sort.Strings(insert)
		ret := stringSliceDiff(defaults, insert)
		return insert, ret
	case columnsBlacklist:
		insert := make([]string, 0, len(noDefaults)+len(nonZeroDefaults))
		insert = append(insert, noDefaults...)
		insert = append(insert, nonZeroDefaults...)
		insert = removeFromSlice(insert, c.Cols)
		sort.Strings(insert)
		ret := stringSliceDiff(defaults, insert)
		return insert, ret
	case columnsGreylist:
		insert := make([]string, 0, len(noDefaults)+len(nonZeroDefaults)+len(c.Cols))
		insert = append(insert, noDefaults...)
		insert = append(insert, nonZeroDefaults...)
		insert = append(insert, c.Cols...)
		insert = deduplicate(insert)
		sort.Strings(insert)
		ret := stringSliceDiff(defaults, insert)
		return insert, ret
	}
	return nil, nil
}

func (c Columns) UpdateColumnSet(allColumns, pkeyCols []string) []string {
	switch c.Kind {
	case columnsNone:
		return nil
	case columnsInfer:
		result := removeFromSlice(copySlice(allColumns), pkeyCols)
		sort.Strings(result)
		return result
	case columnsWhitelist:
		result := make([]string, len(c.Cols))
		copy(result, c.Cols)
		sort.Strings(result)
		return result
	case columnsBlacklist:
		result := removeFromSlice(copySlice(allColumns), pkeyCols)
		result = removeFromSlice(result, c.Cols)
		sort.Strings(result)
		return result
	case columnsGreylist:
		result := removeFromSlice(copySlice(allColumns), pkeyCols)
		result = append(result, c.Cols...)
		result = deduplicate(result)
		sort.Strings(result)
		return result
	}
	return nil
}

func stringSliceDiff(all, exclude []string) []string {
	excludeSet := make(map[string]struct{}, len(exclude))
	for _, e := range exclude {
		excludeSet[e] = struct{}{}
	}
	result := make([]string, 0, len(all))
	for _, s := range all {
		if _, found := excludeSet[s]; !found {
			result = append(result, s)
		}
	}
	return result
}

func removeFromSlice(slice, toRemove []string) []string {
	removeSet := make(map[string]struct{}, len(toRemove))
	for _, r := range toRemove {
		removeSet[r] = struct{}{}
	}
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if _, found := removeSet[s]; !found {
			result = append(result, s)
		}
	}
	return result
}

func copySlice(s []string) []string {
	c := make([]string, len(s))
	copy(c, s)
	return c
}

func deduplicate(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	result := make([]string, 0, len(s))
	for _, v := range s {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
