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

CLI command reference is auto-generated. See the markdown files in `docs/` (e.g. [`docs/qage.md`](docs/qage.md)) for the latest command help.

Generate / refresh locally:
```bash
qage docs  # Creates/updates markdown files in ./docs/
ls docs/*.md
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

qage combines two cryptographic components in a hybrid KEM:

| Component      | Purpose | Status |
| -------------- | ------- | ------ |
| X25519         | Classical ECDH security | Widely deployed |
| ML-KEM-768     | Post-quantum KEM (Kyber level 3) | NIST PQC selection |

The shared secret is derived from *both* encapsulations; an attacker must successfully break both to recover the file key. This follows the standard hybrid rationale: security degrades only if **both** primitives fail.

⚠️ Disclaimer: While ML-KEM (Kyber) is selected by NIST, real-world PQ threats and potential side-channel / implementation bugs can exist. Treat this as an additional defense layer, not a silver bullet. Review the code and perform your own audits before protecting extremely sensitive data.

## Plugin Usage (age integration)

The optional `age-plugin-qage` binary allows the standard `age` tool to encrypt/decrypt using qage recipients transparently.

Build it:
```bash
go build -o age-plugin-qage ./cmd/age-plugin-qage
mv age-plugin-qage $(go env GOPATH)/bin/  # ensure it's on PATH
```

Then `age` will automatically invoke it when encountering `qage` recipients.

## Testing

```bash
go test ./...          # Run all tests
go test -race ./...    # Race detector
qage selftest          # Built-in validation
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.