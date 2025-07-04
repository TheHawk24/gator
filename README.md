# Gator (RSS Feed Aggregator)

# Reuirements
* Postgres
* Go

# How to Install
```
go install github.com/TheHawk24/gator@latest
```

# Configuration
1. Create a `~/.gatorconfig.json` file
2. This json file should contain this fields, "db_url" and "current_user_name". db_url value should be url to connect to database. The current_user_name value is configurable
```
{
    "db_url": <url>,
    "current_user_name: <name>
}
```
3. Save the file

# How to use
- **register** - Register a user, expects a single argument. `gator register <name>`
- **login** - Login with a user, expects a single argument. `gator login <name>`
- **reset** - Delete all users, expects no arguments. `gator reset`
- **users** - List all users, expects no arguments. `gator users`
- **agg** - Collect feeds every interval. `gator agg <time>`. Example ("1s", "10m" "2h")
- **addfeed** - Add a feed to database, expects two arguments. `gator addfeed <feed name> <url>`
- **feeds** - Display all feeds, expects no arguments. `gator feeds`
- **follow** - Follow a feed, expects a single argument. `gator follow <url>`
- **folloing** - Display all feeds the user is following, expects no arguments. `gator following`
- **unfollow** - Unfollow a feed, expects a single argument. `gator unfollow <url>`
- **browse** - Display posts from all feeds a user follows. expects no arguments. `gator browse`



