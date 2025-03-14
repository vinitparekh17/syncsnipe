# SyncSnipe

A local-first file synchronization and backup tool.

## What is SyncSnipe?

SyncSnipe is a single-binary Go application that keeps your files in sync across directories without cloud dependencies. It watches for changes in real-time, syncs them locally, and aims for simplicity, reliability, and user control—all from your machine.

## Features

- [x] **Web UI** – Serve a frontend to manage sync rules and monitor activity (API ready, UI WIP).
- [ ] **Real-Time Local Sync** – Watches and synchronizes specified directories as changes happen (`create`/`write` supported; `remove`/`rename` in progress).
- [ ] **Automated Backups** – Schedule backups to external drives or other locations (planned).
- [ ] **Conflict Detection** – Spots file differences and logs them (resolution coming soon).
- [ ] **Versioning Support** – Track and restore previous file versions (planned).
- [x] **Sync Profiles** – Define multiple sync setups (e.g., "Work," "Pictures," "Music").
- [ ] **Ignore Patterns** – Exclude files per profile (e.g., `.git`, `*.tmp`, `.DS_Store`).

## How to run it locally

This repo uses [Go 1.24 or higher](https://go.dev/dl/), [Node.js LTS](https://nodejs.org/), with [pnpm](https://pnpm.io/)
```bash
$ git clone git@github.com:vinitparekh17/syncsnipe.git
```
if you have already cloned and present in repo's directory prefer using make commands for better experience
```bash
$ make build-frontend # To ensure we have build dir to be served by our golang server
$ go run main.go web  # sub-commands web/cli whichever you prefers
```

> [!NOTE]
> SyncSnipe is under active development; as there are not published artifacts at this time; bugs and inappropriate behaviours are expected.

## Status overview

- **Done**: Profiles, web serving, DB ops.
- **Next**: Full event support (`remove`, `rename`), conflict resolution, testing (unit/integration), backup scheduling, versioning.
