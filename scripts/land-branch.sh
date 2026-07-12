#!/usr/bin/env bash
# land-branch.sh — Squash-merge a feature branch to main and push.
#
# Usage: bash scripts/land-branch.sh <branch> "<conventional-commit-message>"
#
# This script:
# 1. Validates the branch name and commit message
# 2. Checks out main
# 3. Squash-merges the feature branch
# 4. Pushes to origin
# 5. Cleans up the local feature branch
#
# Environment:
#   GIT_BIGPOWERS_LAND=1 — set automatically to allow protected branch pushes

set -euo pipefail

BRANCH="${1:-}"
MESSAGE="${2:-}"

if [ -z "$BRANCH" ] || [ -z "$MESSAGE" ]; then
  echo "Usage: $0 <branch> \"<conventional-commit-message>\""
  exit 1
fi

# Validate conventional commits format
if ! echo "$MESSAGE" | grep -qE '^(feat|fix|perf|docs|chore|style|refactor|test|ci|build|revert)(\(.+\))?: .+'; then
  echo "❌ Commit message must follow Conventional Commits format:"
  echo "   <type>(<scope>): <description>"
  echo ""
  echo "   Types: feat, fix, perf, docs, chore, style, refactor, test, ci, build, revert"
  echo "   Got: $MESSAGE"
  exit 1
fi

# Check we're not already on main
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" = "main" ] || [ "$CURRENT_BRANCH" = "master" ]; then
  echo "❌ Already on $CURRENT_BRANCH. Check out a feature branch first."
  exit 1
fi

# Check the target branch exists
if ! git rev-parse --verify "$BRANCH" >/dev/null 2>&1; then
  echo "❌ Branch '$BRANCH' does not exist."
  exit 1
fi

# Stash any uncommitted changes
STASHED=false
if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "📦 Stashing uncommitted changes..."
  git stash push -m "land-branch: auto-stash before merge"
  STASHED=true
fi

# Checkout main and pull latest
echo "📥 Checking out main..."
git checkout main
git pull origin main

# Squash merge the feature branch
echo "🔀 Squash merging '$BRANCH'..."
git merge --squash "$BRANCH"

# Commit with the provided message
echo "📝 Committing..."
git commit -m "$MESSAGE"

# Push to origin
echo "📤 Pushing to origin..."
GIT_BIGPOWERS_LAND=1 git push origin main

# Delete the local feature branch
echo "🧹 Cleaning up local branch '$BRANCH'..."
git branch -d "$BRANCH"

# Restore stash if we had one
if [ "$STASHED" = true ]; then
  echo "📦 Restoring stashed changes..."
  git stash pop
fi

echo "✅ Branch '$BRANCH' landed to main."
