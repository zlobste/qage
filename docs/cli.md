# qage CLI

```
qage adds post-quantum hybrid (X25519+ML-KEM-768) recipients to age.

Usage:
  qage [command]

Available Commands:
  inspect     Show identity metadata
  keygen      Generate a new qage identity
  pub         Print public recipient from identity
  selftest    Run internal quick tests
  version     Show version
```

## Commands

### keygen
Generate a new identity.

Flags:
- `-o, --output` destination file (default `-` stdout).
- `--comment` optional comment appended.
- `--pq-only` (reserved for future PQ-only suite).

### pub
Output the public recipient from an identity file.

Flags:
- `-i, --identity` path or `-` for stdin.

### inspect
Show metadata about an identity (suite, creation placeholder).

### selftest
Run internal quick self-tests (KATs/fuzz smoke in future).

### version
Print version string with commit if available.
