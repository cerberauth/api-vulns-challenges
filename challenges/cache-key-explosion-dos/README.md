# Cache-Key Explosion / Cache-Busting DoS

This challenge demonstrates cache-storage exhaustion via unbounded query parameters: `/search?q=<anything>` lets the caller pick an arbitrary `q` value, and the vulnerable cache keys every distinct value as its own entry with no limit — so an attacker can generate unbounded cache entries (storage exhaustion) and force an origin request for every single unique value at volume (origin-hammering DoS), defeating the cache's purpose entirely.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /search?q=<value>` — returns a JSON body and reports `X-Distinct-Cache-Entries` plus a running `origin_requests_so_far` counter

## Exploiting it

```bash
for i in $(seq 1 50); do
  curl -s -D - "http://localhost:8080/search?q=$i" -o /dev/null | grep -E "X-Distinct-Cache-Entries|X-Cache"
done
# every request is a MISS and the distinct-entry / origin-request counters grow without bound
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: every distinct "q" value gets its own, unbounded cache entry
go run main.go serve --vulnerable=true

# fixed: the number of distinct cache entries per path is capped; overflow values share a single bucket
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
