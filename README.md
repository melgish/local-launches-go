# local-launches-go

A rewrite of my website scraper in Go. Periodically scrapes
[https://spaceflightnow.com/launch-schedule](https://spaceflightnow.com/launch-schedule)
and filters the list to what's happening locally.

# run locally

```sh
go run .
```

# build and run

To build the project locally:

```
go build -o local-launches
./local-launches
```

# Docker

To build and run with Docker:

```
docker build -t local-launches .
docker run -p 8080:8080 local-launches
```

# Configuration

You can configure the server using the following environment variables:

- `HTTP_PORT`: Port for the HTTP server (default: `8080`).
- `LOG_REQUESTS`: Set to `true` to log incoming HTTP requests (default:
  `false`).
- `REFRESH`: Scrape interval, e.g. `4h`, `30m` (default: `4h`).

The site periodically scrapes the Spaceflight Now launch schedule and displays
upcoming local launches.

# License

MIT License.
