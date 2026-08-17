# gator

`gator` is a command line RSS feed aggregator. You register users, add RSS feeds,
follow the feeds you care about, and let a long-running aggregator collect posts
into Postgres so you can browse them from your terminal.

## Requirements

You need two things installed before you can run `gator`:

- **Go** (1.26.5 or newer) — needed to install the CLI.
- **PostgreSQL** (v15 or newer) — `gator` stores users, feeds, and posts there.

Make sure your Postgres server is running and that you can connect to it, for
example:

```bash
psql "postgres://postgres:postgres@localhost:5432/gator"
```

## Install

Install the CLI with `go install`:

```bash
go install github.com/its-me-sv/gator@latest
```

This compiles a standalone `gator` binary into your `$GOBIN` (usually
`~/go/bin`). Go programs are statically compiled binaries, so once it's
installed you can run `gator` anywhere — the Go toolchain isn't needed at
runtime.

If `gator` isn't found afterwards, add your Go bin directory to your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

> `go run .` from a clone of this repo works too, but that's only for
> development. `gator` is the thing you actually ship and run.

## Database setup

Create a database for the app:

```bash
createdb gator
```

The schema lives in [sql/schema/](sql/schema/) as [goose](https://github.com/pressly/goose)
migrations. From a clone of the repo, apply them with:

```bash
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

## Config file

`gator` reads a JSON config file named `.gatorconfig.json` from your home
directory. Create it yourself:

```bash
touch ~/.gatorconfig.json
```

And put your connection string in it:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

- `db_url` — the Postgres connection string `gator` connects to.
- `current_user_name` — the logged-in user. Leave it empty; `gator` rewrites
  this field for you when you `register` or `login`.

## Usage

Every command is run as `gator <command> [args]`.

Register yourself and add a feed:

```bash
gator register lane
gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml
```

Then start the aggregator in one terminal and browse in another:

```bash
gator agg 1m      # runs until you stop it with Ctrl+C
gator browse 5
```

### Commands

| Command | Arguments | What it does |
| --- | --- | --- |
| `register` | `<name>` | Creates a new user and logs in as them. |
| `login` | `<name>` | Switches the current user to an existing user. |
| `users` | — | Lists all users, marking the current one. |
| `reset` | — | Deletes every user (and their feeds, follows, and posts). |
| `addfeed` | `<name> <url>` | Adds a feed and follows it. Requires a logged-in user. |
| `feeds` | — | Lists every feed, with the user who added it. |
| `follow` | `<url>` | Follows an existing feed. Requires a logged-in user. |
| `following` | — | Lists the feeds the current user follows. |
| `unfollow` | `<url>` | Unfollows a feed. Requires a logged-in user. |
| `agg` | `<time_between_reqs>` | Continuously fetches feeds, waiting the given duration between requests (e.g. `10s`, `1m`, `1h`). |
| `browse` | `[limit]` | Prints the latest posts from feeds you follow (default limit: 2). |

`agg` is a long-running process — it keeps scraping the least recently fetched
feed on every tick until you stop it with `Ctrl+C`. Be polite to the sites you
scrape and use an interval measured in minutes rather than seconds.
