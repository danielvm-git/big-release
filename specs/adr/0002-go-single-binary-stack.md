# ADR 0002: Go Single-Binary Stack Selection

## Status
Accepted

## Context
`big-release` is intended to automate releases across 35+ repositories written in varying languages (Go, TypeScript, Python, Rust, PHP, Swift, GDScript). The previous de facto standard (`semantic-release`) is built on Node.js, requiring `npm`, `npx`, or Node runtime environments even inside pure Python or Go repositories.

## Decision
We choose **Go** as the implementation language for `big-release`:
1. **Single Static Binary:** Compiles into a standalone binary with zero external runtime dependencies.
2. **Cross-Compilation:** Native Go toolchain compiles across `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, and `windows/amd64`.
3. **Execution Speed & Low Memory:** Instant startup (<20ms) and minimal RAM footprint compared to VM/interpreter runtimes (Node, Python).
4. **Reliable Standard Library:** Deep support for process execution, file system operations, and networking.

## Consequences
- **Positive:**
  - CI runners do not need Node.js or Python installed to run releases.
  - Releases execute within seconds with predictable memory usage.
  - Native concurrency and race detection (`-race`) prevent async state bugs.
- **Negative:**
  - Plugin architecture cannot rely on dynamic `eval` or runtime JavaScript module imports; plugins are compiled-in or interact via structured CLI subprocess contracts.
