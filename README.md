# HostBind 

**Local dev environments that AI agents can actually trust.**
*Install once. Every project, every worktree, every agent gets its own conflict-free ports and URLs — and agents ask instead of guessing.*

Status: Concept → v0.1 spike · License: Apache-2.0 (or MIT) · Platforms: macOS, Linux, Windows (v0.1: macOS/Linux, Windows by v0.3)

---

## 1. Problem (short version)

Run 3 projects and 3 AI agents at once and you get: two apps on `:3000`, `VITE_API_URL` pointing at a dead port, CORS errors, stale OAuth callbacks, and agents that *assume* `localhost:5173`.

Root cause: **nothing owns the answer to "where is service X of project Y right now?"** HostBind owns it.

**Core thesis:** *Agents should discover the environment, not guess it.*

## 2. Positioning (honest)

Port allocation alone is a solved-ish problem, and there are existing tools for stable `.localhost` names, per-worktree ports, and process running (Procfile runners, Tilt, process-compose, mprocs, devcontainers, direnv/mise). **Before v0.1 ships, try each and write `docs/alternatives.md` with a fair comparison table.**

HostBind's differentiator is the combination nobody else focuses on:

1. **Agent-first interface** — `context --json`, skill, `AGENTS.md`, MCP
2. **Cross-project registry** — one source of truth across *all* projects, not one repo
3. **Diagnosis** — `doctor` explains *why* something is unreachable/stale
4. **Zero-config, single binary, local-only**

Positioning line: *"Not a port manager. The service registry your AI agents read before they touch localhost."*

## 3. Principles

1. **AI expresses intent; HostBind decides ports.** Agents never pick ports.
2. **Stable by default, conflict-free always.**
3. **One source of truth**: the registry.
4. **Deterministic core, thin agent layer.**
5. **Never surprise the user**: no silent file rewrites, never kill unrelated processes.
6. **Local-first**: no account, no telemetry by default, no network calls.
7. **Works with existing tools**, doesn't replace them.

## 4. Key design decisions (changed from v1)

| Topic | v1 | v2 decision | Why |
|---|---|---|---|
| Language | undecided | **Go** single static binary (Rust also fine) | No Node/Python dependency, easy cross-platform releases, easy contributor onboarding |
| Daemon | required | **Daemonless in v0.1–v0.2.** SQLite + file locks; CLI does work on demand. Daemon only when proxy/health-watching needs it (v0.4) | Removes install friction and a whole failure class |
| URLs | `.localhost` in MVP | **Ports-only in v0.1.** Proxy + `.localhost` in v0.4 | Proxy brings port 80, HMR/WebSocket, cookie and OAuth complications |
| Port forcing | glossed over | **Adapter ladder** (see §6) | This is the hard part; it's the real core |
| MCP | far future | **v0.3** | Cheap to build, highest agent-adoption payoff |
| Race handling | "prevent races" | **Leases** in registry + OS bind as final arbiter + retry | Check-then-bind gap is real |
| Agent adoption | skill only | Skill **+ `AGENTS.md`/`CLAUDE.md` snippet + shell/script wrappers + `doctor` detects rogue processes** | Skills load on relevance; project instruction files load every session |
| `fix --yes` | open | **Allowlist of fix types, diff logged, `hostbind undo`** | LLM must not make unreviewed arbitrary edits |
| Env prefix | `DEVmesh_` | **`HOSTBIND_`** | Consistency |

## 5. Architecture (v0.1–v0.3)

```
 Agent (Claude Code / Codex / OpenCode / Kilo / Cursor…)
   │  Skill · AGENTS.md · MCP (v0.3)
   ▼
 hostbind CLI  ── human-friendly + --json
   │
   ├── Scanner + Adapters   (detect framework, know how to force a port)
   ├── Allocator            (ranges, leases, persistence)
   ├── Runner               (spawns commands with injected env, tracks PID)
   ├── Doctor               (read-only checks, proposes fixes)
   └── Registry             (SQLite, WAL mode, file-locked)
```

A daemon (v0.4) adds: health watching, reverse proxy, dashboard.

## 6. The port-forcing problem: Adapter ladder

For each service, an adapter picks the **least invasive** working method, in this order:

1. **Env var** (`PORT`, `VITE_PORT`, …) — no file changes
2. **CLI flag append** (`--port 4312`, `-p`, `--host`) — no file changes
3. **Generated overlay config** outside the repo (e.g., a temp config passed via flag)
4. **Config edit with explicit consent** — shown as a diff, reversible with `undo`
5. **Reverse proxy** (v0.4) — app binds anywhere, HostBind routes

