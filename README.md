<!-- markdownlint-disable MD033 -->
# ![Logz Banner](docs/assets/top_banner.png)

[![Kubex Go Dist CI](https://github.com/kubex-ecosystem/logz/actions/workflows/kubex_go_release.yml/badge.svg)](https://github.com/kubex-ecosystem/logz/actions/workflows/kubex_go_release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E=1.21-blue)](go.mod)
[![Releases](https://img.shields.io/github/v/release/kubex-ecosystem/logz?include_prereleases)](https://github.com/kubex-ecosystem/logz/releases)

---

**An advanced logging and metrics management tool with native support for Prometheus integration, dynamic notifications, and a powerful CLI.**

---

## **Table of Contents**

1. [About the Project](#about-the-project)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage](#usage)
    - [Basic Logging Examples](#basic-logging-examples)
    - [Advanced Features](#advanced-features)
    - [Production Usage Examples](#production-usage-examples)
    - [Command Line Interface](#command-line-interface)
    - [Configuration](#configuration)
5. [Performance and Benchmarks](#performance-and-benchmarks)
6. [Prometheus Integration](#prometheus-integration)
7. [Roadmap](#roadmap)
8. [Contributing](#contributing)
9. [Contact](#contact)

---

## **About the Project**

Logz is a flexible and powerful solution for managing logs and metrics in modern systems. Built in **Go**, it provides extensive support for multiple notification methods such as **HTTP Webhooks**, **WebSocket real-time streaming**, and **native Prometheus integration**, alongside seamless integration with monitoring systems for advanced observability.

Logz is designed to be robust, highly configurable, and scalable, catering to developers, DevOps teams, and software architects who need a centralized approach to logging, metrics and many other aspects of their systems.

**What makes Logz special?**

- 🚀 **Production Tested**: Successfully handles 500+ concurrent logging operations
- 🌊 **Real-time Streaming**: WebSocket support for live log monitoring
- 📊 **Native Prometheus**: Built-in metrics collection and HTTP endpoint exposure
- 🎯 **Zero Mocking**: All tests use real HTTP servers and WebSocket connections
- 🧪 **Comprehensive Testing**: 11+ integration tests covering all features
- 🔗 **Easy Integration**: Simple API that works with existing Go applications

**Why Logz?**

- 💡 **Ease of Use**: Configure and manage logs effortlessly with intuitive APIs
- 🌐 **Seamless Integration**: Drop-in replacement for standard logging with advanced features
- 🔧 **Extensibility**: Add new notifiers and services as needed
- 🏭 **Enterprise Ready**: Thread-safe, high-performance logging for production workloads

---

## **Features**

✨ **Advanced Logging System**:

- **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR with beautiful emoji formatting
- **Thread-Safe Concurrent Logging**: Handle thousands of simultaneous log operations
- **LogEntry Builder Pattern**: Flexible and intuitive log entry construction
- **Multiple Formatters**: JSON structured logs and colorful text output

🌐 **Real-Time Notifications**:

- **HTTP Webhooks**: Send logs to external services via HTTP POST
- **WebSocket Integration**: Real-time log streaming to connected clients
- **Dynamic Notifier Management**: Add/remove notifiers on the fly
- **Authentication Support**: Secure webhook delivery with token-based auth

📊 **Prometheus Integration**:

- **Native Metrics Support**: Counter, Gauge, and custom metric types
- **HTTP Metrics Endpoint**: Standard `/metrics` endpoint for Prometheus scraping
- **Metric Persistence**: Automatic saving/loading of metrics to/from disk
- **Metric Validation**: Enforces Prometheus naming conventions
- **Whitelist Support**: Control which metrics are exported

🔧 **Developer Experience**:

- **Zero Mocking Required**: All tests use real HTTP servers and WebSocket connections
- **Comprehensive Test Suite**: 11+ integration tests covering all features
- **Simple API**: Easy-to-use interfaces for quick integration
- **Context-Aware Logging**: Rich metadata support for better debugging

🚀 **Production Ready**:

- **High Performance**: Tested with 500+ concurrent operations
- **Memory Efficient**: Optimized mutex usage and resource management
- **Error Handling**: Robust error handling and recovery mechanisms
- **Configurable**: Flexible configuration options for different environments

---

## **Installation**

Requirements:

- **Go** version 1.19 or later.
- Prometheus (optional for advanced monitoring).

```bash
# Clone this repository
git clone https://github.com/kubex-ecosystem/logz.git

# Navigate to the project directory
cd logz

# Build the binary using make
make build

# Install the binary using make
make install

# (Optional) Add the binary to the PATH to use it globally
export PATH=$PATH:$(pwd)
```

---

## **Usage**

## Basic Logging Examples

### Simple Logging

```go
package main

import "github.com/kubex-ecosystem/logz"

func main() {
    // Create a new logger instance
    log := logger.NewLogger("my-app")

    // Basic logging with emoji formatting
    log.InfoCtx("Application started successfully", nil)
    log.WarnCtx("This is a warning message", nil)
    log.ErrorCtx("Something went wrong", nil)
}
```

**Output:**

```plaintext
 [INFO]  ℹ️  - Application started successfully
 [WARN]  ⚠️  - This is a warning message
 [ERROR] ❌  - Something went wrong
```

### Logging with Metadata

```go
package main

import "github.com/kubex-ecosystem/logz"

func main() {
    log := logger.NewLogger("my-service")

    // Set global metadata
    log.SetMetadata("service", "user-api")
    log.SetMetadata("version", "1.2.3")

    // Log with additional context
    log.InfoCtx("User login successful", map[string]interface{}{
        "user_id":    12345,
        "ip_address": "192.168.1.100",
        "duration":   "250ms",
    })
}
```

### Using LogEntry Builder Pattern

```go
package main

import (
    "fmt"
    "github.com/kubex-ecosystem/logz/internal/core"
)

func main() {
    // Create structured log entries
    entry := core.NewLogEntry().
        WithLevel(core.INFO).
        WithMessage("Payment processed successfully").
        AddMetadata("transaction_id", "txn_123456").
        AddMetadata("amount", 99.99).
        AddMetadata("currency", "USD").
        SetSeverity(1)

    // Use the entry with formatters
    formatter := core.NewJSONFormatter()
    output := formatter.Format(entry)
    fmt.Println(output)
}
```

## Advanced Features

### HTTP Webhook Notifications

```go
package main

import (
    "github.com/kubex-ecosystem/logz"
    "github.com/kubex-ecosystem/logz/internal/core"
)

func main() {
    // Create logger with HTTP notifier
    log := logger.NewLogger("webhook-app")

    // Add HTTP webhook notifier
    httpNotifier := core.NewHTTPNotifier(
        "https://hooks.slack.com/services/YOUR/WEBHOOK/URL",
        "your-auth-token",
    )

    // Create log entry
    entry := core.NewLogEntry().
        WithLevel(core.ERROR).
        WithMessage("Critical system error detected").
        AddMetadata("severity", "high").
        AddMetadata("component", "database")

    // Send notification
    err := httpNotifier.Notify(entry)
    if err != nil {
        log.ErrorCtx("Failed to send webhook", map[string]interface{}{
            "error": err.Error(),
        })
    }
}
```

### WebSocket Real-Time Logging

```go
package main

import (
    "fmt"
    "time"
    "github.com/kubex-ecosystem/logz/internal/core"
)

func main() {
    // Create WebSocket notifier
    wsNotifier := core.NewWebSocketNotifier("ws://localhost:8080/logs", nil)

    // Create and send log entry
    entry := core.NewLogEntry().
        WithLevel(core.INFO).
        WithMessage("Real-time log update").
        AddMetadata("timestamp", time.Now().Unix()).
        AddMetadata("event", "user_action")

    err := wsNotifier.Notify(entry)
    if err != nil {
        fmt.Printf("WebSocket notification failed: %v\n", err)
    }
}
```

### Prometheus Metrics Integration

```go
package main

import "github.com/kubex-ecosystem/logz/internal/core"

func main() {
    // Get Prometheus manager instance
    prometheus := core.GetPrometheusManager()

    // Add various metrics
    prometheus.AddMetric("http_requests_total", 100, map[string]string{
        "method": "GET",
        "status": "200",
    })

    prometheus.AddMetric("response_time_seconds", 0.045, map[string]string{
        "endpoint": "/api/users",
    })

    // Increment counter
    prometheus.IncrementMetric("api_calls_total", 1)

    // Start HTTP server to expose /metrics endpoint
    prometheus.StartHTTPServer(":2112")
}
```

### Concurrent Logging (Production Ready)

```go
package main

import (
    "sync"
    "github.com/kubex-ecosystem/logz"
)

func main() {
    log := logger.NewLogger("concurrent-app")
    var wg sync.WaitGroup

    // Simulate high-traffic logging
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            log.InfoCtx("Processing request", map[string]interface{}{
                "request_id": id,
                "worker":     "goroutine",
                "status":     "processing",
            })
        }(i)
    }

    wg.Wait()
    log.InfoCtx("All requests processed", nil)
}
```

## Production Usage Examples

### Web API with Logging and Metrics

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/kubex-ecosystem/logz"
    "github.com/kubex-ecosystem/logz/internal/core"
)

func main() {
    // Initialize logger and metrics
    log := logger.NewLogger("api-server")
    prometheus := core.GetPrometheusManager()

    // Set global metadata
    log.SetMetadata("service", "user-api")
    log.SetMetadata("version", "1.0.0")

    // Start metrics server
    go prometheus.StartHTTPServer(":2112")

    // Setup Gin router
    r := gin.Default()

    // Middleware for logging and metrics
    r.Use(func(c *gin.Context) {
        start := time.Now()

        // Process request
        c.Next()

        // Log request details
        duration := time.Since(start)
        log.InfoCtx("HTTP Request", map[string]interface{}{
            "method":     c.Request.Method,
            "path":       c.Request.URL.Path,
            "status":     c.Writer.Status(),
            "duration":   duration.String(),
            "client_ip":  c.ClientIP(),
            "user_agent": c.Request.UserAgent(),
        })

        // Update metrics
        prometheus.AddMetric("http_requests_total", 1, map[string]string{
            "method": c.Request.Method,
            "status": fmt.Sprintf("%d", c.Writer.Status()),
        })

        prometheus.AddMetric("http_request_duration_seconds",
            duration.Seconds(), map[string]string{
            "endpoint": c.Request.URL.Path,
        })
    })

    // API endpoints
    r.GET("/users/:id", getUserHandler(log))
    r.POST("/users", createUserHandler(log))

    log.InfoCtx("Server starting on :8080", nil)
    r.Run(":8080")
}

func getUserHandler(log logger.LogzLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Param("id")

        log.InfoCtx("Fetching user", map[string]interface{}{
            "user_id": userID,
            "action":  "get_user",
        })

        // Simulate database fetch
        user := map[string]interface{}{
            "id":   userID,
            "name": "John Doe",
            "email": "john@example.com",
        }

        c.JSON(http.StatusOK, user)
    }
}

func createUserHandler(log logger.LogzLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.InfoCtx("Creating new user", map[string]interface{}{
            "action": "create_user",
        })

        c.JSON(http.StatusCreated, gin.H{"status": "created"})
    }
}
```

### Microservice with WebSocket Notifications

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/kubex-ecosystem/logz"
    "github.com/kubex-ecosystem/logz/internal/core"
)

type OrderService struct {
    log      logger.LogzLogger
    notifier core.Notifier
}

func NewOrderService() *OrderService {
    log := logger.NewLogger("order-service")

    // Setup WebSocket notifier for real-time updates
    wsNotifier := core.NewWebSocketNotifier("ws://monitoring:8080/orders", nil)

    return &OrderService{
        log:      log,
        notifier: wsNotifier,
    }
}

func (s *OrderService) ProcessOrder(ctx context.Context, orderID string) error {
    s.log.InfoCtx("Processing order", map[string]interface{}{
        "order_id": orderID,
        "status":   "started",
    })

    // Simulate order processing steps
    steps := []string{"validation", "payment", "inventory", "fulfillment"}

    for i, step := range steps {
        time.Sleep(100 * time.Millisecond) // Simulate work

        // Log each step
        s.log.InfoCtx("Order step completed", map[string]interface{}{
            "order_id": orderID,
            "step":     step,
            "progress": fmt.Sprintf("%d/%d", i+1, len(steps)),
        })

        // Send real-time notification
        entry := core.NewLogEntry().
            WithLevel(core.INFO).
            WithMessage("Order progress update").
            AddMetadata("order_id", orderID).
            AddMetadata("step", step).
            AddMetadata("completed", i+1).
            AddMetadata("total", len(steps))

        s.notifier.Notify(entry)
    }

    s.log.InfoCtx("Order completed successfully", map[string]interface{}{
        "order_id": orderID,
        "status":   "completed",
    })

    return nil
}
```

### Error Handling with Webhook Alerts

```go
package main

import (
    "errors"
    "github.com/kubex-ecosystem/logz"
    "github.com/kubex-ecosystem/logz/internal/core"
)

type PaymentService struct {
    log           logger.LogzLogger
    alertNotifier core.Notifier
}

func NewPaymentService() *PaymentService {
    log := logger.NewLogger("payment-service")

    // Setup webhook for critical alerts
    alertNotifier := core.NewHTTPNotifier(
        "https://hooks.slack.com/services/TEAM/WEBHOOK/TOKEN",
        "Bearer slack-token-123",
    )

    return &PaymentService{
        log:           log,
        alertNotifier: alertNotifier,
    }
}

func (s *PaymentService) ProcessPayment(amount float64, cardToken string) error {
    s.log.InfoCtx("Processing payment", map[string]interface{}{
        "amount":     amount,
        "card_token": cardToken[:8] + "****", // Masked for security
    })

    // Simulate payment processing
    if amount > 10000 {
        // Critical error - send alert
        entry := core.NewLogEntry().
            WithLevel(core.ERROR).
            WithMessage("High-value payment failed - requires manual review").
            AddMetadata("amount", amount).
            AddMetadata("severity", "critical").
            AddMetadata("requires_action", true)

        s.alertNotifier.Notify(entry)

        s.log.ErrorCtx("Payment failed - amount too high", map[string]interface{}{
            "amount": amount,
            "reason": "exceeds_limit",
        })

        return errors.New("payment amount exceeds limit")
    }

    s.log.InfoCtx("Payment processed successfully", map[string]interface{}{
        "amount": amount,
        "status": "success",
    })

    return nil
}
```

## Command Line Interface

Here are some examples of commands you can execute with Logz's CLI:

```bash
# Log at different levels
logz info --msg "Starting the application."
logz error --msg "Database connection failed."

# Start the detached service
logz start

# Stop the detached service
logz stop

# Watch logs in real-time
logz watch
```

### CLI Usage Examples

Here are some practical examples of how to use `logz` to log messages and enhance your application's logging capabilities:

#### 1. Log a Debug Message with Metadata

```bash
logz debug \
--msg 'Just an example for how it works and show logs with this app.. AMAZING!! Dont you think?' \
--output "stdout" \
--metadata requestId=12345,user=admin
```

**Output:**

```plaintext
[2025-03-02T04:09:16Z] 🐛 DEBUG - Just an example for how it works and show logs with this app.. AMAZING!! Dont you think?
                     {"requestId":"12345","user":"admin"}
```

#### 2. Log an Info Message to a File

```bash
logz info \
--msg "This is an information log entry!" \
--output "/path/to/logfile.log" \
--metadata sessionId=98765,location=server01
```

#### 3. Log an Error Message in JSON Format

```bash
logz error \
--msg "An error occurred while processing the request" \
--output "stdout" \
--format "json" \
--metadata errorCode=500,details="Internal Server Error"
```

**Output (JSON):**

```json
{
  "timestamp": "2025-03-02T04:10:52Z",
  "level": "ERROR",
  "message": "An error occurred while processing the request",
  "metadata": {
    "errorCode": 500,
    "details": "Internal Server Error"
  }
}
```

---

#### The image below shows the CLI in action

It is demonstrating how to log messages at different levels and formats

<source src="docs/assets/demo.mp4" type="video/mp4">
  <p>Your browser does not support the video tag.</p>
</source>

---

### Description of Commands and Flags

- **`--msg`**: Specifies the log message.
- **`--output`**: Defines where to output the log (`stdout` for console or a file path).
- **`--format`**: Sets the format of the log (e.g., `text` or `json`).
- **`--metadata`**: Adds metadata to the log entry in the form of key-value pairs.

### Configuration

Logz uses a JSON or YAML configuration file to centralize its setup. The file is automatically generated on first use or can be manually configured at:
`~/.kubex/logz/config.json`.

**Example Configuration**:

```json
{
  "port": "2112",
  "bindAddress": "0.0.0.0",
  "logLevel": "info",
  "notifiers": {
    "webhook1": {
      "type": "http",
      "webhookURL": "https://example.com/webhook",
      "authToken": "your-token-here"
    }
  }
}
```

---

## Performance and Benchmarks

Logz has been thoroughly tested for production use:

- ✅ **Concurrent Operations**: Successfully handles 500+ simultaneous logging operations
- ✅ **Zero Race Conditions**: Thread-safe design with optimized mutex usage
- ✅ **Memory Efficient**: Minimal memory overhead with smart resource management
- ✅ **Network Resilient**: Robust error handling for HTTP and WebSocket failures
- ✅ **Prometheus Ready**: Metrics collection with minimal performance impact

### Real-World Usage

Currently deployed and battle-tested in production systems including:

- **GOBE Backend System**: Full-featured backend with MCP (Model Context Protocol) support
- **High-Traffic APIs**: REST APIs with thousands of requests per minute
- **Microservice Architectures**: Distributed systems with real-time monitoring
- **CI/CD Pipelines**: Automated deployment and monitoring workflows

---

## **Prometheus Integration**

Once started, Logz exposes metrics at the endpoint:

```plaintext
http://localhost:2112/metrics
```

**Example Prometheus Configuration**:

```yaml
scrape_configs:
  - job_name: 'logz'
    static_configs:
      - targets: ['localhost:2112']
```

---

## **Roadmap**

🔜 **Upcoming Features**:

- Support for additional notifier types (e.g., Slack, Discord, and email).
- Integrated monitoring dashboard.
- Advanced configuration with automated validation.

---

## **Contributing**

Contributions are welcome! Feel free to open issues or submit pull requests. Check out the [Contributing Guide](docs/CONTRIBUTING.md) for more details.

---

## **Contact**

💌 **Developer**:

Rafael Mori

- 🌐 [Portfolio](https://rafa-mori.dev)
- 🔗 [LinkedIn](https://www.linkedin.com/in/rafa-mori/)
- 📧 [Email](mailto:faelmori@gmail.com)
- 💼 Follow me on GitHub:
  - [faelmori](https://github.com/faelmori)
  - [kubex-ecosystem](https://github.com/kubex-ecosystem)

---

***I'm open to new work opportunities and collaborations. If you find this project interesting, don't hesitate to reach out!***
