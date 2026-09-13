# Request Smuggling / Desync

This challenge implements a minimal raw HTTP/1.1 server (bypassing Go's `net/http`, which already rejects ambiguous requests) to demonstrate a TE.CL-style desync: when both `Content-Length` and `Transfer-Encoding: chunked` are present, the request boundary is determined by the chunked body, and any bytes appended after the terminating chunk are buffered on the connection and answered as if they were a second, smuggled request.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: a request with both Content-Length and Transfer-Encoding: chunked is accepted,
# and content appended after the chunked body is processed as a smuggled second request
go run main.go serve --vulnerable=true

# fixed: any request with both headers present is rejected with 400, per RFC 7230 §3.3.3
go run main.go serve --vulnerable=false
```

Example probe with `netcat`, sending a chunked body that terminates early and smuggles a second request:

```bash
printf 'POST / HTTP/1.1\r\nHost: localhost\r\nContent-Length: 4\r\nTransfer-Encoding: chunked\r\n\r\n0\r\n\r\nG' | nc localhost 8080
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
