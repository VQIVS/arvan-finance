#!/bin/bash

# Simple test runner script for SMS service tests

set -e

echo "SMS Service Test Runner"
echo "======================="

# Change to project directory (go up one level from tests folder)
cd "$(dirname "$0")/.."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# Function to run basic tests
run_basic_tests() {
    echo "Running basic unit tests..."
    if go test ./tests -v; then
        echo "Basic tests passed!"
    else
        echo "Basic tests failed!"
        exit 1
    fi
}

# Function to run tests with coverage
run_coverage_tests() {
    echo "Running tests with coverage..."
    if go test ./tests -cover -coverprofile=coverage.out; then
        echo "Coverage tests completed!"
        
        # Generate HTML coverage report
        echo "Generating HTML coverage report..."
        go tool cover -html=coverage.out -o coverage.html
        echo "Coverage report generated: coverage.html"
    else
        echo "Coverage tests failed!"
        exit 1
    fi
}

# Function to run benchmark tests
run_benchmarks() {
    echo "Running benchmark tests..."
    if go test -bench=. -benchmem ./tests; then
        echo "Benchmark tests completed!"
    else
        echo "Benchmark tests failed!"
        exit 1
    fi
}

# Function to run load tests
run_load_tests() {
    echo "Running load tests..."
    if go test -bench=BenchmarkHighVolume -benchtime=10s -benchmem ./tests; then
        echo "Load tests completed!"
    else
        echo "Load tests failed!"
        exit 1
    fi
}

# Function to run stress tests
run_stress_tests() {
    echo "Running stress tests..."
    if go test -bench=BenchmarkStress -benchtime=30s -benchmem ./tests; then
        echo "Stress tests completed!"
    else
        echo "Stress tests failed!"
        exit 1
    fi
}

# Function to run memory benchmarks
run_memory_tests() {
    echo "Running memory efficiency tests..."
    if go test -bench=BenchmarkMemory -benchmem ./tests; then
        echo "Memory tests completed!"
    else
        echo "Memory tests failed!"
        exit 1
    fi
}

# Function to run concurrent benchmarks
run_concurrent_tests() {
    echo "Running concurrency benchmarks..."
    if go test -bench=BenchmarkConcurrent -cpu=1,2,4,8 -benchmem ./tests; then
        echo "Concurrency tests completed!"
    else
        echo "Concurrency tests failed!"
        exit 1
    fi
}

# Function to run performance profiling
run_profiling() {
    echo "Running performance profiling..."
    echo "CPU Profiling..."
    go test -bench=BenchmarkConcurrentOperations -cpuprofile=cpu.prof ./tests
    echo "Memory Profiling..."
    go test -bench=BenchmarkMemoryUsage -memprofile=mem.prof ./tests
    echo "Profiles generated: cpu.prof, mem.prof"
    echo "To analyze profiles:"
    echo "  go tool pprof cpu.prof"
    echo "  go tool pprof mem.prof"
}

# Function to run specific test
run_specific_test() {
    local test_name="$1"
    echo "Running specific test: $test_name"
    if go test -run "$test_name" ./tests -v; then
        echo "Test '$test_name' passed!"
    else
        echo "Test '$test_name' failed!"
        exit 1
    fi
}

# Function to run race condition tests
run_race_tests() {
    echo "Running tests with race detection..."
    if go test -race ./tests; then
        echo "Race condition tests passed!"
    else
        echo "Race condition detected!"
        exit 1
    fi
}

# Function to clean test artifacts
clean_artifacts() {
    echo "Cleaning test artifacts..."
    rm -f coverage.out coverage.html cpu.prof mem.prof
    echo "Artifacts cleaned!"
}

# Main menu
show_help() {
    echo "Usage: $0 [OPTION]"
    echo 
    echo "Options:"
    echo "  basic      Run basic unit tests"
    echo "  coverage   Run tests with coverage report"
    echo "  bench      Run benchmark tests"
    echo "  load       Run load tests (high volume)"
    echo "  stress     Run stress tests (extreme load)"
    echo "  memory     Run memory efficiency tests"
    echo "  concurrent Run concurrency benchmarks"
    echo "  profile    Run performance profiling"
    echo "  race       Run tests with race detection"
    echo "  all        Run all test types"
    echo "  performance Run all performance tests (bench, load, stress, memory, concurrent)"
    echo "  clean      Clean test artifacts"
    echo "  specific   Run specific test (example: $0 specific TestCreateUser)"
    echo "  help       Show this help message"
}

# Parse command line arguments
case "${1:-help}" in
    "basic")
        run_basic_tests
        ;;
    "coverage")
        run_coverage_tests
        ;;
    "bench")
        run_benchmarks
        ;;
    "load")
        run_load_tests
        ;;
    "stress")
        run_stress_tests
        ;;
    "memory")
        run_memory_tests
        ;;
    "concurrent")
        run_concurrent_tests
        ;;
    "profile")
        run_profiling
        ;;
    "race")
        run_race_tests
        ;;
    "all")
        echo "Running complete test suite..."
        run_basic_tests
        echo
        run_coverage_tests
        echo
        run_benchmarks
        echo
        run_race_tests
        echo
        echo "All tests completed successfully!"
        ;;
    "performance")
        echo "Running performance test suite..."
        run_benchmarks
        echo
        run_load_tests
        echo
        run_stress_tests
        echo
        run_memory_tests
        echo
        run_concurrent_tests
        echo
        echo "Performance tests completed successfully!"
        ;;
    "clean")
        clean_artifacts
        ;;
    "specific")
        if [ -z "$2" ]; then
            echo "Error: Please provide test name as second argument"
            echo "Example: $0 specific TestCreateUser"
            exit 1
        fi
        run_specific_test "$2"
        ;;
    "help"|*)
        show_help
        ;;
esac