Adapter interface (keep it tiny so contributors can add one in an afternoon):

```
detect(dir)      → bool + confidence
services(dir)    → [{name, type, defaultCommand}]
portArgs(port)   → {env: {...}, args: [...]}
envFor(service, registry) → {VITE_API_URL: ..., ...}
health(service)  → probe spec (optional)
```

Adapters ship in-tree as YAML/declarative where possible (most frameworks need only "which env/flag sets the port"), and code only when necessary. **Declarative adapters = the easiest contribution path.**

v0.1 adapters: `generic` (env `PORT`), `vite`, `next`, `fastapi/uvicorn`, `express`.
Later: Django, Flask, Astro, Nuxt, SvelteKit, Remix, Rails, Laravel, Spring, Go, Rust, Bun, Deno, docker-compose.

## 7. Registry & identity

- SQLite, WAL mode, one file at `~/.hostbind/registry.db`. Atomic allocation via transaction + lock.
- **Project identity** = `git remote URL + relative path` if available, else absolute path. Two clones of one repo are distinct *instances* of one project. Moved folder → `hostbind doctor` offers re-link.
- **Leases:** a port is *reserved* with a TTL before the process binds; the OS bind is the final arbiter; on bind failure, reallocate and retry.
- **Garbage collection:** `hostbind gc` removes entries for deleted paths/worktrees; `doctor` warns on orphaned entries.
- Registry schema versioned + migrations from day one.
- Generated env files (`.env.hostbind`) are **gitignored automatically** and never committed — ports differ per machine.

Port ranges (configurable in `~/.hostbind/config.yaml`):

```yaml
port_ranges: { web: 4300-4799, api: 5300-5799, worker: 6300-6799, other: 7300-7799 }
```

## 8. Commands

**v0.1 (spike → usable):**

```bash
hostbind run [--name api] -- <command>   # allocate port, inject env, run, register
hostbind ls                              # all projects/services/ports/status
hostbind context [--json]                # what agents call first
hostbind port <service>                  # print the port/URL (scriptable)
hostbind stop [service]
hostbind install-skill                   # skill + AGENTS.md snippet
hostbind gc
```

**v0.2:** `init`/`scan` (autodetect → `.hostbind/config.yaml`), `up` (start all services in a project), `doctor` (read-only), env sync via `.env.hostbind`, `restart`, `status`.

**v0.3:** MCP server (`hostbind mcp`), `fix` (allowlisted, with diff + `undo`), Windows, worktree support, `logs`.

**v0.4:** daemon, reverse proxy + `*.localhost` URLs, health monitoring, dependency ordering, Docker Compose integration, local dashboard.

Every command supports `--json`; exit codes are stable and documented (agents depend on them).

## 9. Environment contract

```
HOSTBIND_PROJECT           solar-mapper
HOSTBIND_INSTANCE          main            # worktree/branch instance
HOSTBIND_SERVICE           api
HOSTBIND_PORT              5314
HOSTBIND_URL               http://localhost:5314
HOSTBIND_<SERVICE>_URL     e.g. HOSTBIND_WEB_URL, HOSTBIND_API_URL
```

Framework-specific vars (e.g. `VITE_API_URL`, `NEXT_PUBLIC_API_URL`) come from adapters, written to `.env.hostbind` — never into the user's own `.env` unless they opt in.
Values are labeled **managed / generated / detected / user-defined**. Secrets are never printed or logged.

## 10. Agent integration (the growth engine)

**Three layers, so agents comply even if one is ignored:**

1. **`AGENTS.md` / `CLAUDE.md` snippet** (auto-added by `install-skill`, with consent) — always in context:
   > This repo uses HostBind. Never choose or assume dev ports (3000, 5173, 8000…). Run `hostbind context --json` for service URLs. Start services with `hostbind up` / `hostbind run`. If something is unreachable, run `hostbind doctor`.
2. **Skill** (`skills/hostbind-local-development/SKILL.md`) with references for commands, env contract, troubleshooting.
3. **MCP tools** (v0.3): `hostbind_context`, `hostbind_start`, `hostbind_stop`, `hostbind_doctor`, `hostbind_logs`.

**Rules taught to agents:** never invent ports; query before starting; use provided URLs; don't hand-edit managed env; run `doctor` on failure; never `kill` unknown processes.

**Compliance measurement:** publish an open **agent-compliance benchmark** (`bench/`): scripted tasks like "start the app and call the API" run against several agents, measuring how often they use HostBind vs guess a port. This is both a quality tool and great content for launch posts.

