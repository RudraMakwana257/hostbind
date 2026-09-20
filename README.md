# HostBind 🚀

[![Go Report Card](https://goreportcard.com/badge/github.com/RudraMakwana257/hostbind)](https://goreportcard.com/report/github.com/RudraMakwana257/hostbind)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Local dev environments that AI agents can actually trust.**

Run 3 projects and 3 AI agents at once and you get: two apps fighting for `:3000`, `VITE_API_URL` pointing at a dead port, CORS errors, and AI agents that blindly *assume* `localhost:5173`.

**HostBind** fixes this. It is a local-first, single-binary service registry and port allocator. 
*Install once. Every project, every worktree, every agent gets its own conflict-free ports.*

> **Core thesis:** Agents should discover the environment, not guess it.

---

## ⚡ Quickstart

### 1. Install
*(Pre-compiled binaries coming soon)*
```bash
go install github.com/RudraMakwana257/hostbind/cmd/hostbind@latest
```

### 2. Run a service
Instead of running `npm run dev`, run:
```bash
hostbind run --name web
```
HostBind will automatically detect your framework (Next, Vite, Generic), allocate a strictly free port (e.g., `4300`), inject it safely, and start the process.

### 3. Let AI Agents know where things are
```bash
hostbind context --json
```
Output:
```json
{
  "project": "my-app",
  "instance": "main",
  "services": {
    "web": "http://localhost:4300"
  }
}
```

---

## 🤖 AI Agent Integration

HostBind is built **agent-first**. To stop your AI coding assistants (Claude, Cursor, Copilot, etc.) from guessing wrong ports:

Run this command in your project root:
```bash
hostbind install-skill
```
This injects a standard `AGENTS.md` file that teaches any AI agent to query `hostbind context` before attempting to interact with localhost.

---

## 🛠 Features (v0.1)

- **Daemonless & Fast:** Powered by a pure-Go SQLite WAL-mode registry (`~/.hostbind/registry.db`). No background daemon required.
- **Adapter Ladder:** Injects ports via `PORT` env vars or CLI arguments (`--port`). It **never** modifies your source code silently.
- **Framework Support:** Out-of-the-box support for Next.js, Vite, and a Generic fallback.
- **Environment Sync:** Auto-generates `.env.hostbind` and updates `.gitignore`.

## 🤝 Contributing

We want adapters for **every framework on earth** (FastAPI, Express, Django, Laravel, Spring, Astro, Nuxt, etc.). 
Writing an adapter takes 10 minutes. See our `internal/adapters/` directory for examples!

1. Fork the repo.
2. Add your framework to `internal/adapters/`.
3. Submit a PR!

## License
MIT
