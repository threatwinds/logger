# CLAUDE.md - Logger Library

## Library Overview

Logger is a **structured logging library** with configurable output formats, severity levels, and retry mechanisms. It provides both direct logging and error handling capabilities with UUID tracking.

**Status**: Legacy library - **prefer `go-sdk/catcher` for new development**. Logger remains in use for backward compatibility in existing services.

**Key Features:**
- Structured logging (JSON, CSV, text formats)
- Configurable severity levels (DEBUG to ALERT)
- File rotation with lumberjack
- Retry mechanisms for operations
- Gin framework integration
- UUID-based error tracking

**Migration Note**: New services should use `github.com/threatwinds/go-sdk/catcher` which provides enhanced error handling and logging. Logger is maintained for existing code.

## Installation

```go
import "github.com/threatwinds/logger"
```

## Core Features

### 1. Structured Logging

Log messages with automatic severity classification, file paths, line numbers, and timestamps.

**Severity Levels:**
- **100-199**: DEBUG
- **200-299**: INFO
- **300-399**: NOTICE
- **400-499**: WARNING
- **500-501**: ERROR
- **502-508**: CRITICAL
- **509-510**: ALERT

### 2. Error Handling

Create structured errors with UUIDs, status codes, and messages. Supports error checking and Gin integration.

### 3. Retry Operations

Built-in retry functions with configurable delays and exception handling (similar to catcher).

### 4. Multiple Output Formats

- **JSON**: Structured logging for production
- **CSV**: Tabular format for analysis
- **Text**: Human-readable format for debugging

### 5. Log Rotation

Automatic log file rotation using lumberjack (5MB max size, 100 backups, 30-day retention).

## Usage Examples

### Basic Logger Initialization

```go
package main

import (
    "github.com/threatwinds/logger"
    "time"
)

func main() {
    // Create logger with default configuration
    log := logger.NewLogger(nil)

    // Log informational message
    log.Info("service started on port %d", 8080)

    // Create custom logger
    customLog := logger.NewLogger(&logger.Config{
        Format:  "json",        // json, text, csv
        Level:   400,           // Only log WARNING and above
        Output:  "stdout",      // stdout or file path
        Retries: 5,
        Wait:    1 * time.Second,
    })

    customLog.LogF(500, "database connection failed: %s", "timeout")
}
```

### Logging with Different Severity Levels

```go
log := logger.NewLogger(nil)

// Informational logs (200)
log.Info("user %s logged in", "john@example.com")

// Custom severity logging
log.LogF(100, "Debug: cache hit for key %s", "user:123")
log.LogF(300, "Notice: API rate limit approaching")
log.LogF(400, "Warning: slow query detected: %dms", 2500)
log.LogF(500, "Error: failed to process request")
log.LogF(502, "Critical: database connection lost")
log.LogF(509, "Alert: system disk space < 5%%")
```

### Error Handling

```go
import (
    "github.com/threatwinds/logger"
    "github.com/gin-gonic/gin"
)

func processRequest(log *logger.Logger) error {
    // Create error with status code mapping
    err := log.ErrorF("user %s not found", "john@example.com")
    if err != nil {
        // err.UUID contains unique error identifier
        // err.Status contains HTTP status code
        // err.Message contains formatted message
        return err
    }

    return nil
}

// Use with Gin
func handleRequest(c *gin.Context, log *logger.Logger) {
    err := processRequest(log)
    if err != nil {
        if logErr := logger.ToError(err); logErr != nil {
            logErr.GinError(c) // Sets headers and status
            return
        }
    }

    c.JSON(200, gin.H{"status": "ok"})
}
```

### Status Code Mapping

```go
// Map error messages to specific HTTP status codes
log := logger.NewLogger(&logger.Config{
    Format: "json",
    Level:  400,
    Output: "stdout",
    StatusMap: map[int][]string{
        404: {"not found", "does not exist"},
        403: {"forbidden", "unauthorized"},
        400: {"invalid", "malformed"},
        500: {"internal error", "failed"},
    },
})

// This will map to 404
err := log.ErrorF("user not found")
// err.Status == 404

// This will map to 403
err = log.ErrorF("forbidden access to resource")
// err.Status == 403

// No match defaults to 500
err = log.ErrorF("unexpected error")
// err.Status == 500
```

### File Logging with Rotation

```go
// Log to file with automatic rotation
log := logger.NewLogger(&logger.Config{
    Format: "json",
    Level:  200, // Log INFO and above
    Output: "/var/log/myapp/app.log", // File path instead of stdout
})

// Logs will be rotated automatically:
// - Max size: 5MB per file
// - Max backups: 100 files
// - Max age: 30 days
// - Files: app.log, app.log.1, app.log.2, etc.
```

### Retry Operations

