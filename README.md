<div align="center">

# Bloombox

![Bloombox Cover](assets/cover.png)

A free, open-source email validation API made for Indie Hackers — no auth, no limits, no paywalls, no data stored. Just fast, accurate checks for syntax, domains, MX records, and throwaways. Refreshed daily.

</div>

## Features

- **Real-time email validation API** with JSON responses
- **Batch email validation** (up to 100 emails per request)
- **Multiple validation checks**:
  - **Syntax validation** - RFC 5322 compliant email format checking
  - **MX record validation** - DNS MX record verification
  - **SMTP validation** - Real-time mailbox verification
- **Disposable email detection** - Blocks temporary email services
  - **Free email provider detection** - Identifies free email domains
  - **Role-based email detection** - Blocks generic role emails (admin@, info@, etc.)
  - **Banned words filtering** - Custom word blacklist in email usernames
  - **Email blacklist** - Custom email address blacklist
  - **Domain blacklist** - Custom domain blacklist
  - **Gravatar validation** - Checks if email has a Gravatar account
- **Configurable validators** - Enable/disable specific validators at runtime
- **Built-in caching** with LRU eviction for performance
- **Concurrent processing** with rate limiting
- **Health monitoring** and validation statistics
- **Custom filter support** - Cuckoo filters for large datasets with configurable false positive rates

## Installation

```bash
go mod download
```

## Data Files

The service includes several pre-configured data files in the `/data` directory for various validation purposes:

### Available Data Files

| File             | Size  | Domains | Description                                         |
| ---------------- | ----- | ------- | --------------------------------------------------- |
| `disposable.txt` | 51KB  | 4,020   | Disposable email domains (temporary email services) |
| `free.txt`       | 117KB | 8,766   | Free email provider domains (Gmail, Yahoo, etc.)    |
| `hubspot.txt`    | 70KB  | 4,769   | Non-company domains (personal email domains)        |
| `skiplist.txt`   | 124KB | 9,273   | Merged list of free and hubspot domains             |

### Data File Usage

These files are used by the corresponding validators:

- **`disposable.txt`** → `disposable` validator
- **`free.txt`** → `free` validator
- **`hubspot.txt`** → Can be used for `blacklist_domains` validator
- **`skiplist.txt`** → Comprehensive blacklist for `blacklist_domains` validator

### Configuration

To use these data files, set the appropriate environment variables:

```bash
export DISPOSABLE_EMAILS_FILE=data/disposable.txt
export FREE_EMAILS_FILE=data/free.txt
export BLACKLIST_DOMAINS_FILE=data/skiplist.txt
```

The files contain one domain per line and are automatically loaded when the corresponding validators are enabled.

## Usage

### Start the server

> **Note:** SMTP validation may require rotating proxies to avoid being flagged by upstream mail servers. Some providers may block repeated validation attempts from the same IP address.

```bash
go run cmd/server/main.go
```

The server runs on port 8080 by default. Set the `PORT` environment variable to use a different port.

### API Endpoints

| Method | Endpoint            | Description                                 |
| ------ | ------------------- | ------------------------------------------- |
| `GET`  | `/`                 | API documentation and available validators  |
| `POST` | `/validate`         | Validate single email address               |
| `POST` | `/batch`            | Validate multiple email addresses (max 100) |
| `GET`  | `/validators`       | List available validators and their status  |
| `PUT`  | `/validators/:name` | Enable/disable specific validator           |
| `GET`  | `/health`           | Health check endpoint                       |

### Examples

#### Single Email Validation

```bash
curl -X POST http://localhost:8080/validate \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "validators": ["syntax", "mx", "smtp"]
  }'
```

Response:

```json
{
  "email": "user@example.com",
  "timestamp": "2024-01-15T10:30:00Z",
  "duration": "1.2s",
  "results": {
    "syntax": {
      "valid": true,
      "message": "Valid email syntax",
      "duration": "0.1ms"
    },
    "mx": {
      "valid": true,
      "message": "MX records found",
      "duration": "50ms"
    },
    "smtp": {
      "valid": true,
      "message": "Mailbox can receive emails",
      "duration": "1.1s"
    }
  },
  "is_valid": true,
  "summary": {
    "is_disposable": false,
    "is_free": false,
    "is_role": false
  }
}
```

