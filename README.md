# Gator

Gator is a command-line blog aggregator that allows users to add and follow RSS feed URLs, creating a custom feed of posts from various sources.

## Requirements

To use Gator, ensure the following are installed:

- **PostgreSQL** - Gator uses a PostgreSQL database to store user and feed information.
- **Go** - Gator is built with Go, so it must be installed on your system.

## Installation

Install Gator by running:

```bash
go install github.com/Nukambe/gator@latest
```

This will install the `gator` CLI in your Go bin directory.

## Configuration

Gator requires a configuration file to connect to your database. Create a file named `.gatorconfig.json` in your home directory with this structure:

```json
{
  "db_url": "postgres://username:password@localhost:5432/database_name?sslmode=disable",
  "current_user_name": "null"
}
```

Replace `username`, `password`, and `database_name` with your PostgreSQL credentials.

## Usage

Run Gator with:

```bash
gator
```

### Commands

Gator provides several commands to manage feeds and user accounts. Here’s a summary:

- **User Management**
    - `login` - Log into an existing account.
    - `register` - Register a new account.
    - `reset` - Reset the database, deleting all users.
    - `users` - View a list of registered users.

- **Feed Aggregation**
    - `agg` - Aggregates all posts from followed feeds.

- **Feed Management**
    - `addfeed <url>` - Adds a new RSS feed URL to follow.
    - `feeds` - Lists all available feeds.
    - `follow <feed_url>` - Follow a feed by its URL.
    - `following` - View feeds you’re currently following.
    - `unfollow <feed_url>` - Unfollow a feed by its URL.

- **Browsing**
    - `browse` - Browse posts from followed feeds.