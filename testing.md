# Testing Guide

This document provides instructions for running tests in the Weather Notifier application. The application includes three types of tests: Unit tests, Integration tests, and End-to-End (E2E) tests.

## Quick Start

**Run all tests with one command:**
```bash
chmod +x scripts/test-runner.sh

./scripts/test-runner.sh all
```

## Prerequisites

Make sure you have the following installed on your machine:
- **Git**
- **Docker and Docker Compose**
- **Go 1.24.3+** (for local testing)
- **Node.js 18+** (for E2E tests)

## Running Tests

### 🚀 Quick Commands

```bash
# Setup test environment (run once)
./scripts/test-runner.sh setup

# Run all tests
./scripts/test-runner.sh all

# Run specific test types
./scripts/test-runner.sh unit
./scripts/test-runner.sh integration
./scripts/test-runner.sh e2e
```

## Coverage Reports

### Generate Coverage
```bash
# Generate combined coverage report
./scripts/test-runner.sh coverage
```

## Continuous Integration

### GitHub Actions Workflows

The project includes separate workflows for each test type:

```
.github/workflows/
├── unit-tests.yml         # Unit tests only
├── integration-tests.yml  # Integration tests only
├── e2e-tests.yml          # E2E tests only
└── all-tests.yml          # All tests + coverage
```
