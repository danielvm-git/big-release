# BigBase Design System

A design system for the **BigBase Admin UI** — the control panel of BigBase, a self-hosted, single-binary Backend-as-a-Service platform (auth, database, storage, git, deploy/Sites, functions, messaging, CI/CD, monitoring). The system gives design agents everything needed to build interfaces and assets that look and feel like BigBase: **Appwrite Console clarity, BigBase indigo identity.**

The brand decision is non-negotiable: accent is **indigo `#4F46E5`** (not Appwrite pink), UI type is **Inter**, code is **Fira Code**, tone is approachable and guided.

---

## Sources

This system was built by studying the actual BigBase codebase and Appwrite's open-source design language. The reader is encouraged to explore these to build higher-fidelity BigBase designs:

- **BigBase repo** — https://github.com/danielvm-git/bigbase
  - `ui/src/index.css` — the existing ~1,400-line Appwrite-inspired token system (indigo accent, neutral 0–900, semantic roles). Mirrored into `source/bigbase-ui/` here.
  - `ui/src/components/` — `Button`, `Card`, `Input`, `Badge`, `PageHeader`, `EmptyState`, `Tabs`.
  - `ui/src/pages/` — Dashboard, DeployPage, Login, Functions, DataStudio, SQL, Users, etc.
  - `specs/006-appwrite-look-and-feel.md` — the token migration plan (Pink → BigBase indigo).
  - `specs/013-deploy.md` — the Sites/Deploy journey, build detection, and API routes.
- **Appwrite Console** (UX-pattern reference, *not* copied) — https://github.com/appwrite/console
- **Pink design system** (token/component reference) — https://github.com/appwrite/pink · docs https://pink.appwrite.io
- **Appwrite Sites quick-start** (canonical journey) — https://appwrite.io/docs/products/sites/quick-start

> The original BigBase UI is React 19 + Vite + React Router 7, embedded at `/admin/` with hash routing. We kept that reality (React + TypeScript), kept the API shapes, and rebuilt the *visual + journey* layer.

---

## What we extracted from Appwrite Console (and how BigBase adapts it)

| Appwrite Console pattern | BigBase adaptation |
|---|---|
| Sidebar grouped by **product area** with an icon system | Grouped **Overview / Build / Data / Auth / DevOps**, Lucide icons, indigo-tint active state |
| Guided **Create** wizards (sites, functions) with a step rail | `WizardSteps` rail: Source → Configure → Deploy, no dead ends |
| Friendly **empty states** with an illustration-like icon + single CTA | `EmptyState` — icon chip, invitation title, payoff body, one primary button |
| **Card + list** views with status pills | `SiteCard` grid + status `Badge` semantics (ready/building/failed) |
| Low-noise dense pages — quiet neutrals, color = meaning | Neutral 0–900 ramp; accent reserved for primary action/active/focus |
| Settings density, progressive disclosure | Collapsible build settings & env vars in the wizard |
| **Inter** typography, calm spacing rhythm | Inter UI + 4px spacing base + emphasized easing |
| Real **dark mode** | Single `[data-theme="dark"]` role remap |

The full IA rationale, token mapping, component inventory, journey maps, a11y, responsive, and React refactor plan live in **`SYSTEM_DESIGN.md`**.

---

## Content fundamentals

How BigBase writes copy:

