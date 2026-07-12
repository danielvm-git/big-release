# Shared pattern list and matcher for git hook guardrails.
# Source from block-dangerous-git.sh only.

GIT_GUARDRAILS_PATTERNS=(
  "git reset --hard"
  "git clean -fd"
  "git clean -f"
  "git branch -D"
  "git checkout \\."
  "git restore \\."
  "push --force"
  "reset --hard"
)

# Protected branches — direct commits/pushes forbidden unless GIT_BIGPOWERS_LAND=1
GIT_GUARDRAILS_PROTECTED_BRANCHES="main|master"

# Print first matching pattern to stdout; exit 0 if dangerous, 1 if safe.
git_guardrails_first_match() {
  local cmd="$1"
  local p
  for p in "${GIT_GUARDRAILS_PATTERNS[@]}"; do
    if echo "$cmd" | grep -qE "$p"; then
      printf '%s' "$p"
      return 0
    fi
  done
  return 1
}

# Check if command targets a protected branch (push or commit on main/master).
# Returns 0 (match) with branch name on stdout, or 1 (no match).
git_guardrails_protected_branch() {
  local cmd="$1"
  local branch

  # Match: git push origin main/master, git push main/master
  if branch=$(echo "$cmd" | grep -oE "push\s+(-[a-z]+\s+)?(origin\s+)?(${GIT_GUARDRAILS_PROTECTED_BRANCHES})\b" | grep -oE "${GIT_GUARDRAILS_PROTECTED_BRANCHES}$"); then
    printf '%s' "$branch"
    return 0
  fi

  # Match: git commit ... on protected branch (current branch check done in caller)
  if echo "$cmd" | grep -qE "^git commit\b"; then
    local current
    current=$(git branch --show-current 2>/dev/null || true)
    if echo "$current" | grep -qE "^(${GIT_GUARDRAILS_PROTECTED_BRANCHES})$"; then
      printf '%s' "$current"
      return 0
    fi
  fi

  return 1
}
