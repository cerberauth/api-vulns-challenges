# stale-while-revalidate / stale-if-error Not Honored

This challenge demonstrates a cache that advertises `stale-while-revalidate` and `stale-if-error` in its `Cache-Control` header, but never actually honors them: once `max-age` expires, it always synchronously re-fetches from the origin, defeating both the performance benefit (serving stale instantly while revalidating in the background) and the resilience benefit (serving stale instead of surfacing an origin outage).

## How to run it

```bash
go run main.go serve
```

## Endpoints

- `GET /` — returns a timestamped body with `Cache-Control: max-age=2, stale-while-revalidate=5, stale-if-error=5`
- `GET /toggle-outage` — flips the simulated origin between up and down, to exercise `stale-if-error`

## Exploiting it / confirming the finding

```bash
# populate the cache
curl -s http://localhost:8080/

# wait past max-age (2s) but within stale-while-revalidate (5s)
sleep 3
curl -s -D - http://localhost:8080/ -o /dev/null   # expect X-Cache: STALE served instantly (fixed) vs MISS (vulnerable)

# simulate an origin outage and confirm stale-if-error behavior
curl -s http://localhost:8080/toggle-outage
sleep 3
curl -s -D - http://localhost:8080/ -o /dev/null   # expect a stale 200 (fixed) vs a 500 (vulnerable)
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: stale-while-revalidate and stale-if-error are advertised but never honored
go run main.go serve --vulnerable=true

# fixed: stale content is served within the declared windows while revalidating in the background / masking origin errors
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
