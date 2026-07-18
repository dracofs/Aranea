# Aranea

Aranea is an experimental distributed web crawler written in Go. Multiple
processes can share a Redis-backed frontier, claim pages without duplicating
work, and discover links concurrently.

## Current milestone

- Atomic URL deduplication and queue insertion in Redis.
- Reliable queue semantics: claimed pages remain in an in-flight list until a
  worker acknowledges them.
- Configurable worker count, crawl depth, request timeout, and same-host scope.
- URL resolution, fragment removal, and HTTP(S)-only filtering.
- Graceful shutdown on `SIGINT` or `SIGTERM`.
- Focused unit tests for crawl coordination, parsing, fetching, and URLs.

## Run locally

Requirements: Go 1.26+ and Redis 6.2+ (the frontier uses `BLMOVE`).

```sh
go test ./...
go run ./cmd/crawler \
  --seed https://example.com \
  --namespace example-crawl \
  --workers 5 \
  --max-depth 2 \
  --reset
```

`--reset` deletes only the three Redis keys under the selected namespace. Do
not pass it when joining an active crawl. To add another worker process, run the
same command with the same namespace and omit `--reset`.

If a process is interrupted while it owns pages, stop all workers and pass
`--recover` once to return unacknowledged pages to the queue.

Use `go run ./cmd/crawler --help` for all options. The crawler stays on the seed
hostname by default; pass `--same-host=false` to allow external links.

## Roadmap

1. Honor `robots.txt`, add per-host rate limits, and implement retries/backoff.
2. Add Redis integration tests and crash-safe leased claims.
3. Export Prometheus metrics and provide a Grafana dashboard.
4. Store parsed pages and index metadata in PostgreSQL.
5. Containerize the services and deploy worker replicas.
