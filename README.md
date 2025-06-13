![Logz Banner](docs/assets/top_banner.png)

---

**An advanced logging and metrics management tool with native support for Prometheus integration, dynamic notifications, and a powerful CLI.**

---

## **Table of Contents**
1. [About the Project](#about-the-project)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage](#usage)
    - [CLI](#cli)
    - [Configuration](#configuration)
5. [Prometheus Integration](#prometheus-integration)
6. [Roadmap](#roadmap)
7. [Contributing](#contributing)
8. [Contact](#contact)

---

## **About the Project**
Logz is a flexible and powerful solution for managing logs and metrics in modern systems. Built in **Go**, it provides extensive support for multiple notification methods such as **HTTP Webhooks**, **ZeroMQ**, and **DBus**, alongside seamless integration with **Prometheus** for advanced monitoring.

Logz is designed to be robust, highly configurable, and scalable, catering to developers, DevOps teams, and software architects who need a centralized approach to logging, metrics and many other aspects of their systems.

**Why Logz?**
- 💡 **Ease of Use**: Configure and manage logs effortlessly.
- 🌐 **Seamless Integration**: Easily integrates with Prometheus and other systems.
- 🔧 **Extensibility**: Add new notifiers and services as needed.

---

## **Features**
✨ **Dynamic Notifiers**:
- Support for multiple notifiers simultaneously.
- Centralized and flexible configuration via JSON or YAML.

📊 **Monitoring and Metrics**:
- Exposes Prometheus-compatible metrics.
- Dynamic management of metrics with persistence support.

💻 **Powerful CLI**:
- Straightforward commands to manage logs and services.
- Extensible for additional workflows.

🔒 **Resilient and Secure**:
- Validates against Prometheus naming conventions.
- Distinct modes for standalone and service execution.

---

## **Installation**
Requirements:
- **Go** version 1.19 or later.
- Prometheus (optional for advanced monitoring).

```bash
# Clone this repository
git clone https://github.com/rafa-mori/logz.git

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

### CLI
Here are some examples of commands you can execute with Logz’s CLI:

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

### **Usage Examples**

Here are some practical examples of how to use `logz` to log messages and enhance your application's logging capabilities:

#### **1. Log a Debug Message with Metadata**

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

#### **2. Log an Info Message to a File**

```bash
logz info \
--msg "This is an information log entry!" \
--output "/path/to/logfile.log" \
--metadata sessionId=98765,location=server01
```

#### **3. Log an Error Message in JSON Format**

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

#### The image below shows the CLI in action, demonstrating how to log messages at different levels and formats:

![img.png](docs/assets/cli_print.png)

---

### **Description of Commands and Flags**
- **`--msg`**: Specifies the log message.
- **`--output`**: Defines where to output the log (`stdout` for console or a file path).
- **`--format`**: Sets the format of the log (e.g., `text` or `json`).
- **`--metadata`**: Adds metadata to the log entry in the form of key-value pairs.

---

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

## **Prometheus Integration**
Once started, Logz exposes metrics at the endpoint:
```
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
[Rafael Mori](mailto:faelmori@gmail.com)
💼 [Follow me on GitHub](https://github.com/rafa-mori)
I'm open to new work opportunities and collaborations. If you find this project interesting, don’t hesitate to reach out!
