# ADR-001: Use Go as Primary Programming Language

## Status
Accepted

## Context
We need to choose a programming language for building a weather notification API service. The service must handle HTTP requests, database operations, scheduled jobs, and email notifications.

## Decision
We will use Go 1.24.3 as the primary programming language for the backend service.

## Alternatives Considered
- Node.js: Rejected because its single-threaded nature limits concurrent job processing
- PHP: Rejected because it lacks native concurrency support and has poor performance for background job processing

## Consequences
**What becomes easier:**
- Fast compilation and execution with single binary deployment
- Simple deployment process for containerization
- Handling concurrent requests and background jobs with goroutines
- Code reliability improves due to strong typing that prevents runtime errors

**What becomes more difficult:**
- Learning curve increases for developers who are not familiar with Go
- Error handling requires more verbose code patterns
- Limited generics support may require code duplication in some situations