## 11. Doctor & Fix

`doctor` (read-only) checks: registry vs reality, process alive, port reachable, HTTP health, env values pointing at stale ports, CORS/proxy origins referencing old ports, services running *outside* HostBind, orphaned entries, port conflicts with named PID/command.

Output states **what's wrong · why · what would change · how to fix**, and has `--json`.

`fix`: interactive by default. `--yes` only applies **allowlisted** safe fix types (env URL updates in generated files, registry repairs). Every change writes a diff to `~/.hostbind/history/` and is reversible via `hostbind undo`. Anything touching user-authored files always requires explicit confirmation.

## 12. Safety & trust model

- **Detected → Approved → Executed** commands. Cloning an untrusted repo and running `hostbind up` must not silently execute arbitrary scripts: first run shows the command list and asks for approval; approvals stored per project *and per command hash* (re-prompt if it changes).
- Never kill processes HostBind didn't start; offer alternatives instead.
- Destructive commands (`clean --all`, `undo`) print exactly what they'll remove.
- Registry is user-only permissions (0600); no network access; no telemetry (optional opt-in anonymous metrics later, off by default, documented).
- `SECURITY.md` with private disclosure process.

## 13. Worktrees & multi-agent (v0.3)

Each git worktree = a separate **instance** with its own allocations, and env `HOSTBIND_INSTANCE=<branch>`. `hostbind worktree add <branch>` wraps `git worktree add` and pre-allocates. Parallel agents in parallel worktrees never collide. Cleanup via `gc` when worktrees are removed.

## 14. Other things v1 missed (now decided)

- **Monorepos:** scanner reads workspaces (`pnpm-workspace.yaml`, `package.json` workspaces, `turbo.json`, `nx.json`) and creates one service per runnable package; the user can prune.
- **Non-HTTP services** (Postgres, Redis, gRPC, brokers): allocated as `tcp` type; exposed as `HOSTBIND_<SERVICE>_HOST/PORT`; `internal: true` hides from URL listing.
- **Docker Compose:** v0.2 = detect published-port collisions and report; v0.4 = override-file generation with allocated ports.
- **OAuth / callback URLs:** many providers only accept exact `http://localhost:<port>`. Ports-only v0.1 avoids the `.localhost` problem; `doctor` flags callback URLs mismatching the allocated port. Document a recommended pattern.
- **Uninstall/reset:** `hostbind uninstall` removes skill files, snippets, registry, and shell hooks, and lists everything it touches. `hostbind reset` clears the registry only.
- **What HostBind writes into repos:** only `.hostbind/config.yaml` (optional, commit-friendly, no ports) and a `.gitignore` entry for `.env.hostbind`. Documented in `docs/what-we-write.md`.
- **Testing:** allocator concurrency tests (N processes racing), fake-process integration tests, per-adapter golden tests, CI matrix on macOS/Linux/Windows.
- **Naming:** "mesh" suggests service-mesh networking. Check npm/GitHub/trademark conflicts and domain availability before launch; consider alternatives with a plainer meaning. Decide *before* the first public commit — renames after stars are painful.

## 15. Install

```bash
brew install hostbind              # macOS/Linux
curl -fsSL https://.../install.sh | sh
scoop install hostbind             # Windows
npm i -g hostbind                  # thin wrapper that downloads the binary
go install …                      # for Go users
```

Then `hostbind install-skill` detects Claude Code, Codex, OpenCode, Kilo (and later Cursor, Windsurf, Gemini CLI, Copilot) and installs the right files, showing exactly what it will write first.

## 16. Repo structure

```
hostbind/
├── README.md  LICENSE  CONTRIBUTING.md  CODE_OF_CONDUCT.md
├── SECURITY.md  CHANGELOG.md  GOVERNANCE.md  ROADMAP.md
├── docs/  (concepts, architecture, alternatives, adapters-guide,
│           agent-integration, troubleshooting, what-we-write)
├── cmd/hostbind/          # CLI entry
├── internal/
│   ├── registry/  allocator/  runner/  scanner/  doctor/  envgen/  mcp/
├── adapters/             # declarative YAML + optional Go
│   ├── vite.yaml  next.yaml  fastapi.yaml  express.yaml …
├── skills/hostbind-local-development/
│   ├── SKILL.md  references/
├── integrations/         # per-agent installers (claude, codex, opencode, kilo)
├── bench/                # agent-compliance benchmark
├── examples/             # vite+fastapi, next+express, monorepo, worktrees
├── tests/
└── .github/  (issue templates, PR template, CI, release workflow, FUNDING)
```

