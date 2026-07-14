# Gator

## Environment Setup

Gator requires:

- Go 1.26 or newer
- A running instance of a [PostgreSQL](https://www.postgresql.org/) database

Optionally, you can install [Goose](https://github.com/pressly/goose) and run it
in the [sql/schema](https://github.com/ashy558/bootdev-gator/sql/schema) folder
to set up the `gator` database.

You can install the latest version of `goose` by running:

```shell
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Installation

You can install the latest version of `gator` by running:

```shell
go install github.com/ashy558/main/cmd/gator@latest
```

## Building

Set the `LD_FLAGS` with meta information like the version or the commit:

```shell
export LD_FLAGS="-w -s -X main.Version=$(git describe --tags | cut -c 2-) -X main.BuildDate=$(date "+%F-%T") -X main.Commit=$(git rev-parse --verify HEAD) -X main.Mode=prod";
```

Then build the gator binary like this:

```shell
go build -ldflags="$LD_FLAGS" -o gator
```

## Usage

```help
Usage: gator COMMAND

Examples:
  gator addfeed HackerNews https://hnrss.org/
  gator agg 5m
  gator browse 10
  gator feeds
  gator follow https://hnrss.org/
  gator following
  gator help
  gator login boots
  gator register boots
  gator reset
  gator unfollow https://hnrss.org/
  gator users

Commands:
  addfeed NAME URL  Add a new feed
  agg INTERVAL      Fetch new posts from feeds
  browse [LIMIT]    Browse posts from followed feeds
  feeds             Print the feeds added
  follow URL        Follow an existing feed
  following         Print the feeds followed
  help              Print this help message
  login USERNAME    Login as a registered user
  register USERNAME Register a new user
  reset             Truncate all user data (development only)
  unfollow URL      Unfollow a followed feed
  users             Print all registered users
```
