# Aranea

A concurrent web crawler written in Go. Workers pull URLs from a Redis queue, fetch pages, extract links, and enqueue unseen URLs while exporting Prometheus metrics and visualizing them on Grafana Cloud.

## Tech stack

| Piece | What it does |
| --- | --- |
| **Go** | Crawler runtime (`cmd/crawler`) |
| **goquery** | HTML parsing / link extraction |
| **Redis** | URL queue (`crawler_queue`) and visited set (`visited_urls`) |
| **Prometheus** | Scrapes crawler metrics and can remote-write them |
| **Grafana Cloud** | Optional dashboard for those metrics |
| **Docker Compose** | Runs Redis and Prometheus |

Each start clears the Redis queue and visited set, then re-seeds from the URL you pass in. Worker goroutines pull URLs, fetch pages, normalize discovered links, and enqueue ones that have not been seen.

Metrics served at `:2112/metrics`:

- `aranea_pages_crawled_total`
- `aranea_fetch_errors_total`
- `aranea_queue_depth`
- `aranea_fetch_duration_seconds`

## Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) (for Redis and Prometheus)

## Quick start

```bash
git clone https://github.com/dracofs/Aranea.git
cd Aranea
```

### 1. Start Redis (required)

```bash
docker compose up -d redis
```

Redis listens on `localhost:6379`. Override with `REDIS_ADDR` or `REDIS_PASSWORD` if needed.

### 2. Run the crawler

Pass a starter URL with `-seed` (or as the first argument):

```bash
go run ./cmd/crawler -seed https://redis.io/
```

```bash
go run ./cmd/crawler https://en.wikipedia.org/wiki/Web_crawler
```

Optional: `-workers` (default `5`).

```bash
go run ./cmd/crawler -seed https://redis.io/ -workers 8
```

You should see workers logging URLs. The seed must be an absolute `http://` or `https://` URL.

### 3. Metrics (optional)

Raw metrics while the crawler is running:

```bash
curl http://localhost:2112/metrics
```

Local Prometheus UI: [http://localhost:9090](http://localhost:9090) (check **Status → Targets**; `aranea` should be **UP**).

To scrape and forward metrics, copy `.env.example` to `.env` and fill in Grafana Cloud credentials from your stack’s Prometheus details:

```bash
cp .env.example .env
```

```env
GRAFANA_CLOUD_USERNAME=
GRAFANA_CLOUD_PASSWORD=
```

Then:

```bash
docker compose up -d prometheus
```

Prometheus scrapes `host.docker.internal:2112` (the crawler on your machine) and remote-writes `aranea_*` series to Grafana Cloud. Query them in **Grafana Cloud Explore**, for example:

```promql
aranea_pages_crawled_total
```

If you skip Grafana Cloud, leave `.env` empty and you can still use Redis plus `curl` on `:2112`. Prometheus will fail remote_write without valid credentials.

## Layout

```
cmd/crawler/          entrypoint
internal/crawler/     crawl loop
internal/fetcher/     HTTP fetch
internal/parser/      link extraction
internal/queue/       Redis list
internal/dedupe/      Redis visited set
internal/metrics/     Prometheus instrumentation
internal/utils/       URL normalization
docker-compose.yml    Redis + Prometheus
prometheus.yaml       scrape + remote_write
```

## License

MIT. See [LICENSE](LICENSE).
