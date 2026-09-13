# WAF / Filtering Behavior Validation

This challenge demonstrates a WAF/filtering layer in front of a search endpoint. The vulnerable server performs no filtering at all, letting a baseline set of SQLi/XSS/path-traversal payloads through unmodified.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /search?q=<payload>` — echoes the query back

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: payloads such as "' OR 1=1--", "<script>alert(1)</script>" or "../../etc/passwd" pass through
go run main.go serve --vulnerable=true

# fixed: the same baseline payload set is blocked with 403, including encoded variants
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
