# Database Setup with SQLC + Goose

This project uses:

- [sqlc](https://sqlc.dev/) for type-safe SQL queries
- [goose](https://github.com/pressly/goose) for database migrations

## Directory Structure

```
db/
├── migrations/         # Goose migrations (up/down)
│   └── 00001_create_posts.sql
├── queries/           # SQL queries organized by concern
│   └── posts.sql
├── migrate.go         # Migration runner
├── db.go             # SQLC generated
└── models.go         # SQLC generated
```

## Setup

### 1. Install Tools

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 2. Generate SQLC Code

```bash
cd app/example-app
sqlc generate
```

## Migrations with Goose

### Create a New Migration

```bash
cd app/example-app
goose -dir db/migrations create add_comments sql
```

This creates: `db/migrations/YYYYMMDDHHMMSS_add_comments.sql`

### Migration File Format

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    content TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comments;
-- +goose StatementEnd
```

### Run Migrations Manually

```bash
# Migrate up
goose -dir db/migrations sqlite3 ./app.db up

# Migrate down (rollback one)
goose -dir db/migrations sqlite3 ./app.db down

# Check status
goose -dir db/migrations sqlite3 ./app.db status

# Reset database
goose -dir db/migrations sqlite3 ./app.db reset
```

### Migrations in Code

Migrations run automatically on app startup via `db.RunMigrations()`:

```go
database, err := sql.Open("sqlite3", "./app.db")
if err != nil {
    log.Fatal(err)
}

// Runs all pending migrations
if err := db.RunMigrations(database); err != nil {
    log.Fatal(err)
}
```

## SQLC Usage

```go
queries := db.New(database)

// Use generated functions
post, err := queries.CreatePost(ctx, db.CreatePostParams{
    Title:   "Hello World",
    Content: "First post!",
    Author:  "John",
})
```

## Workflow: Adding a New Feature

1. **Create migration:**

   ```bash
   goose -dir db/migrations create create_comments sql
   ```

2. **Write migration SQL** (with up/down)

3. **Create queries file:**

   ```bash
   touch db/queries/comments.sql
   ```

4. **Write SQLC queries**

5. **Generate code:**

   ```bash
   sqlc generate
   ```

6. **Run migrations:**
   ```bash
   # Automatic on app start, or manually:
   goose -dir db/migrations sqlite3 ./app.db up
   ```

## Best Practices

- ✅ Always write both Up and Down migrations
- ✅ One migration file per logical change
- ✅ Use sequential numbering (00001, 00002, etc.)
- ✅ Organize queries by domain (posts.sql, comments.sql)
- ✅ Embed migrations with `//go:embed` for deployment
- ✅ Test rollbacks before production
- ✅ Keep migrations idempotent when possible
