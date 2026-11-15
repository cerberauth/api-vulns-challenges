# TLS Certificate Misconfiguration

This challenge demonstrates an HTTPS API served with a certificate no client should trust: it is self-signed, it expired a year ago, and it was issued for `internal.example.com` rather than the host being called. Clients that skip certificate verification to "make it work" are open to interception.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: self-signed, expired certificate for internal.example.com
go run main.go serve --vulnerable=true

# fixed: valid certificate for localhost / 127.0.0.1, signed by the challenge CA
go run main.go serve --vulnerable=false
```

The challenge CA is derived from a fixed seed, so `go run main.go ca` always prints the same certificate. Trust it to scan the fixed server without any certificate finding:

```bash
go run main.go ca > ca.pem
SSL_CERT_FILE=ca.pem curl https://localhost:8080
```

## How to exploit it

```bash
openssl s_client -connect localhost:8080 -servername localhost </dev/null 2>/dev/null | openssl x509 -noout -subject -issuer -dates

vulnapi scan curl https://localhost:8080 -k
# vulnerable: tls.certificate_untrusted, tls.certificate_expired, tls.certificate_hostname_mismatch

SSL_CERT_FILE=ca.pem vulnapi scan curl https://localhost:8080
# fixed: no certificate finding
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