- **Person & tense:** second person, present tense. The product talks *to* the developer. "Where's your code?", "Your site is live", "Connect a Git repository and BigBase builds, deploys, and serves it."
- **Tone:** approachable and guided — a knowledgeable teammate, not a manual. Confident but never jargon-heavy. The goal of every screen is "get something running."
- **Casing:** sentence case for headings and body. UPPERCASE only for overlines / section labels (`OVERVIEW`, `LATEST COMMIT`) and badges. Title Case inside status badges (Ready, Building, Failed).
- **Buttons:** verb-first and specific — *Create site*, *Authorize GitHub*, *Deploy*, *Redeploy*, *Add domain*, *Save changes*. Never "Submit" or "OK".
- **Empty states:** title is an invitation ("Create your first site"); body states the payoff; one primary CTA.
- **Errors:** plain and actionable, surfaced near the field or as a toast — "Password must be at least 6 characters.", "Build failed — exit code 1. Check logs." No raw stack traces.
- **Status words:** single words — Ready / Building / Failed / Pending — always paired with a color *and* a dot/spinner.
- **Emoji:** essentially none. One celebratory exception is tolerated at a genuine milestone ("Your site is live 🎉"). Never decorative, never in nav or labels.
- **Numbers & code:** monospace (Fira Code) for URLs, commit SHAs, branches, commands, IDs.

---

## Visual foundations

- **Color:** an indigo-on-neutral system. Accent `#4F46E5` (600 hover `#4338CA`, 700 active `#3730A3`) is reserved for primary buttons, active nav, focus rings, and links — never decorative fills. Everything else is the **neutral 0–900** ramp. Status uses green/amber/red/blue, each as a solid + a ~10–12% tinted background + a darkened (light) or lightened (dark) foreground. No gradients in chrome; the only gradients are the playful **site-card thumbnails** (indigo, slate, rose, emerald) standing in for site screenshots.
- **Accent themes:** indigo is the default, but the accent is a swappable layer. **12 month themes** (`data-accent="january…december"`) replace the brand ramp + accent role tokens; each is verified **WCAG 2.1 AA** (≥4.5:1) on light and dark surfaces — light/pastel seeds flip button text to near-black, and the June "Rainbow" theme uses a gradient fill with a text-shadow while keeping a solid-purple AA fallback for links/nav/focus. Switch them from the sidebar footer; the choice persists to `localStorage('bigbase-theme')`. Full contrast table + mechanism in `SPECS.md §3`.
- **Typography:** **Inter** for all UI (400/450/500/600/700, tracking −0.02em on headings); **Fira Code** for code/SQL/logs/identifiers. Type scale 12 → 40px. Generous line-height (1.55) on body.
- **Spacing & layout:** 4px base scale. Content max-width ~1180px beside a 252px sidebar. Calm, roomy density — Appwrite-like breathing room, not a cramped admin grid. Flex/grid with `gap` everywhere.
- **Backgrounds:** flat. App background is `--bg-default` (near-white / near-black); surfaces are `--bg-surface`. No textures, no patterns, no full-bleed photography. Imagery is limited to the indigo brand mark and gradient site thumbnails.
- **Corner radii:** xs 4 (badges/sm buttons), s 8 (buttons/inputs), m 12 (cards), l 16 (login/modals), full (pills/avatars).
- **Cards:** `--bg-surface` + 1px `--border-default` + `--radius-m` + soft `--shadow-s`. They lift on hover (`translateY(-2px)` + `--shadow-m`) only when interactive (e.g. site cards). No colored left-border accents.
- **Borders & shadows:** hairline neutral borders do most of the structural work; shadows are subtle (xs–l) and reserved for elevation (cards, dropdowns, toasts, modals). A dedicated `--focus-ring` (3px indigo @ 18%) communicates focus.
- **Motion:** quick and purposeful. Durations 150–300ms; the signature easing is `cubic-bezier(.32,.72,0,1)` (emphasized) for things that travel (toasts, progress, wizard steps). Fades + small translates; no bounces, no springy overshoot. Skeleton shimmer for loading.
- **Hover / press states:** hover lightens surfaces (`--overlay-hover`) or shifts the accent one step darker; primary buttons darken 500→600. Press nudges buttons ~0.5px down (no scale-shrink). Active nav gets an indigo-tint background + indigo text.
- **Transparency & blur:** sparing. Tinted backgrounds (accent/status @ 10–12%) and overlay scrims for modals. No glassmorphism/backdrop-blur in the core chrome.
- **Iconography:** Lucide (see below). 2px stroke, round caps — friendly and open.
- **Imagery vibe:** there is almost none by design. Where a "screenshot" is implied (site cards, site detail hero), we use a flat brand-colored gradient block — cool indigos as the default, with slate/rose/emerald accents per site.

