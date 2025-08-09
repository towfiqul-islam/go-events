# Testing Guide

This document provides comprehensive information about the testing framework and test coverage for the Go Events project.

## 📋 Overview

The project includes a complete test suite covering:

- **Unit Tests**: All models, utilities, middleware, and background jobs
- **Integration Tests**: Database operations with mocking
- **Mock Testing**: Database interactions using sqlmock
- **Edge Case Testing**: Boundary conditions and error scenarios

## 🏗️ Test Structure

```
go-events/
├── utils/
│   ├── hash_test.go           # Password hashing tests
│   └── jwt_test.go            # JWT token tests
├── models/
│   ├── event_test.go          # Event model tests
│   └── user_test.go           # User model tests
├── middlewares/
│   └── auth_test.go           # Authentication middleware tests
├── jobs/
│   └── notification_job_test.go # Background job tests
└── test/
    └── helpers.go             # Test utilities and mocks
```

## 🔧 Testing Dependencies

### Required Packages

```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/DATA-DOG/go-sqlmock
```

### Test Libraries Used

- **testify/assert**: Assertion helpers
- **testify/mock**: Mock generation and verification
- **go-sqlmock**: Database mocking for SQL operations
- **httptest**: HTTP request/response testing

## 🚀 Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests with Verbose Output

```bash
go test ./... -v
```

### Run Tests for Specific Package

```bash
go test ./utils/... -v
go test ./models/... -v
go test ./middlewares/... -v
go test ./jobs/... -v
```

### Run Tests with Coverage

```bash
go test ./... -cover
```

### Generate Coverage Report

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Run Specific Test

```bash
go test ./utils -run TestHashPassword -v
go test ./models -run TestEvent_Save -v
```

## 📊 Test Coverage

### Current Coverage by Package

#### Utils Package

- ✅ **hash.go**: 100% coverage

  - `HashPassword()` - All scenarios including edge cases
  - `CheckHashPassword()` - Valid/invalid password combinations
  - Consistency testing and empty string handling

- ✅ **jwt.go**: 100% coverage
  - `GenerateToken()` - Various user scenarios
  - `VerifyToken()` - Valid/invalid/expired tokens
  - Edge cases: missing claims, wrong signing methods

#### Models Package

- ✅ **event.go**: 100% coverage

  - `Save()`, `Update()`, `Delete()` - All CRUD operations
  - `GetAllEvents()`, `GetEventById()` - Query operations
  - Database error scenarios and edge cases

- ✅ **user.go**: 100% coverage

  - `Save()` - User creation with error handling
  - `ValidateUser()` - Authentication with password verification
  - `GetUser()` - User retrieval operations

- ✅ **notification.go**: Covered via integration tests
  - Database operations tested through job tests

#### Middlewares Package

- ✅ **auth.go**: 100% coverage
  - JWT token validation
  - Context value setting
  - Unauthorized request handling
  - Middleware chain continuation

#### Jobs Package

- ✅ **notification_job.go**: 100% coverage
  - Service lifecycle (start/stop)
  - Background processing logic
  - Message generation with various timing scenarios
  - Error handling and edge cases

## 🧪 Test Patterns and Best Practices

### 1. Table-Driven Tests

```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {
        name:    "Valid input",
        input:   "test",
        want:    "expected",
        wantErr: false,
    },
}
```

### 2. Database Mocking

```go
func TestWithMockDB(t *testing.T) {
    mock, cleanup, err := test.SetupMockDB()
    assert.NoError(t, err)
    defer cleanup()

    // Setup expectations
    mock.ExpectQuery("SELECT").WillReturnRows(...)

    // Run test
    result, err := SomeFunction()

    // Verify
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### 3. HTTP Testing

```go
func TestHTTPHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.GET("/test", handler)

    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