## 17. Open-source growth plan

**Make it worth starring — the first 60 seconds decide everything.**

1. **README with a 20-second demo GIF**: two projects that both want `:3000` → `hostbind run` on both → `hostbind context --json` → agent uses the right URL. Show the *before* pain, then the fix.
2. **One-command quickstart** that works in <2 minutes on a fresh machine.
3. **Badges, comparison table** (honest), clear "Who is this for / not for".
4. **Ship v0.1 small but rock-solid.** A tool that never breaks beats a tool with 50 features.

**Make it easy to contribute:**

- `good first issue` labels seeded up front (each new framework adapter = one issue).
- **Adapter contribution guide**: "add a framework in 30 minutes" with a template and golden test.
- `CONTRIBUTING.md` with local setup in ≤3 commands, `make test`, `make lint`.
- Fast PR review SLA (e.g., 48h first response), friendly `CODE_OF_CONDUCT.md`, `GOVERNANCE.md` (how maintainers are added), public `ROADMAP.md`, GitHub Discussions on.
- Credit contributors in release notes; all-contributors bot.
- Conventional commits, automated releases (goreleaser), Dependabot.

**Distribution:**

- Launch: Hacker News (Show HN), Reddit (r/webdev, r/programming, r/ClaudeAI, r/LocalLLaMA-adjacent communities), X/LinkedIn, dev.to, Product Hunt.
- **Content that spreads:** "Why your AI agent keeps guessing localhost:5173" blog post with benchmark numbers; short screen recordings of 3 agents in parallel without collisions.
- Submit to awesome-lists (awesome-cli-apps, awesome-claude-code and similar agent-tool lists), skill/MCP directories, Homebrew, Scoop, npm.
- Integrations as growth: PRs/docs for each agent ecosystem; a ready `AGENTS.md` template people can copy even without installing.
- Dogfood publicly: use it on your own projects (FTM CRM, LOKNOMIC, Solar Mapper) and show real logs.

**Sustainability:** GitHub Sponsors / Open Collective later; keep core free forever; avoid relicensing surprises. Any future paid team features stay separate from the local core.

## 18. Roadmap

| Version | Goal | Contents |
|---|---|---|
| **v0.1** (2 wks) | Prove agents stop guessing ports | `run`, `ls`, `context --json`, `port`, `stop`, `gc`, generic+vite+next+fastapi+express adapters, skill + AGENTS.md snippet, macOS/Linux |
| **v0.2** | Zero-config | `scan/init/up`, config file, `doctor` (read-only), `.env.hostbind`, monorepo scan, compose collision report |
| **v0.3** | Agent-native | MCP server, `fix`/`undo`, worktrees, `logs`, Windows, more agent installers |
| **v0.4** | Stable URLs | daemon, reverse proxy + `*.localhost`, health checks, dependencies, dashboard |
| **v1.0** | Trust | Stable JSON/CLI/adapter API guarantees, security review, docs complete |

## 19. Validation plan (do this first)

1. Use the alternatives for a day each; write down what's missing.
2. Post the problem (not the product) in 3–4 communities; count who describes it unprompted.
3. Track your own pain: how many times/week does a port/URL mismatch cost you time?
4. Build the v0.1 spike; give it to 5–10 developers; watch how agents behave (bench).
5. **Kill/pivot criteria:** if agents ignore HostBind even with all three layers, or existing tools already cover ≥80% of the value, narrow to the piece that's genuinely unique (likely the agent-context/doctor layer).

## 20. Success metrics

Technical: port conflicts prevented · agent-compliance rate (bench) · allocation reliability under concurrency · time-to-first-`up` · detection false-positive rate · doctor issues correctly diagnosed.
Community: weekly active installs (from package-manager stats, no telemetry) · adapters contributed by non-maintainers · issue response time · repeat contributors.

## 21. Non-goals

Cloud/dashboard SaaS, accounts, billing, Kubernetes, replacing Docker Compose, IDE, team features — not until the local core is proven.

## 22. Final definition

> **HostBind is a local-first, single-binary service registry and port allocator for multi-project development. It gives every service a persistent, conflict-free endpoint, keeps environment configuration in sync, diagnoses stale setups, and gives AI coding agents reliable structured context — so agents discover the environment instead of guessing it.**

**README pitch (3 lines):**
Run many projects, worktrees, and AI agents on one machine without port chaos.
`hostbind run` → conflict-free ports. `hostbind context` → agents know where everything is. `hostbind doctor` → know why it broke.
Local-first. One binary. No account.
