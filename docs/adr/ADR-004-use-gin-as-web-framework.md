# ADR-004: Use Gin as Web Framework

## Status
Accepted

## Context
We need to choose a web framework for building REST API endpoints and handling HTTP middleware for the weather notification service.

## Decision
We will use Gin as the web framework for HTTP server implementation.

## Alternatives Considered
- Echo: Rejected because it has a smaller ecosystem and less community adoption
- Chi: Rejected because it functions more as a router than a full-featured framework
- Standard Library: Rejected because it lacks built-in features like JSON binding and middleware

## Consequences
**What becomes easier:**
- Development accelerates with built-in features (routing, middleware, binding)
- HTTP operations perform excellently
- Common web application needs are met with comprehensive middleware ecosystem

**What becomes more difficult:**
- Flexibility decreases for complex routing requirements