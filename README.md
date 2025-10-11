## Installation

Install directly with Go:

Make sure you have sqlc and air (air is optional) installed.

```bash
go install github.com/nick-friedrich/beesting/cmd/beesting@latest
```

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

```bash
go install github.com/air-verse/air@latest
```

After installation, you can use `beesting` from anywhere.
If it don't work check your $GOPATH and make sure you have the `bin` directory in your path. Also try to source your `.zshrc` or `.bashrc` file or restart your terminal.

```bash
beesting new my-app
beesting dev my-app
```
