#!/bin/bash
# Dogfood Walkthrough Script
# This script automates the npm demo package creation and release process

set -euo pipefail

echo "=== Big-Release Dogfood Walkthrough ==="
echo ""

# Configuration
REPO_DIR="big-release-demo-npm"
REPO_URL="https://github.com/danielvm-git/big-release-demo-npm.git"
BIG_RELEASE="/Users/danielvm/Developer/big-release/bin/big-release"

# Step 1: Create repository
echo "Step 1: Creating repository..."
mkdir -p "$REPO_DIR"
cd "$REPO_DIR"
git init
git checkout -b main

# Step 2: Initialize npm package
echo "Step 2: Initializing npm package..."
npm init -y

# Step 3: Create package.json
echo "Step 3: Creating package.json..."
cat > package.json << 'EOF'
{
  "name": "@danielvm-git/big-release-demo",
  "version": "0.0.0",
  "description": "Demo package for dogfooding big-release",
  "main": "index.js",
  "scripts": {
    "test": "echo \"OK\""
  },
  "keywords": ["demo", "big-release"],
  "author": "danielvm-git",
  "license": "MIT",
  "repository": {
    "type": "git",
    "url": "https://github.com/danielvm-git/big-release-demo-npm.git"
  }
}
EOF

# Step 4: Create big-release config
echo "Step 4: Creating big-release config..."
cat > .big-release.yml << 'EOF'
branches:
  - name: main

tagFormat: "v${version}"

plugins:
  - changelog
  - git

publishers:
  npm:
    enabled: true
EOF

# Step 5: Create initial files
echo "Step 5: Creating initial files..."
cat > index.js << 'EOF'
function greet(name) {
  return `Hello, ${name}!`;
}

module.exports = { greet };
EOF

cat > CHANGELOG.md << 'EOF'
# Changelog

All notable changes to this project will be documented in this file.
EOF

# Step 6: Initial commit
echo "Step 6: Creating initial commit..."
git add -A
git commit -m "feat: initial project setup"

# Step 7: Add remote
echo "Step 7: Adding remote..."
git remote add origin "$REPO_URL"

# Step 8: Validate configuration
echo "Step 8: Validating configuration..."
$BIG_RELEASE validate

# Step 9: Check health
echo "Step 9: Checking health..."
$BIG_RELEASE health

# Step 10: Dry-run release
echo "Step 10: Running dry-run release..."
$BIG_RELEASE release --dry-run

echo ""
echo "=== Setup Complete ==="
echo ""
echo "Repository created at: $(pwd)"
echo ""
echo "Next steps:"
echo "1. Push to GitHub: git push -u origin main"
echo "2. Set NPM_TOKEN: export NPM_TOKEN='your-token'"
echo "3. Run real release: $BIG_RELEASE release"
echo ""
echo "For more details, see: specs/dogfood-plan-npm.md"
