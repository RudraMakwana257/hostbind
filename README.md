# HostBind 🚀

[![NPM Version](https://img.shields.io/npm/v/hostbind)](https://www.npmjs.com/package/hostbind)
[![Go Report Card](https://goreportcard.com/badge/github.com/RudraMakwana257/hostbind)](https://goreportcard.com/report/github.com/RudraMakwana257/hostbind)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Local dev environments that AI agents can actually trust.**

Run 3 projects and 3 AI agents at once and you get: two apps fighting for `:3000`, `VITE_API_URL` pointing at a dead port, CORS errors, and AI agents that blindly *assume* `localhost:5173`.

**HostBind** fixes this. It is a local-first, single-binary service registry and port allocator. 
*Install once. Every project, every worktree, every agent gets its own conflict-free ports.*

> **Core thesis:** Agents should discover the environment, not guess it.

---

## ⚡ Quickstart

### 1. Install (Global)
The easiest way to install (works on Windows, Mac, and Linux):
```bash
npm install -g hostbind
```
*(No Node? You can also download the binary directly from GitHub Releases or run `go install github.com/RudraMakwana257/hostbind/cmd/hostbind@latest`)*

### 2. Run a service
Instead of running `npm run dev` or `python main.py`, run:
```bash
hostbind run --name web
```
HostBind will automatically detect your framework (Next, Vite, Python, Generic), allocate a strictly free port, inject it safely, and start the process.

### 3. See what's running
```bash
hostbind ls
```
*(Prints a clean table of all running projects, PIDs, and ports).*

---

## 🤖 AI Agent Integration (MCP & Skills)

HostBind is built **agent-first**. To stop your AI coding assistants (Claude, Cursor, Copilot, etc.) from guessing wrong ports:

Run this command in your project root:
```bash
hostbind install-skill
```
This injects a standard `AGENTS.md` file that teaches any AI agent to query `hostbind context` before attempting to interact with localhost. HostBind also runs an internal **MCP (Model Context Protocol) Server** for deep native integration with Claude and Cursor.

---

## 🛠 Features

- **Daemonless Registry:** Powered by a pure-Go SQLite WAL-mode registry. No background daemon required.
- **Reverse Proxy:** Run `hostbind proxy --port 8080` to access services beautifully via `http://api.my-app.localhost:8080`.
- **Doctor Diagnostics:** Run `hostbind doctor` to automatically find and fix orphaned ports or zombie PIDs.
- **Security Check:** Prevents AI agents from injecting dangerous shell commands via interactive terminal prompts.
- **Adapter Ladder:** Out-of-the-box support for Next.js, Vite, Python (FastAPI/Flask), and a Generic fallback. 

## 🤝 Contributing

We want adapters for **every framework on earth** (Express, Django, Laravel, Spring, Astro, Nuxt, etc.). 
Writing an adapter takes 10 minutes. See our `internal/adapters/` directory for examples!

1. Fork the repo.
2. Add your framework to `internal/adapters/`.
3. Submit a PR!

## License
MIT
