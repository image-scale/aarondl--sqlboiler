# Acceptance Criteria

## Task 1: Runtime support library — DONE
## Task 2: Query building system — DONE

## Task 3: Reflection-based result binding

### Acceptance Criteria
- [ ] MakeStructMapping creates a map from column names to encoded field paths using struct tags (tag "orm" with field name, or snake_case of Go field name)
- [ ] Tag "-" excludes a field from mapping; tag ",bind" enables recursive descent into nested structs
- [ ] BindMapping produces an ordered slice of field path encodings given column names and a struct mapping
- [ ] PtrsFromMapping extracts pointer-to-field interfaces suitable for rows.Scan given a struct value and mapping
- [ ] ValuesFromMapping extracts field values given a struct value and mapping
- [ ] Bind scans SQL rows into a single struct pointer (*Struct) returning sql.ErrNoRows if empty
- [ ] Bind scans SQL rows into a slice of struct pointers (*[]*Struct) appending all matching rows
- [ ] Bind scans SQL rows into a slice of structs (*[]Struct) appending all matching rows
- [ ] Binding returns an error for invalid target types (non-pointer, non-struct, etc.)
- [ ] TitleToSnake converts Go TitleCase field names to snake_case (e.g., UserID -> user_id, HTTPSEnabled -> https_enabled)
- [ ] Equal compares two values considering nil, []byte, time.Time, and numeric types
- [ ] Assign sets a value on a destination using sql.Scanner/driver.Valuer when available
- [ ] IsNil checks if a value is nil including through Valuer interface
