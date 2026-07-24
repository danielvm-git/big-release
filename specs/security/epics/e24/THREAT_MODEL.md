# Threat Model — E24 pnpm Support

> **Epic:** e24 — pnpm Support  
> **Date:** 2026-07-24  
> **Reviewer:** security-review skill (build-epic Step 0)  
> **Scope:** New pnpm publisher mirroring npm; npm Detect precedence; shared nodeutil extraction  
> **Confidence threshold:** 8/10 (per skill mandate)

---

## Executive Summary

| Severity | Count | CWE Coverage |
|----------|-------|--------------|
| HIGH     | 0     | — |
| MEDIUM   | 2     | CWE-78, CWE-88 |
| LOW      | 2     | CWE-22, CWE-532 |

Overall residual risk after mitigations: **LOW**. The pnpm publisher reuses the same ExecCommand + package-name validation pattern as npm (already hardened against flag injection).

---

## Surface Area

| Component | Trust boundary | Sensitive data |
|-----------|----------------|----------------|
| `pnpm publish` / `pnpm view` via `ExecCommand` | Host → pnpm CLI → registry | Registry auth (`.npmrc` / env; not read by big-release directly) |
| `package.json` read/write | Working directory | Package name/version only |
| Detect markers (`pnpm-lock.yaml`, `pnpm-workspace.yaml`) | Filesystem presence only | None |
| npm Detect exclusion when pnpm markers exist | Publisher selection | Prevents double-publish |

### Data Flow

```
package.json name/version
        │
        ▼
┌─────────────────────────┐
│  pnpm Publisher          │
│  Prepare: mutate version │
│  Publish: pnpm publish   │
│    --no-git-checks       │
│    [--tag <channel>]     │
│  Verify: pnpm view name  │
│    version               │
└──────┬──────────────────┘
       │ ExecCommand (injectable)
       ▼
   pnpm CLI → npm-compatible registry
```

---

## FINDING E24-PNPM-01 — MEDIUM (Confidence 9/10)

### Command injection via crafted package name in `pnpm view`

**File:** `internal/publishers/pnpm/pnpm.go` (Verify)  
**Severity:** MEDIUM  
**Category:** OS Command Injection / Argument Injection  
**CWE:** CWE-78 / CWE-88

#### Description

`Verify` passes the package name from `package.json` into `pnpm view <name> version`. A crafted name (e.g. `--help` or `-c`) could alter CLI argument parsing if not validated.

#### Mitigation (required in implementation)

Mirror npm: validate with `isValidPackageName` before exec; prefer argv that separates options from operands (`pnpm view <name> version` with validated name only). Do not shell-interpolate.

#### Status

**Accepted with mitigation** — same pattern as npm publisher.

---

## FINDING E24-PNPM-02 — MEDIUM (Confidence 8/10)

### Double-publish if both npm and pnpm Detect

**File:** npm + pnpm `Detect()`  
**Severity:** MEDIUM  
**Category:** Business Logic / Integrity  
**CWE:** CWE-841 (adjacent)

#### Description

pnpm projects typically have `package.json` plus `pnpm-lock.yaml`. Without exclusion, both publishers Detect and may publish twice.

#### Mitigation

npm `Detect` returns false when `pnpm-lock.yaml` or `pnpm-workspace.yaml` exists. Document Detect precedence in s04 docs.

#### Status

**Accepted with mitigation** — implemented in e24s01 or e24s03.

---

## FINDING E24-PNPM-03 — LOW (Confidence 8/10)

### Working-directory package.json path traversal

**CWE:** CWE-22  

Prepare/Verify operate on `package.json` in the process CWD only (relative path). No user-supplied path. Residual risk low if CWD is attacker-controlled (CI trust model same as npm).

---

## FINDING E24-PNPM-04 — LOW (Confidence 8/10)

### Auth token leakage via stderr

**CWE:** CWE-532  

Publish wraps stderr into errors. pnpm may echo registry URLs; tokens should not appear if using standard `.npmrc`. Avoid logging full env. Same posture as npm.

---

## Residual Risk

| Story | Risk after mitigations |
|-------|------------------------|
| e24s01 | LOW — new publisher + npm Detect exclusion |
| e24s02 | LOW — tests only |
| e24s03 | LOW–MEDIUM — shared nodeutil; regression risk on npm |
| e24s04 | NONE — docs |

**Gate:** Proceed. No HIGH findings ≥ 8 confidence unresolved.
