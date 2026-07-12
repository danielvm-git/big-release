# Tech Stack

## Language & Runtime
- **Go** (single binary distribution)
- **YAML** (configuration)
- **Bash** (scripts)

## Dependencies
- **cobra** — CLI framework
- **gopkg.in/yaml.v3** — YAML parsing

## Architecture
Three layers:
1. **Algorithm** (`internal/algorithm/`) — commit analysis, version calculation, notes generation
2. **Git Operations** (`internal/git/`) — commit retrieval, tag creation, push/pull
3. **Publishers** (`internal/publishers/`) — language-specific package publishing
4. **CLI** (`cmd/big-release/`) — command-line interface

## Build
- `make build` — compile binary
- `make test` — run tests
- `make lint` — run linter
- `make preflight` — full verification stack
