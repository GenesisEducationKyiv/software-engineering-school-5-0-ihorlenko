# ADR-002: Use PostgreSQL as Primary Database

## Status
Accepted

## Context
We need to choose a database system for storing user data, subscription information, and managing relational data for the weather notification service.

## Decision
We will use PostgreSQL as the primary database for the application.

## Alternatives Considered
- MySQL: Rejected because it has less robust JSON support and stricter SQL mode requirements
- SQLite: Rejected because it lacks concurrent write support needed for multi-user subscriptions
- MongoDB: Rejected because a relational data model better fits user-subscription relationships

## Consequences
**What becomes easier:**
- Data consistency and integrity improve with ACID compliance
- Production deployment benefits from mature backup and replication solutions
- Query optimization improves with rich indexing capabilities
- Future extension possibilities increase with native JSON support for weather data storage

**What becomes more difficult:**
- Memory usage increases compared to lighter databases
- Performance optimization requires more configuration
- Operational overhead increases for simple use cases