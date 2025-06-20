# Secure Email Validator

A professional command-line tool and HTTP API server written in Go for comprehensive email security validation, including domain security features like DNSSEC and STARTTLS support.

## 🚀 Features

- **Professional Email Validation**: RFC 5322 compliant email format verification
- **Domain Security Analysis**: Comprehensive security feature checking
- **Gmail Normalization**: Smart handling of Gmail aliases and dot notation
- **MX Record Verification**: Validates mail exchange record availability
- **DNSSEC Validation**: Checks DNS Security Extensions support
- **STARTTLS Verification**: Ensures encrypted email transmission capability
- **JSON API Output**: Machine-readable results for integration
- **HTTP REST API**: Professional web service interface
- **Cross-Platform Support**: Linux, macOS, and Windows compatibility
- **Detailed Reporting**: Comprehensive validation insights

## 📦 Installation

### From Source

```bash
git clone https://github.com/yourusername/secure-email-validator.git
cd secure-email-validator
make build
```

### Install System-wide

```
make install
```

## 🔧 Usage

### CLI Mode

#### Basic Validation
```
./bin/secure-email-validator -email user@company.com
```

#### Detailed Analysis
```
./bin/secure-email-validator -email user@company.com -verbose
```

#### JSON Output for Integration
```
./bin/secure-email-validator -email user@company.com -json
```

### Server Mode (REST API)

#### Start HTTP Server
```bash
./bin/secure-email-validator -server -port 8080
```

#### API Endpoints

**Email Validation**
```bash
curl "http://localhost:8080/validate?email=user@company.com&verbose=true"
```

**Service Health Check**
```bash
curl "http://localhost:8080/health"
```

## 📋 Professional Use Cases

### Enterprise Email Validation
```bash
# Corporate email verification
$ ./bin/secure-email-validator -email employee@company.com -verbose
✅ Email 'employee@company.com' is valid and secure
✨ Reason: Email is valid and domain supports secure mail delivery

--- Security Analysis ---
Domain: company.com
Has MX Record: true
Has DNSSEC: true
Primary MX Server: mail.company.com
Supports STARTTLS: true
```

### Integration Example (JSON)
```bash
$ ./bin/secure-email-validator -email user@domain.com -json
{
  "valid": true,
  "reason": "Email is valid and domain supports secure mail delivery",
  "normalized_email": "user@domain.com",
  "domain": "domain.com",
  "has_mx_record": true,
  "has_dnssec": true,
  "primary_mx_server": "mx.domain.com",
  "supports_starttls": true
}
```

## 🔒 Security Validation Features

1. **MX Record Analysis**: Verifies mail server availability
2. **DNSSEC Compliance**: Checks DNS security implementation
3. **STARTTLS Support**: Validates encrypted transmission capability
4. **Email Normalization**: Handles provider-specific formatting rules
5. **Domain Security Assessment**: Comprehensive security posture evaluation

## ⚙️ System Requirements

- Go 1.21 or higher
- `dig` utility (for DNSSEC validation)
- Network connectivity for DNS/SMTP verification

## 🛠️ Development Commands

```bash
make build      # Build application
make build-all  # Cross-platform compilation
make test       # Run test suite
make clean      # Clean build artifacts
make install    # System-wide installation
make dev        # Development server
make fmt        # Code formatting
make vet        # Code analysis
make lint       # Linting
```

## 🌐 API Integration

Perfect for:
- **Web Applications**: REST API integration
- **Microservices**: Email validation service
- **CI/CD Pipelines**: Automated email verification
- **Enterprise Systems**: Bulk email validation

## 📊 Exit Codes

- `0`: Email validation successful
- `1`: Email validation failed or security issues detected

## 🤝 Contributing

We welcome contributions! Please:

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Add comprehensive tests
5. Follow Go best practices
6. Submit a pull request

## 📄 License

MIT License - Professional use encouraged
