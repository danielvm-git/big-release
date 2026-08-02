# release.bigbase.click

The launch and documentation site for big-release. A self-contained Go app that serves a polished single-page site showcasing:

- Live proof (4 canary apps with real version badges)
- How big-release works (3-step release flow)
- Installation instructions (macOS, Linux, CI)
- Per-language quickstart walkthroughs (Go, Node, Python, PHP) with setup files and copy-to-agent Markdown
- Comparison with semantic-release
- Supported languages and registries

## Build & Run

```bash
# Local dev
go run .

# Tests
go test ./... -count=1 -timeout 120s

# Vet
go vet ./...
```

The app reads `VERSION` at runtime (set by big-release in CI) and serves the embedded `index.html` on the PORT env var (default 8080).

## Deployment

Deployed to `https://release.bigbase.click` via BigBase on every merge to main in the big-release repo.
