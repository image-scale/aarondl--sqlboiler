# Goal

## Project
sqlboiler — a Go project.

## Description
SQLBoiler is a database-first ORM code generator for Go. It reads database schemas and generates type-safe Go code tailored to the specific schema. The project provides:

1. A runtime support library (`boil` package) with database executor interfaces, column selection strategies, lifecycle hooks, debug logging, and global state management.
2. A query building system (`queries` package) with composable query modifiers, SQL generation for SELECT/UPDATE/DELETE, reflection-based result binding, and eager loading of relationships.
3. Database type mappings (`types` package) for exotic column types like PostgreSQL arrays, JSON, hstore, and arbitrary-precision decimals.
4. A driver abstraction layer (`drivers` package) defining interfaces for database drivers, table/column/key/relationship metadata types, relationship detection from foreign keys, and support for binary out-of-process drivers.
5. An import management system (`importers` package) for controlling which Go imports are needed in generated code.
6. A code generation engine (`boilingcore` package) that orchestrates template execution, configuration, alias resolution, and output file generation.

## Scope
- ~62 production source files to implement
- ~34 test files to write
- Reproduce all core source code, tests, and configuration
