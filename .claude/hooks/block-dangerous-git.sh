#!/usr/bin/env bash
# Blocks dangerous git commands for Claude Code, Cursor, and Gemini CLI hooks.
# Requires jq on PATH.
# GIT_GUARDRAILS_MODE: claude (default) | cursor | gemini
#   claude/cursor: stderr message, exit 2 on block, exit 0 on allow
#   gemini: JSON with decision on stdout, exit 0 always (allow or deny)
# GIT_BIGPOWERS_LAND=1: allows protected branch operations (set by land-branch.sh)

set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/git-guardrails-core.sh
. "$SCRIPT_DIR/lib/git-guardrails-core.sh"

INPUT=$(cat)
COMMAND=$(echo "$INPUT" | jq -r '.command // .tool_input.command // empty')
MODE="${GIT_GUARDRAILS_MODE:-claude}"

if [ -z "$COMMAND" ]; then
  if [ "$MODE" = "gemini" ]; then
    echo '{"decision":"allow"}'
  fi
  exit 0
fi

# Check dangerous patterns first
if PATTERN=$(git_guardrails_first_match "$COMMAND"); then
  REASON="BLOCKED: '$COMMAND' matches dangerous pattern '$PATTERN'. The user has prevented you from doing this."
  case "$MODE" in
    gemini)
      jq -nc --arg reason "$REASON" '{decision: "deny", reason: $reason}'
      exit 0
      ;;
    claude|cursor|*)
      echo "$REASON" >&2
      exit 2
      ;;
  esac
fi

# Check protected branch (allow only with GIT_BIGPOWERS_LAND=1)
if [ "${GIT_BIGPOWERS_LAND:-0}" != "1" ]; then
  if BRANCH=$(git_guardrails_protected_branch "$COMMAND"); then
    REASON="BLOCKED: Direct commits/pushes to protected branch '$BRANCH' are forbidden. Use a feature branch and land-branch.sh."
    case "$MODE" in
      gemini)
        jq -nc --arg reason "$REASON" '{decision: "deny", reason: $reason}'
        exit 0
        ;;
      claude|cursor|*)
        echo "$REASON" >&2
        exit 2
        ;;
    esac
  fi
fi

if [ "$MODE" = "gemini" ]; then
  echo '{"decision":"allow"}'
fi
exit 0
