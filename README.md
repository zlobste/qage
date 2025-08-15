# qage: Post-Quantum Age Encryption

[![CI](https://github.com/zlobste/qage/actions/workflows/ci.yml/badge.svg)](https://github.com/zlobste/qage/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Post-quantum secure recipients for [age](https://age-encryption.org) encryption using hybrid X25519 + ML-KEM-768 cryptography.

## Features

- 🔒 **Post-quantum security** with hybrid X25519 + ML-KEM-768
- 🔄 **Drop-in age compatibility** - works with existing age workflows  
- 📦 **Go library** for easy integration
- 🔌 **Age plugin** support

## Quick Start

```bash
# Install
go install github.com/zlobste/qage/cmd/qage@latest

# Generate a key
qage keygen -o ~/.age/qage-key

# Get the public recipient  
qage pub -i ~/.age/qage-key
# → qage1abc...xyz

# Encrypt with age
age -R qage1abc...xyz -o secret.age secret.txt

# Decrypt
age -d -i ~/.age/qage-key secret.age
```

## Documentation

- **[CLI Reference](docs/cli.md)** - Complete command documentation

Generate CLI docs locally:
```bash
qage docs  # Creates markdown files in ./docs/
```

## Go Library

```go
import "github.com/zlobste/qage/pkg/qage"

// Generate identity
identity, err := qage.NewIdentity()

// Get recipient for encryption
recipient := identity.Recipient()

// Use with age
recipients := []age.Recipient{recipient}
identities := []age.Identity{identity}
```

## Security

qage combines two cryptographic systems:
- **X25519** (classical security) 
- **ML-KEM-768** (post-quantum security)

An attacker must break *both* to compromise your data.

## Testing

```bash
go test ./...      # Run all tests
qage selftest      # Built-in validation  
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.