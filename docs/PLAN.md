# HostBind v0.1 Execution Plan

This document outlines the step-by-step plan to build the v0.1 spike of HostBind, strictly adhering to the `HostBind-SPEC.md`. Following this step-by-step approach prevents hallucinations, ensures high code quality, and maintains a clean architecture.

## Phase 1: Project Setup (Go & CLI Framework)
- [ ] Initialize Go module (`go mod init github.com/RudraMakwana257/hostbind`).
- [ ] Set up the standard Go project structure (`cmd/`, `internal/`, `docs/`, `adapters/`).
- [ ] Integrate **Cobra** for robust CLI command parsing and **Viper** for configuration.
- [ ] Create skeleton commands for `run`, `ls`, `context`, `port`, `stop`.

## Phase 2: Core Engine (Registry & Allocator)
- [ ] **Registry (`internal/registry`)**: Implement a local SQLite database (using `mattn/go-sqlite3` or a pure Go alternative like `modernc.org/sqlite` to keep it a simple static binary). Must support WAL mode, file locking, and atomic allocations.
- [ ] **Allocator (`internal/allocator`)**: Logic to find available ports within defined ranges (e.g., 4300-4799 for web) and reserve them via "leases" in the registry.

## Phase 3: The Adapter Ladder (`internal/adapters`)
- [ ] Create the Adapter interface (`detect`, `services`, `portArgs`, `envFor`).
- [ ] Implement declarative/YAML loading for simple adapters.
- [ ] Implement the first set of v0.1 adapters:
  - `generic` (injects `PORT`)
  - `vite` (injects `VITE_PORT`)
  - `next` (injects `PORT`)
  - `fastapi`
  - `express`

## Phase 4: Runner & Environment Sync
- [ ] **Runner (`internal/runner`)**: Logic to spawn child processes, inject the allocated ports via environment variables, and track PIDs.
- [ ] **EnvGen**: Logic to generate and sync the `.env.hostbind` file locally so AI agents and users can read the environment, and auto-add it to `.gitignore`.

## Phase 5: CLI Command Implementation
- [ ] Wire up `hostbind run`: (Detect -> Allocate -> Generate Env -> Run -> Register).
- [ ] Wire up `hostbind ls`: Display running services from the registry.
- [ ] Wire up `hostbind context --json`: Output machine-readable state for AI agents.
- [ ] Wire up `hostbind stop`: Cleanly kill the tracked PID and release the registry lease.

## Phase 6: Testing & Quality Assurance
- [ ] Write Go unit tests for the Allocator (ensuring no race conditions).
- [ ] Write integration tests for the Runner (spawning a fake process).
- [ ] Add linting (`golangci-lint`) and strict error handling (no silent failures).
