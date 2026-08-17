# gator

A command-line RSS feed aggregator. Register a user, add feeds, follow the ones you care about, and leave the aggregator running to pull posts into Postgres. Read them back from your terminal whenever you want.

Built while working through the [boot.dev](https://boot.dev) Go course.

## Contents

- [Requirements](#requirements)
- [Setup](#setup)
- [Usage](#usage)
- [Commands](#commands)
- [Tech stack](#tech-stack)
- [Project layout](#project-layout)
- [Development](#development)
- [Notes & warnings](#notes--warnings)

## Requirements

- **Go** 1.26.5+ — to install the CLI.
- **PostgreSQL** 15+ — where users, feeds, and posts live.
- **[goose](https://github.com/pressly/goose)** — to run the migrations.

## Setup

Install the binary:

```bash
go install github.com/its-me-sv/gator@latest
```

If the `gator` command isn't found, add Go's bin directory to your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Create the database and apply the schema (migrations live in [sql/schema/](sql/schema/), so run this from a clone):

```bash
createdb gator
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

Then write `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Leave `current_user_name` empty — `gator` rewrites it whenever you `register` or `login`.

## Usage

```bash
gator register lane
gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml
```

Start the aggregator in one terminal, browse in another:

```bash
gator agg 1m      # runs until Ctrl+C
gator browse 5
```

## Commands

| Command | Arguments | What it does |
| --- | --- | --- |
| `register` | `<name>` | Creates a user and logs in as them. |
| `login` | `<name>` | Switches to an existing user. |
| `users` | — | Lists all users, marking the current one. |
| `reset` | — | Deletes every user, and their feeds, follows, and posts. |
| `addfeed` | `<name> <url>` | Adds a feed and follows it. Needs a logged-in user. |
| `feeds` | — | Lists every feed with the user who added it. |
| `follow` | `<url>` | Follows an existing feed. Needs a logged-in user. |
| `following` | — | Lists the feeds you follow. |
| `unfollow` | `<url>` | Unfollows a feed. Needs a logged-in user. |
| `agg` | `<time_between_reqs>` | Fetches feeds on a loop, waiting the given duration between requests (`10s`, `1m`, `1h`). |
| `browse` | `[limit]` | Prints the latest posts from feeds you follow. Defaults to 2. |

## Tech stack

- **Go** — the CLI itself, standard library only for HTTP and XML parsing.
- **PostgreSQL** — storage.
- **[sqlc](https://sqlc.dev)** — generates type-safe Go from the SQL in
  [sql/queries/](sql/queries/).
- **[goose](https://github.com/pressly/goose)** — schema migrations.
- **[lib/pq](https://github.com/lib/pq)** — Postgres driver.
- **[google/uuid](https://github.com/google/uuid)** — primary keys.

## Project layout

```text
main.go                  wiring, command registration, auth middleware
commands.go              command registry
handler_*.go             one file per command group (users, feeds, follows, posts, agg)
rss_feed.go              fetching and parsing RSS
printers.go              terminal output formatting
internal/config/         reads and writes ~/.gatorconfig.json
internal/database/       sqlc-generated code — don't edit by hand
sql/schema/              goose migrations
sql/queries/             sqlc query source
```

## Development

Run straight from a clone instead of installing:

```bash
go run . users
```

After changing anything in `sql/queries/` or `sql/schema/`, regenerate the database layer:

```bash
sqlc generate
```

New migration files go in `sql/schema/` with a `NNN_name.sql` prefix and the usual goose `-- +goose Up` / `-- +goose Down` markers.

## Notes & warnings

- **`reset` really does delete everything.** No confirmation prompt, no undo. It's there for development.
- **`agg` is a long-running process.** Every tick it scrapes the least recently fetched feed. Use minutes, not seconds — hammering someone's server will get you blocked, and rightly so.
- **There's no auth.** `login` just writes a name to your config file, so any user can act as any other. Fine locally, not fine anywhere else.
- **The connection string sits in plaintext** in `~/.gatorconfig.json`. Don't point it at anything you care about.
