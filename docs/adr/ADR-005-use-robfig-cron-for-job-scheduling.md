# ADR-005: Use Robfig Cron for Job Scheduling

## Status
Accepted

## Context
We need to implement scheduled weather notifications that run at specific intervals to send weather updates to subscribed users.

## Decision
We will use the robfig/cron/v3 library for scheduling weather notification jobs.

## Alternatives Considered
- External Job Queue: Rejected because it creates unnecessary complexity for simple time-based scheduling needs
- System Cron: Rejected because it increases deployment complexity and lacks application integration
- Ticker-based Solution: Rejected because implementing cron-like scheduling becomes too complex

## Consequences
**What becomes easier:**
- Integration becomes simple within the application process
- Timing control becomes precise with cron expression scheduling
- Testing and debugging of scheduled jobs becomes easier
- Scheduling becomes consistent across regions with built-in timezone handling

**What becomes more difficult:**
- Job persistence disappears if application crashes