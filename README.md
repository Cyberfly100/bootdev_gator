# Gator

Gator is a command-line blog aggregator built with Go and PostgreSQL built as part of the boot.dev backend dev course. It allows users to register, follow RSS feeds, aggregate posts, and browse their personalized feed.

## Installation (developers)

### Prerequisites

#### Go
Install using webi: `curl -sS https://webi.sh/golang | sh`

#### PostgreSQL
- **Linux**: `sudo apt update && sudo apt install postgresql postgresql-contrib`
- **macOS**: `brew install postgresql`
- **Windows**: Download and install from [postgresql.org](https://www.postgresql.org/download/)

#### Goose (Database Migration Tool)
- Install with: `go install github.com/pressly/goose/v3/cmd/goose@latest`

## Usage

### Setup

1. Start PostgreSQL: `./startDB.sh` (or `sudo systemctl start postgresql` on Linux)
2. Create a database and user in PostgreSQL.
3. Run database migrations using goose: `cd sql/schema && goose postgres "<db_url>" up` (replace `<db_url>` with your database connection string)
4. Build the project: `go build -o gator`

Once the project is built, the Go toolchain is no longer needed for running the application. However, PostgreSQL must remain installed and running, as the application requires a database connection to function.

### Configuration

Create a config file at `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://username:password@localhost:5432/dbname?sslmode=disable",
  "current_user_name": ""
}
```

Replace `username`, `password`, and `dbname` with your PostgreSQL credentials.

### Running the Program

Run commands with: `./gator <command> [args]`

Available commands:
- `register <username>`: Register a new user
- `login <username>`: Log in as an existing user
- `reset`: Delete all users from the users db (for testing)
- `users`: List all users (highlight current user)
- `agg <time>`: Aggregate posts from followed feeds (e.g., `agg 30s`)
- `addfeed <name> <url>`: Add a new RSS feed (requires login)
- `feeds`: List all feeds
- `follow <url>`: Follow a feed (requires login)
- `following`: List followed feeds for current user (requires login)
- `unfollow <url>`: Unfollow a feed (requires login)
- `browse [limit]`: Browse posts from followed feeds (default limit 2, requires login)