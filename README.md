# 🐝 Beesting

A simple CLI tool to manage Go applications in a monorepo structure.

## Installation

Install directly with Go:

```bash
go install github.com/nick-friedrich/beesting/cmd/beesting@latest
```

After installation, you can use `beesting` from anywhere:

```bash
beesting new my-app
beesting dev my-app
```

## Commands

### `new` - Create a new application

Create a new application in the `app/` directory with a basic `main.go` template.

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
    main.go
```

### `dev` - Run an application in development mode

Start an application in development mode by running its `main.go` file.

```bash
beesting dev <app-name>
```

**Example:**

```bash
beesting dev my-api
```

This runs the application using `go run`.

---

### Local Development

If you're working on the beesting project itself, you can use the Makefile:

```bash
make new my-app
make dev my-app
```
