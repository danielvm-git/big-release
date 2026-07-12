# Getting Started with big-release

This tutorial walks you through setting up big-release for your first project.

## Prerequisites

- Git installed
- A Git repository with conventional commits
- (Optional) Node.js for npm publishing
- (Optional) Python for PyPI publishing
- (Optional) Rust toolchain for crates.io publishing

## Step 1: Install big-release

### macOS

```bash
curl -sL https://github.com/danielvm-git/big-release/releases/latest/download/big-release-darwin-arm64 -o big-release
chmod +x big-release
sudo mv big-release /usr/local/bin/
```

### Linux

```bash
curl -sL https://github.com/danielvm-git/big-release/releases/latest/download/big-release-linux-amd64 -o big-release
chmod +x big-release
sudo mv big-release /usr/local/bin/
```

### Verify Installation

```bash
big-release --version
```

## Step 2: Configure Your Project

Create `.big-release.yml` in your project root:

```yaml
branches:
  - main
  - next
  - "N.x"
  - name: beta
    prerelease: true

tagFormat: "v${version}"

plugins:
  - changelog
  - git
  - github
```

## Step 3: Validate Configuration

```bash
big-release validate
```

## Step 4: Run a Dry Release

```bash
big-release --dry-run
```

This will show you what would happen without making any changes.

## Step 5: Make Your First Release

### Commit a Feature

```bash
git commit -m "feat: add initial functionality"
```

### Run big-release

```bash
big-release
```

This will:
1. Analyze your commits
2. Determine the next version (1.0.0)
3. Generate release notes
4. Create a git tag
5. Publish to configured registries

## Step 6: Verify the Release

Check your GitHub repository for:
- New git tag (v1.0.0)
- GitHub Release with release notes
- Published package (if configured)

## Next Steps

- [Configuration Reference](../reference/configuration.md)
- [Publisher Guide](../how-to/publishers/README.md)
- [Plugin Development](../how-to/develop-plugins.md)

## Troubleshooting

### "No releasable commits found"

This means your commits don't follow Conventional Commits. Check:
- Your commit messages start with `feat:`, `fix:`, `perf:`, etc.
- You're not using `[skip release]` in your commits

### "Branch not configured for release"

Add your branch to the configuration:

```yaml
branches:
  - your-branch-name
```

### "Authentication failed"

Ensure you have push access to the repository and proper credentials configured.
