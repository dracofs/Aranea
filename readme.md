# Aranea

A concurrent web crawler written in Go. Workers pull URLs from a Redis queue, fetch pages, extract links, and enqueue unseen URLs while exporting Prometheus metrics.

## Tech stack

| Layer | Choice |
| --- | --- |
| Language | [Go](https://go.dev/) |
| Queue & dedupe | [Redis](https://redis.io/) via [go-redis](https://github.com/redis/go-redis) |
| HTML parsing | [goquery](https://github.com/PuerkitoBio/goquery) |
| Metrics | [Prometheus](https://prometheus.io/) (`client_golang`) |
| Infra | Docker Compose (Redis) |

## Prerequisites

- Go 1.22+ (module targets Go 1.26)
- Docker & Docker Compose (for Redis)

## Quick start

### 1. Start Redis

```bash
docker compose up -d
```

Redis listens on `localhost:6379` by default.

### 2. Run the crawler

```bash
go run ./cmd/crawler
```

By default this seeds `https://redis.io/` and starts **5** workers (configured in `cmd/crawler/main.go`).

### 3. Watch metrics

The crawler serves Prometheus metrics at:

```text
http://localhost:2112/metrics
```

Useful series:

- `aranea_pages_crawled_total`
- `aranea_fetch_errors_total`
- `aranea_queue_depth`
- `aranea_fetch_duration_seconds`

To scrape with Prometheus, point it at [`prometheus.yaml`](./prometheus.yaml).

## Configuration

| Variable | Default | Description |
| --- | --- | --- |
| `REDIS_ADDR` | `localhost:6379` | Redis host:port |
| `REDIS_PASSWORD` | _(empty)_ | Redis password, if any |

Example:

```bash
REDIS_ADDR=localhost:6379 go run ./cmd/crawler
```

## Project layout

```text
cmd/crawler/          # entrypoint
internal/crawler/     # crawl loop & worker orchestration
internal/fetcher/     # HTTP fetch
internal/parser/      # link extraction
internal/queue/       # Redis URL queue
internal/dedupe/      # Redis visited-set
internal/metrics/     # Prometheus metrics
internal/redisclient/ # Redis connection helper
internal/utils/       # URL normalization
```

## License

See [LICENSE](./LICENSE).
