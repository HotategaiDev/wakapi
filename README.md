# 📈 wakapi
**A minimalist, self-hosted WakaTime-compatible backend for coding statistics**

![Wakapi screenshot](https://anchr.io/i/zCVbN.png)

[![Buy me a coffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://buymeacoff.ee/n1try)

## Prerequisites
* Go >= 1.10 (with `$GOPATH` properly set)
* A MySQL database

## Usage
* Create an empty MySQL database
* Get code: `go get github.com/n1try/wakapi`
* Go to project root: `cd "$GOPATH/src/github.com/n1try/wakapi"`
* Install dependencies: `go get -d ./...`
* Copy `.env.example` to `.env` and set database credentials
* Set target port in `config.ini`
* Build executable: `go build`
* Run server: `./wakapi`
* Edit your local `~/.wakatime.cfg` file
  * `api_url = https://your.server:someport/api/heartbeat`
  * `api_key = the_api_key_printed_to_the_console_after_starting_the_server`
* Open [http://localhost:3000](http://localhost:3000) in your browser

### User Accounts
* When starting wakapi for the first time, a default user _**admin**_ with password _**admin**_ is created. The corresponding API key is printed to the console.
* Additional users, at the moment, can be added only via SQL statements on your database, like this:
    * Connect to your database server: `mysql -u yourusername -p -H your.hostname` (alternatively use GUI tools like _MySQL Workbench_)
    * Select your database: `USE yourdatabasename;`
    * Add the new user: `INSERT INTO users (id, password, api_key) VALUES ('your_nickname', MD5('your_password'), '728f084c-85e0-41de-aa2a-b6cc871200c1');` (the latter value should be a random [UUIDv4](https://tools.ietf.org/html/rfc4122), as can be found in your `~/.wakatime.cfg`)

## Best Practices
It is recommended to use wakapi behind a **reverse proxy**, like [Caddy](https://caddyserver.com) or _nginx_ to enable **TLS encryption** (HTTPS).
However, if you want to expose your wakapi instance to the public anyway, you need to set `listen = 0.0.0.0` in `config.ini`

## Todo
* Persisted summaries / aggregations (for performance)
* User sign up and log in
* Additional endpoints for retrieving statistics data
* Enhanced UI
  * Loading spinner
  * Responsiveness
* Support for SQLite database
* Dockerize
* Unit tests

## Important Note
**This is not an alternative to using WakaTime.** It is just a custom, non-commercial, self-hosted application to collect coding statistics using the already existing editor plugins provided by the WakaTime community. It was created for personal use only and with the purpose of keeping the sovereignity of your own data. However, if you like the official product, **please support the authors and buy an official WakaTime subscription!**

## License
GPL-v3 @ [Ferdinand Mütsch](https://muetsch.io)
