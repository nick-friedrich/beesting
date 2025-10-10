# 🐝 Beesting

A simple CLI tool to manage Go applications in a monorepo structure.

## Installation

Install directly with Go:

```bash
go install github.com/nick-friedrich/beesting/cmd/beesting@latest
```

After installation, you can use `beesting` from anywhere.
If it don't work check your $GOPATH and make sure you have the `bin` directory in your path. Also try to source your `.zshrc` or `.bashrc` file or restart your terminal.

```bash
beesting new my-app
beesting dev my-app
```

## Commands

### `new` - Create a new application

Create a new application in the `app/` directory with a basic HTTP server template.

```bash
beesting new <app-name>
```

**Example:**

```bash
beesting new my-api
```

This creates:

```
app/
  my-api/
    main.go  # Simple HTTP server on :8080
```

### `dev` - Run an application in development mode

Start an application in development mode with hot-reloading.

```bash
beesting dev <app-name>
```

**Example:**

```bash
beesting dev my-api
# Server starts on http://localhost:8080
```

This automatically uses [Air](https://github.com/air-verse/air) for hot-reloading if installed, otherwise falls back to `go run`.

When you edit and save files, the server automatically rebuilds and restarts!

**Install Air for hot-reloading:**

```bash
go install github.com/air-verse/air@latest
```

---

### Local Development

If you're working on the beesting project itself, you can use the Makefile:

```bash
make new my-app
make dev my-app
```
