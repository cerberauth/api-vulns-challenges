# Nginx-Specific Misconfigurations

This challenge demonstrates the classic Nginx off-by-slash alias traversal: a `location` directive without a trailing slash combined with an `alias` directive lets a path like `/files../secret/flag.txt` escape the intended alias root, because Nginx strips only the literal `/files` prefix and appends the remainder verbatim to the alias path.

## How to run it

```bash
docker build -f Dockerfile -t nginx-alias-traversal ../..
docker run --rm -p 8080:8080 nginx-alias-traversal
```

## Modes

Unlike the Go-based challenges in this repository, this challenge is a real Nginx server, so the mode is toggled with the `VULNERABLE` environment variable at container startup (defaults to `true`):

```bash
# vulnerable: location /files (no trailing slash) + alias escapes the alias root
docker run --rm -p 8080:8080 -e VULNERABLE=true nginx-alias-traversal

# fixed: location /files/ (trailing slash) requires the path to stay under /files/
docker run --rm -p 8080:8080 -e VULNERABLE=false nginx-alias-traversal
```

```bash
curl http://localhost:8080/files/index.html        # public file, always reachable
curl http://localhost:8080/files../secret/flag.txt  # only reachable in vulnerable mode
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