```go
import (
    "github.com/threatwinds/logger"
    "time"
)

log := logger.NewLogger(nil)

// Retry with configuration
err := log.Retry(func() error {
    return connectToDatabase()
}, "connection refused")

if err != nil {
    log.LogF(500, "failed to connect after retries: %v", err)
}

// Infinite retry until success or exception
log.InfiniteRetry(func() error {
    return connectToRedis()
}, "auth failed")

// Infinite loop for background tasks
log.InfiniteLoop(func() error {
    return processQueue()
}, "shutdown")

// Retry only on specific error
log.InfiniteRetryIfXError(func() error {
    return healthCheck()
}, "timeout")
```

### Error Exception Checking

```go
err := doSomething()

// Check if error contains specific messages
if logger.Is(err, "not found", "does not exist") {
    // Handle not found error
}

// For logger.Error types
if logErr := logger.ToError(err); logErr != nil {
    if logErr.IsAny("timeout", "connection refused") {
        // Retry logic
    }
}
```

### Complete Service Example

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/threatwinds/logger"
)

var log *logger.Logger

func main() {
    // Initialize logger
    log = logger.NewLogger(&logger.Config{
        Format: "json",
        Level:  200,
        Output: "/var/log/myservice/app.log",
        StatusMap: map[int][]string{
            404: {"not found"},
            400: {"invalid", "malformed"},
        },
    })

    log.Info("service starting on port 8080")

    // Setup routes
    r := gin.Default()
    r.GET("/users/:id", getUser)

    // Fatal logs and exits on error
    if err := r.Run(":8080"); err != nil {
        log.Fatal("failed to start server: %v", err)
    }
}

func getUser(c *gin.Context) {
    userID := c.Param("id")

    user, err := fetchUser(userID)
    if err != nil {
        if logErr := logger.ToError(err); logErr != nil {
            logErr.GinError(c)
            return
        }

        // Handle non-logger errors
        log.LogF(500, "unexpected error: %v", err)
        c.AbortWithStatus(500)
        return
    }

    log.Info("user %s fetched successfully", userID)
    c.JSON(200, user)
}

func fetchUser(userID string) (interface{}, error) {
    user := db.GetUser(userID)
    if user == nil {
        return nil, log.ErrorF("user %s not found", userID)
    }

    return user, nil
}
```

## Configuration

### Config Structure

```go
type Config struct {
    Format    string           // "json", "text", "csv"
    Level     int              // Minimum severity to log (100-510)
    Output    string           // "stdout" or file path
    Retries   int              // Number of retries (default: 5)
    Wait      time.Duration    // Wait between retries (default: 1s)
    StatusMap map[int][]string // Message → status code mapping
}
```

### Default Configuration

```go
&logger.Config{
    Format:    "json",
    Level:     400,  // WARNING and above
    Output:    "stdout",
    Retries:   5,
    Wait:      1 * time.Second,
    StatusMap: map[int][]string{},
}
```

### Format Examples

**JSON Format:**
```json
{
  "timestamp": "2025-01-14T10:00:00.000Z",
  "severity": "ERROR",
  "path": "/app/service.go",
  "line": 123,
  "uuid": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "status": 500,
  "message": "database connection failed"
}
```

**Text Format:**
```
2025-01-14T10:00:00.000Z ERROR f47ac10b-58cc-4372-a567-0e02b2c3d479 /app/service.go 123 database connection failed
```

**CSV Format:**
```csv
2025-01-14T10:00:00.000Z,ERROR,f47ac10b-58cc-4372-a567-0e02b2c3d479,/app/service.go,123,database connection failed
```

## Common Patterns

### Service Initialization with Logging

```go
func initService() error {
    log := logger.NewLogger(&logger.Config{
        Format: "json",
        Level:  200,
        Output: "/var/log/service/app.log",
    })

    log.Info("initializing service version %s", version)

    // Retry critical connections
    err := log.InfiniteRetry(func() error {
        return connectDatabase()
    }, "auth failed")

    if err != nil {
        return log.ErrorF("initialization failed: %v", err)
    }

    log.Info("service initialized successfully")
    return nil
}
```

### Background Task with Logging

```go
func backgroundWorker(log *logger.Logger) {
    log.InfiniteLoop(func() error {
        msg, err := queue.Pop()
        if err != nil {
            log.LogF(500, "queue error: %v", err)
            return nil // Continue loop despite error
        }

        if msg != nil {
            err = processMessage(msg)
            if err != nil {
                log.LogF(500, "processing failed for message %s: %v", msg.ID, err)
            } else {
                log.Info("processed message %s", msg.ID)
            }
        }

        return nil
    }, "shutdown")
}
```

## Testing

```go
import (
    "github.com/threatwinds/logger"
    "testing"
)

