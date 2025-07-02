#!/bin/bash

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Running integration tests from project root...${NC}"

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "Project root: $PROJECT_ROOT"
cd "$PROJECT_ROOT"

if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: go.mod not found. Make sure you're running this from the project root.${NC}"
    exit 1
fi

if [ ! -d "migrations" ]; then
    echo -e "${RED}Error: migrations directory not found.${NC}"
    exit 1
fi

export MIGRATIONS_PATH="$(pwd)/migrations"
export TEST_MIGRATIONS_PATH="$(pwd)/migrations"

echo -e "${BLUE}Running integration tests...${NC}"
go test -v -race -coverprofile=coverage-integration.out ./tests/integration/...

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Integration tests completed successfully${NC}"
else
    echo -e "${RED}Integration tests failed${NC}"
    exit 1
fi