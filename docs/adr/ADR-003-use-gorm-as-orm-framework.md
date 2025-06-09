# ADR-003: Use GORM as ORM Framework

## Status
Accepted

## Context
We need to choose a database access layer for Go that provides type-safe database operations, relationship management, and migration capabilities for our PostgreSQL database.

## Decision
We will use GORM as the ORM framework for database operations.

## Alternatives Considered
- Raw SQL with database/sql: Rejected because it requires more boilerplate code and lacks type safety

## Consequences
**What becomes easier:**
- Development time reduces with automatic struct-to-table mapping
- Relationship queries become simpler with preloading capabilities
- Error handling becomes consistent across database operations
- Database migrations and foreign key relationships get built-in support

**What becomes more difficult:**
- Learning curve increases for GORM-specific patterns and conventions
- Performance may decrease compared to raw SQL for complex queries