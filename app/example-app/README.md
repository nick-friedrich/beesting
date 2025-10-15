# Example API

A Go web application with Tailwind CSS, SQLC, and Goose migrations.

## Setup

### 1. Install Dependencies

```bash
# Go dependencies
go mod tidy

# Node.js dependencies (for Tailwind CSS)
npm install
```

### 2. CSS Building

CSS is automatically built in development mode (when `GO_ENV != "production"`). In production, ensure `static/output.css` is pre-built.

### 3. Run Migrations

```bash
# Migrations run automatically on startup
goose -dir db/migrations sqlite3 ./app.db up
```

### 4. Run the Application

```bash
# Development mode (recommended)
make dev

# Or using npm directly
npm run dev

# Or using root Makefile
cd ../.. && make dev example-api

# Or directly
go run main.go
```

## Development Workflow

1. **All Changes**: Run `make dev` in one terminal
2. **Database Changes**: Create new migrations with `goose create`
3. **CSS Changes**: Edit `input.css` - Tailwind watch rebuilds automatically
4. **Go Changes**: Air hot reloads automatically

## Project Structure

```
app/example-api/
├── db/                    # Database layer
│   ├── migrations/        # Goose migrations
│   ├── queries/          # SQLC queries
│   └── *.go             # Generated SQLC code
├── handler/              # HTTP handlers
├── pkg/web/             # Web utilities (templates)
├── static/              # Static assets
│   └── output.css       # Generated Tailwind CSS
├── input.css           # Tailwind CSS input
├── templates/           # HTML templates
├── main.go             # Application entry point
└── package.json        # Node.js dependencies
```

## Tech Stack

- **Backend**: Go + Chi router
- **Database**: SQLite + SQLC + Goose
- **Frontend**: HTML templates + Tailwind CSS
- **Development**: Air (hot reload) + Tailwind CLI (CSS watch)

## API Endpoints

- `GET /` - Home page
- `GET /health` - Health check
- `GET /posts` - List posts
- `POST /posts` - Create post
- `GET /posts/{id}` - Get post
- `PUT /posts/{id}` - Update post
- `DELETE /posts/{id}` - Delete post
- `POST /posts/{id}/publish` - Publish post
