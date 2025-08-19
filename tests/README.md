# SMS Billing Service Tests

Unit tests and benchmarks for the SMS billing service.

## Quick Start

```bash
# Run all tests
./tests/run_tests.sh basic

# Run with coverage
./tests/run_tests.sh coverage

# Run benchmarks
./tests/run_tests.sh bench

# Run everything
./tests/run_tests.sh all
```

## Files

- `user_service_test.go` - Unit tests (69.2% coverage)
- `benchmark_test.go` - Performance benchmarks
- `mocks.go` - Mock implementations
- `run_tests.sh` - Test runner script

## Common Commands

```bash
./tests/run_tests.sh basic      # Unit tests
./tests/run_tests.sh coverage   # Coverage report
./tests/run_tests.sh bench      # Benchmarks
./tests/run_tests.sh load       # Load testing
./tests/run_tests.sh stress     # Stress testing
./tests/run_tests.sh race       # Race detection
./tests/run_tests.sh all        # Everything
./tests/run_tests.sh help       # Show all options
```

## Test Coverage

- ✅ User creation and retrieval
- ✅ Balance credit/debit operations
- ✅ SMS status updates
- ✅ Error handling and validation
- ✅ Concurrent operations
- ✅ Performance benchmarks

## Requirements

- Go 1.19+
- No external dependencies (uses mocks)
