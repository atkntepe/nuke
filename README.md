# nuke

A minimal CLI tool to list, select, and kill processes running on development ports, with optional cache clearing for common package managers and build tools.

**macOS only** — uses `lsof` which is available by default.

## Install

```bash
go install github.com/atkntepe/nuke@latest
```

Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is in your `PATH`.

## Usage

### Interactive mode

```bash
nuke
```

Scans dev ports (3000–9999), shows an interactive list. Use arrow keys to navigate, `Space` to select, `Enter` to kill, `q`/`Esc` to quit, `a` to select all.

```
  [ ] 3000   node          vite dev server
  [x] 4000   python        flask run
  [ ] 8080   ruby          rails server
```

### Kill a specific port

```bash
nuke 3000
nuke 8080 --force
```

### List active ports

```bash
nuke list
nuke list --all
```

Non-interactive, pipe-friendly output:

```
PORT   PID     NAME      COMMAND
3000   48291   node      vite dev server
4000   51033   python    flask run
```

### Clear dev caches

```bash
nuke cache
nuke cache --all
```

Detects installed tools and lets you pick which caches to clear:

| Target | Method |
|--------|--------|
| npm | `npm cache clean --force` |
| yarn | `yarn cache clean` |
| pnpm | `pnpm store prune` |
| bun | `bun pm cache rm` |
| vite | delete `node_modules/.vite` |
| next | delete `.next/cache` |
| turbo | delete `.turbo` |
| go | `go clean -cache` |
| gradle | delete `~/.gradle/caches` |
| maven | delete `~/.m2/repository` |
| docker | `docker system prune -f` |

## Flags

| Flag | Description |
|------|-------------|
| `--all` | Scan all ports (1024–65535) instead of 3000–9999, or clear all caches without prompting |
| `--force` | Use SIGKILL instead of SIGTERM |
