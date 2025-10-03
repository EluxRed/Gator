# Gator
RSS feed aggregator: a CLI local tool

Postgres and Go are required to install and run the program.

1. You can install Gator in your environment with:
```bash
go install github.com/EluxRed/Gator@latest
```

2. Create a config file at ~/.gatorconfig.json (in your home directory) with your own username and password for Postgres and you can set the current_user_name as you wish:
```bash
// json
{
  "db_url": "postgres://USER:PASS@localhost:5432/gator?sslmode=disable",
  "current_user_name": "EluxRed"
}
```

3. Start the Postgres service in the background (needed to access the database):
  for Linux/WSL (Debian)
```bash
sudo service postgresql start
```
  for macOS (Homebrew)
```bash
brew services start postgresql@16
```
  for Windows
```bash
net start postgresql-x64-16
```

4. Prepare the database (repo cloning and goose installation are required):
```bash
git clone https://github.com/EluxRed/Gator
cd Gator
go install github.com/pressly/goose/v3/cmd/goose@latest
createdb gator
goose -dir sql/schema postgres "postgres://USER:PASS@localhost:5432/gator?sslmode=disable" up
```

5. At this point you don't need the repo anymore and you can run the binary with "gator" followed by a command and possibly parameters (if needed).
List of commands:
- **login**: it expects one parameter, the username. It logins with the provided username as current user, but only if the user already exists in the database, saving it in the config file. All other commands will be executed with this user's username. Example:
```bash
gator login EluxRed
```
- **register**: it expects one parameter, the username. It registers a new user into a database. Example:
```bash
gator register EluxRed
```
- **users**: it expects no parameters. It shows all the users registered so far. Example:
```bash
gator users
```
- **addfeed**: it expects two parameters, the feed's name and its URL. It registers the feed into the database and the current user starts following the feed. Example:
```bash
gator addfeed TechCrunch https://techcrunch.com/feed/
```
- **feeds**: it expects no parameters. It shows all the registered feeds. Example:
```bash
gator feeds
```
- **follow**: it expects one parameter, the feed's URL. It lets the current user follow the feed, if it's already registered into the database. Example:
```bash
gator follow https://blog.boot.dev/index.xml
```
- **following**: it expects no parameters. It shows all the feeds that the current user is following. Example:
```bash
gator following
```
- **unfollow**: it expects one parameter, the feed's URL. It lets the current user unfollow the feed, if the user is already following it. Example:
```bash
gator unfollow https://news.ycombinator.com/rss
```
- **agg**: it expects one parameter, the time interval between requests, like 30s or 1m (please do not DOS the services with too many requests). It starts a continuous loop that continuously scrapes posts from the feeds currently present in the database and shows a summary of those posts. This is thought to be run in the background, while you interact with the program in a separate terminal window. Example:
```bash
gator agg 30s
```
- **browse**: it expects one or none parameters, the limit for the number of posts to be shown, which is defaulted to 2 if no parameter is provided. It shows a number of the latest scraped posts, set by the limit provided. Example:
```bash
gator browse 10
```
- **reset**: expects no parameters. Eliminates all the rows in the database, if you want a fresh start. Not reversible. Example:
```bash
gator reset
```