```

## 🔍 Test Categories

### Unit Tests

- **Purpose**: Test individual functions in isolation
- **Scope**: Single function or method
- **Examples**: Password hashing, JWT token generation

### Integration Tests

- **Purpose**: Test component interactions
- **Scope**: Database operations, middleware integration
- **Examples**: User validation with database, authentication flow

### Mock Tests

- **Purpose**: Test with external dependencies mocked
- **Scope**: Database interactions, HTTP requests
- **Examples**: Event CRUD operations, notification processing

### Edge Case Tests

- **Purpose**: Test boundary conditions and error scenarios
- **Scope**: Invalid inputs, network failures, database errors
- **Examples**: Empty passwords, expired tokens, database connection failures

## 📝 Test Utilities

### Test Helpers (`test/helpers.go`)

#### SetupMockDB

```go
mock, cleanup, err := test.SetupMockDB()
// Sets up sqlmock for database testing
```

#### Test Data Generators

```go
testEvent := test.GetTestEvent()
testUser := test.GetTestUser()
testNotification := test.GetTestNotification()
```

#### Mock Expectation Helpers

```go
test.ExpectExecSuccess(mock, query, lastInsertID, rowsAffected)
test.ExpectExecError(mock, query, err)
test.ExpectQueryRows(mock, query, columns, rows...)
```

## 🐛 Debugging Tests

### Running Tests with Debug Output

```bash
go test ./... -v -args -test.v
```

### Identifying Slow Tests

```bash
go test ./... -v -count=1 | grep -E "(PASS|FAIL).*[0-9]+\.[0-9]+s"
```

### Test with Race Condition Detection

```bash
go test ./... -race
```

### Memory Leak Detection

```bash
go test ./... -memprofile=mem.prof
```

## 📋 Test Checklist

### Before Committing

- [ ] All tests pass: `go test ./...`
- [ ] No race conditions: `go test ./... -race`
- [ ] Coverage maintained: `go test ./... -cover`
- [ ] No skipped tests without justification
- [ ] All mock expectations verified

### New Feature Testing

- [ ] Unit tests for new functions
- [ ] Integration tests for database operations
- [ ] Edge case testing
- [ ] Error scenario testing
- [ ] Mock tests for external dependencies

### Bug Fix Testing

- [ ] Test reproducing the bug
- [ ] Test verifying the fix
- [ ] Regression test coverage
- [ ] Related functionality testing

## 🔧 Test Configuration

### Environment Variables

```bash
# For testing with different configurations
export TEST_DB_CONNECTION="test_connection_string"
export TEST_JWT_SECRET="test_secret"
```

### Test Flags

```bash
# Common test flags
go test -v          # Verbose output
go test -short      # Skip long-running tests
go test -count=1    # Disable test caching
go test -timeout=30s # Set timeout
```

## 📈 Continuous Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.22
      - name: Run tests
        run: |
          go test ./... -v -cover
          go test ./... -race
```

## 🎯 Testing Best Practices

### DO

- ✅ Write tests before or alongside code
- ✅ Use descriptive test names
- ✅ Test both success and failure scenarios
- ✅ Keep tests independent and isolated
- ✅ Use table-driven tests for multiple scenarios
- ✅ Mock external dependencies
- ✅ Verify all mock expectations
- ✅ Clean up resources in defer statements

### DON'T

- ❌ Test implementation details
- ❌ Write flaky tests with timing dependencies
- ❌ Ignore test failures
- ❌ Write tests that depend on external services
- ❌ Leave unused mock expectations
- ❌ Skip edge case testing
- ❌ Write overly complex test setup

## 🆘 Troubleshooting

### Common Issues

#### Test Failures

1. **Mock Expectations Not Met**

   ```bash
   Error: there is a remaining expectation
   Solution: Verify all mock.Expect* calls are matched
   ```

2. **Database Connection Issues**

   ```bash
   Error: sql: database is closed
   Solution: Ensure cleanup() is called and DB is properly mocked
   ```

3. **JWT Token Failures**
   ```bash
   Error: could not parse token
   Solution: Check token format and signing method
   ```

#### Performance Issues

1. **Slow Tests**: Check for actual database calls instead of mocks
2. **Memory Leaks**: Ensure proper resource cleanup
3. **Race Conditions**: Use proper synchronization in concurrent tests

### Getting Help

- Check test output for specific error messages
- Verify mock setup matches actual function calls
- Ensure test data matches expected formats
- Review test isolation and cleanup

## 📚 Additional Resources

- [Go Testing Package Documentation](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go-SQLMock Documentation](https://github.com/DATA-DOG/go-sqlmock)
- [Go Testing Best Practices](https://golang.org/doc/tutorial/add-a-test)