func TestErrorCreation(t *testing.T) {
    log := logger.NewLogger(nil)

    err := log.ErrorF("user %s not found", "123")

    if err == nil {
        t.Fatal("expected error")
    }

    if err.Status != 500 {
        t.Errorf("expected status 500, got %d", err.Status)
    }

    if err.UUID == "" {
        t.Error("expected UUID to be set")
    }
}

func TestRetryOperation(t *testing.T) {
    log := logger.NewLogger(&logger.Config{
        Retries: 3,
        Wait:    10 * time.Millisecond,
    })

    attempts := 0
    err := log.Retry(func() error {
        attempts++
        if attempts < 2 {
            return errors.New("temporary error")
        }
        return nil
    }, "permanent error")

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    if attempts != 2 {
        t.Errorf("expected 2 attempts, got %d", attempts)
    }
}
```

## Common Gotchas

### 1. Logger vs Catcher

```go
// ❌ OLD WAY (logger) - Still works but deprecated
import "github.com/threatwinds/logger"
log := logger.NewLogger(nil)
err := log.ErrorF("operation failed")

// ✅ NEW WAY (catcher) - Preferred for new code
import "github.com/threatwinds/go-sdk/catcher"
err := catcher.Error("operation failed", nil, map[string]any{"status": 500})
```

### 2. Log Level Filtering

```go
log := logger.NewLogger(&logger.Config{
    Level: 400, // Only WARNING (400) and above
})

log.Info("this will not be logged")        // Level 200 < 400
log.LogF(300, "this will not be logged")   // Level 300 < 400
log.LogF(400, "this WILL be logged")       // Level 400 = 400
log.LogF(500, "this WILL be logged")       // Level 500 > 400
```

### 3. Fatal Exits Immediately

```go
// Fatal logs and calls os.Exit(1) - cannot recover
log.Fatal("critical error") // Program terminates here
// Code below never executes
```

### 4. StatusMap Only Affects ErrorF

```go
log := logger.NewLogger(&logger.Config{
    StatusMap: map[int][]string{
        404: {"not found"},
    },
})

// StatusMap applies to ErrorF
err := log.ErrorF("user not found") // err.Status == 404

// StatusMap does NOT apply to LogF
log.LogF(500, "user not found") // Status is explicitly 500
```

### 5. File Path Must Exist

```go
// ❌ WRONG - Directory must exist
log := logger.NewLogger(&logger.Config{
    Output: "/nonexistent/dir/app.log", // Will fail
})

// ✅ CORRECT - Ensure directory exists first
os.MkdirAll("/var/log/myapp", 0755)
log := logger.NewLogger(&logger.Config{
    Output: "/var/log/myapp/app.log",
})
```

### 6. Retry Functions Block

```go
// InfiniteRetry blocks until success or exception
log.InfiniteRetry(func() error {
    return connect()
}, "auth failed")

// Run in goroutine for non-blocking
go log.InfiniteRetry(func() error {
    return connect()
}, "auth failed")
```

## Migration to Catcher

For new development, use `go-sdk/catcher` instead of `logger`:

```go
// OLD (logger)
import "github.com/threatwinds/logger"
log := logger.NewLogger(nil)
err := log.ErrorF("operation failed")
log.Info("operation succeeded")

// NEW (catcher)
import "github.com/threatwinds/go-sdk/catcher"
err := catcher.Error("operation failed", nil, map[string]any{"status": 500})
catcher.Info("operation succeeded", map[string]any{"operation": "test"})
```

**Key Differences:**
- Catcher has enhanced error metadata and stack traces
- Catcher's `Info()` is a function, not a method
- Catcher uses `SdkError` instead of `Error`
- Catcher has more retry strategies (e.g., `RetryWithBackoff`)

## Dependencies

### Required:
- `github.com/gin-gonic/gin` - Gin framework integration
- `github.com/google/uuid` - UUID generation
- `gopkg.in/natefinch/lumberjack.v2` - Log rotation

### Import Path:
```go
import "github.com/threatwinds/logger"
```

## Reference

### Related Documentation:
- Main README: `/logger/README.md` (minimal)
- Changelog: `/logger/CHANGELOG.md`
- Successor: `go-sdk/catcher` - Enhanced error handling system

### Used By:
- Legacy services (backward compatibility)
- Services not yet migrated to catcher

### Key Differences from Catcher:

| Feature | Logger | Catcher |
|---------|--------|---------|
| Error Type | `logger.Error` | `catcher.SdkError` |
| Logging | Method-based | Function-based |
| Stack Traces | File + Line only | Full stack trace |
| Error Codes | UUID only | MD5 hash of message |
| Metadata | Limited | Rich metadata in `Args` |
| Retry | Basic functions | Advanced with backoff |
| Status | Recommended: Migrate to catcher | Actively developed |

### When to Use Logger:
- ✅ Maintaining existing code that uses logger
- ✅ Need file rotation (lumberjack integration)
- ❌ New development (use catcher instead)
- ❌ Advanced error handling (use catcher instead)
