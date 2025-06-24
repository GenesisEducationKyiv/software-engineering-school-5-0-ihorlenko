#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${2}${1}${NC}"
}

run_with_status() {
    local description=$1
    local command=$2
    
    print_status "Running $description..." $BLUE
    
    if eval $command; then
        print_status "$description completed successfully" $GREEN
        return 0
    else
        print_status "$description failed" $RED
        return 1
    fi
}

setup_tests() {
    print_status "Setting up test dependencies..." $BLUE
    
    run_with_status "Downloading Go dependencies" "go mod download"
    
    run_with_status "Installing coverage tools" "go install github.com/wadey/gocovmerge@latest"
    
    if [ -d "tests/e2e" ]; then
        cd tests/e2e
        run_with_status "Installing E2E dependencies" "npm ci"
        run_with_status "Installing Playwright browsers" "npx playwright install"
        cd ../..
    fi
    
    print_status "Setup completed successfully" $GREEN
}

run_unit_tests() {
    run_with_status "Unit tests" "go test -v -race -coverprofile=coverage-unit.out ./internal/..."
}

run_integration_tests() {
    check_docker
    
    if [ -f "scripts/run-integration-tests.sh" ]; then
        run_with_status "Integration tests" "chmod +x scripts/run-integration-tests.sh && scripts/run-integration-tests.sh"
    else
        run_with_status "Integration tests" "go test -v -race -coverprofile=coverage-integration.out ./tests/integration/..."
    fi
}

run_e2e_tests() {
    if [ -d "tests/e2e" ]; then
        print_status "Starting application for E2E tests..." $BLUE
        make up > /dev/null 2>&1 || true
        
        local max_attempts=30
        local attempt=1
        
        while [ $attempt -le $max_attempts ]; do
            if curl -f http://localhost:8080/ping > /dev/null 2>&1; then
                print_status "Application is ready" $GREEN
                break
            fi
            
            if [ $attempt -eq $max_attempts ]; then
                print_status "Application failed to start within timeout" $RED
                make down > /dev/null 2>&1 || true
                return 1
            fi
            
            print_status "Waiting for application... (attempt $attempt/$max_attempts)" $YELLOW
            sleep 2
            ((attempt++))
        done
        
        cd tests/e2e
        local e2e_result=0
        run_with_status "E2E tests" "npm test" || e2e_result=1
        cd ../..
        
        print_status "Stopping application..." $BLUE
        make down > /dev/null 2>&1 || true
        
        return $e2e_result
    else
        print_status "E2E tests directory not found, skipping E2E tests" $YELLOW
        return 0
    fi
}

case "${1:-help}" in
    "all")
        print_status "Running all tests..." $BLUE
        
        failed_tests=""
        
        run_unit_tests || failed_tests="$failed_tests unit"
        run_integration_tests || failed_tests="$failed_tests integration"
        run_e2e_tests || failed_tests="$failed_tests e2e"
        
        if [ -n "$failed_tests" ]; then
            print_status "Some tests failed:$failed_tests" $RED
            exit 1
        else
            print_status "All tests passed!" $GREEN
            generate_coverage
        fi
        ;;
    "unit")
        run_unit_tests
        ;;
    "integration")
        run_integration_tests
        ;;
    "e2e")
        run_e2e_tests
        ;;
    *)
        print_status "Unknown command: $1" $RED
        echo ""
        exit 1
        ;;
esac