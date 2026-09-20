---
name: hostbind-local-development
description: Guides the AI agent on how to manage local ports and services using HostBind, preventing port conflicts and guessing.
---

# HostBind Local Development Skill

You are working in a repository that uses **HostBind** for managing local development ports.

**CRITICAL RULE:** NEVER guess, assume, or hardcode development ports (e.g., do not assume Vite is on 5173, Next.js on 3000, or FastAPI on 8000). HostBind assigns ports dynamically to prevent cross-project conflicts.

## How to use HostBind

When the user asks you to start a service, test the API, or debug a frontend, you must follow these steps:

### 1. Check Running Services
Run the following command to see what is already running and get their exact URLs:
```bash
hostbind context --json
```
*Read the JSON output to find the exact `http://localhost:<PORT>` for the service you need to interact with.*

### 2. Start a Service
If the service is not running, start it using HostBind. HostBind will automatically detect the framework (Vite, Next, etc.), assign a free port, and start the process.
```bash
hostbind run --name <service_name>
```
*(e.g., `hostbind run --name web` or `hostbind run --name api`)*

### 3. Stop a Service
If you need to kill a service to clear a port or restart it, use:
```bash
hostbind stop <service_name>
```

### 4. Read Environment Variables
HostBind generates a `.env.hostbind` file in the root of the project. If you need to write a script that connects to the frontend or backend, read the environment variables from this file.

## Troubleshooting

- **Process already running / Port blocked:** Do not run `kill -9` blindly. Run `hostbind ls` to see what HostBind is tracking, and use `hostbind stop <service>` to stop it cleanly.
- **Command fails to start:** Check the framework's native config (like `vite.config.js`). HostBind relies on passing standard arguments (like `--port`) or environment variables (like `PORT`).