#### Batch Email Validation

```bash
curl -X POST http://localhost:8080/batch \
  -H "Content-Type: application/json" \
  -d '{
    "emails": ["user1@example.com", "user2@test.com", "admin@company.org"],
    "validators": ["syntax", "disposable", "free", "role"]
  }'
```

#### Enable/Disable Validator

```bash
curl -X PUT http://localhost:8080/validators/smtp \
  -H "Content-Type: application/json" \
  -d '{"enabled": false}'
```

## Configuration

Configure the service using environment variables:

### Server Configuration

- `PORT` - Server port (default: 8080)

### File-based Validators

- `FREE_EMAILS_FILE` - Path to file containing free email domains (one per line)
- `DISPOSABLE_EMAILS_FILE` - Path to file containing disposable email domains
- `ROLE_EMAILS_FILE` - Path to file containing role-based email usernames
- `BAN_WORDS_FILE` - Path to file containing banned words for email usernames
- `BLACKLIST_EMAILS_FILE` - Path to file containing blacklisted email addresses
- `BLACKLIST_DOMAINS_FILE` - Path to file containing blacklisted domains

### Performance Settings

- `CACHE_SIZE` - LRU cache size (default: 1000)
- `VALIDATION_TIMEOUT` - Overall validation timeout (default: 5s)
- `SMTP_TIMEOUT` - SMTP validation timeout (default: 5s)
- `FALSE_POSITIVE_RATE` - Cuckoo filter false positive rate (default: 0.01)

### SMTP Configuration

- `SMTP_FROM_DOMAIN` - Domain to use for SMTP FROM (default: example.com)
- `SMTP_FROM_EMAIL` - Email to use for SMTP FROM (default: test@example.com)

### Validator Control

- `ENABLED_VALIDATORS` - Comma-separated list of validators to enable (default: syntax)

### Example Configuration

```bash
export PORT=9090
export FREE_EMAILS_FILE=/data/free_emails.txt
export DISPOSABLE_EMAILS_FILE=/data/disposable_emails.txt
export ROLE_EMAILS_FILE=/data/role_emails.txt
export ENABLED_VALIDATORS=syntax,mx,disposable,free,role
export CACHE_SIZE=5000
export VALIDATION_TIMEOUT=10s
```

## Available Validators

| Validator           | Description                         | Requires File | Default  |
| ------------------- | ----------------------------------- | ------------- | -------- |
| `syntax`            | RFC 5322 email format validation    | No            | Enabled  |
| `mx`                | DNS MX record verification          | No            | Disabled |
| `smtp`              | Real-time SMTP mailbox verification | No            | Disabled |
| `disposable`        | Disposable email provider detection | Yes           | Disabled |
| `free`              | Free email provider detection       | Yes           | Disabled |
| `role`              | Role-based email detection          | Yes           | Disabled |
| `banwords`          | Banned words in email username      | Yes           | Disabled |
| `blacklist_emails`  | Blacklisted email addresses         | Yes           | Disabled |
| `blacklist_domains` | Blacklisted domains                 | Yes           | Disabled |
| `gravatar`          | Gravatar account existence check    | No            | Disabled |

## Build

```bash
go build -o bloombox cmd/server/main.go
```

## Dependencies

- **Go 1.24.2+**
- **github.com/hashicorp/golang-lru/v2** - LRU caching
- **github.com/linvon/cuckoo-filter** - Space-efficient filtering
- **go.uber.org/zap** - Structured logging

## Architecture

The service uses a modular architecture with:

- **Validator Interface** - Pluggable validation components
- **Filter System** - Configurable filtering (Map vs Cuckoo Filter)
- **Caching Layer** - LRU cache for validation results
- **Concurrent Processing** - Semaphore-based rate limiting
- **Configuration Management** - Environment-based configuration

## Performance

- **Concurrent validation** with configurable limits
- **LRU caching** for repeated email checks
- **Cuckoo filters** for memory-efficient large datasets
- **Timeout controls** to prevent hanging validations
- **Batch processing** for high-throughput scenarios
