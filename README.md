# Gator

## Introduction

This project was realized by following the boot.dev
[Build a Blog Aggregator in Go](https://www.boot.dev/courses/build-blog-aggregator-golang)
course.

It's an RSS feed aggregator CLI tool, that allows users to:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the
  full post

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

## Ideas to implement

- Add sorting and filtering options to the browse command
- Add pagination to the browse command
- Add concurrency to the agg command so that it can fetch more frequently
- Add a search command that allows for fuzzy searching of posts
- Add bookmarking or liking posts
- Add a TUI that allows you to select a post in the terminal and view it in a
  more readable format (either in the terminal or open in a browser)
- Add an HTTP API (and authentication/authorization) that allows other users to
  interact with the service remotely
- Write a service manager that keeps the agg command running in the background
  and restarts it if it crashes