---

## Iconography

- **System:** [**Lucide**](https://lucide.dev) — open-source (ISC), 2px stroke, round line caps/joins, 24×24 grid. Chosen for a friendly, consistent, *open* set that fits BigBase's approachable tone.
- **Why a substitution:** the original BigBase sidebar used **placeholder single-letter icons** (`H`, `D`, `S`, `λ`…) and Appwrite Console ships a **proprietary icon font** we deliberately don't copy. Lucide replaces both with a cohesive system. **⚠️ Flag:** if BigBase later standardizes on a specific icon font, swap the `assets/icons.jsx` path map.
- **Delivery:** icons are inlined as SVG path data in **`assets/icons.jsx`** (an `<Icon name="rocket" size={18} />` React component), so prototype/kit files stay self-contained with no CDN dependency. The same names map cleanly to the npm `lucide-react` package for production. Key names in use: `layout-dashboard, database, terminal, hard-drive, rocket, box, git-branch, git-pull-request, activity, users, settings, github, globe, plus, check, check-circle, chevron-right/left/down, arrow-left/right, external-link, refresh-cw, search, x, alert-triangle, info, clock, more-horizontal, play, copy, log-out, moon, sun, mail, zap, cpu, bell`.
- **Brand mark:** `assets/bigbase-logo.svg` — an indigo (`#4F46E5`) rounded square (28px radius on a 120px tile) with a white "B", reproducing the existing in-product mark (`.sidebar-logo-icon`). Pair with the Inter 700 "BigBase" wordmark for the lockup.
- **Emoji / unicode as icons:** not used in chrome. (One celebratory 🎉 is allowed in a success moment only.)

---

## Index — what's in this folder

| Path | What it is |
|------|------------|
| `README.md` | This file — context, sources, content + visual foundations, iconography, index |
| `SYSTEM_DESIGN.md` | Full design doc: principles, IA + cross-links, token mapping, component inventory, journeys A–D, a11y, responsive, React refactor plan |
| `SPECS.md` | **Engineering handoff** — token tables, 12-theme WCAG verification, theme switching + localStorage, component-state matrix, responsive breakpoints, dark-mode remap, a11y checklist |
| `colors_and_type.css` | **Source of truth** — all design tokens (color, type, spacing, radius, motion, shadow) + 12 month themes + semantic type styles + dark mode |
| `components.css` | Reusable component classes (button, card, input, badge, tabs, sidebar, choice-card, wizard-steps, site-card, preview-banner, toast, skeleton, …) |
| `assets/bigbase-logo.svg` | Brand mark |
| `assets/icons.jsx` | Lucide `<Icon>` component (inlined SVG paths) |
| `preview/` | Design-system specimen cards (type, color, spacing, components, brand) shown in the Design System tab |
| `ui_kits/admin-console/` | **Interactive prototype / UI kit** — clickable BigBase Console (login, dashboard, Sites flow, site detail, functions, data studio). Open `index.html`. |
| `source/bigbase-ui/` | Read-only mirror of the original BigBase `ui/src` (tokens, components, pages) for reference |
| `SKILL.md` | Agent-skill manifest so this system can be used in Claude Code |

### Fonts
Inter (UI) and Fira Code (code) load automatically from the **Google Fonts CDN** (`@import` at the top of `colors_and_type.css`) — both are free, open-source, and require **no upload or licensing**. Inter is the brand-mandated UI face, so nothing needs swapping. Only if you later want a fully offline/self-hosted build, drop the woff2 files into a `fonts/` folder and replace the `@import` with `@font-face` rules.
