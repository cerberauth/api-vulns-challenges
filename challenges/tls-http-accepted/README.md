# TLS HTTP Accepted

This challenge demonstrates an API that offers HTTPS but still answers plain HTTP on the same port with the same content instead of redirecting to HTTPS. A client that is given, or downgraded to, an `http://` URL sends its credentials and data unencrypted.

## How to run it

```bash
go run main.go serve
```

The server detects the TLS handshake on each connection and serves HTTPS and plain HTTP on the same port, with a self-signed certificate for `localhost`.

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: http://localhost:8080 serves the API
go run main.go serve --vulnerable=true

# fixed: http://localhost:8080 redirects (301) to https://localhost:8080
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
curl -i http://localhost:8080
# vulnerable -> 200 {"message": "ok"}, fixed -> 301 Location: https://localhost:8080/

vulnapi scan curl https://localhost:8080 -k
# vulnerable: tls.http_accepted
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
