# TLS Known CVEs

This challenge demonstrates an HTTPS API whose TLS stack is affected by three well-known vulnerabilities:

- **Heartbleed (CVE-2014-0160)**: a heartbeat request declaring a longer payload than it carries is answered with server memory.
- **DROWN (CVE-2016-0800)**: SSLv2 is still accepted, so the RSA key can be attacked through SSLv2 to decrypt TLS sessions.
- **ROBOT**: the RSA key exchange answers well-formed and malformed PKCS#1 v1.5 padding with different alerts, a Bleichenbacher padding oracle.

Normal clients get a regular TLS 1.2+ connection served by Go's `crypto/tls`. In vulnerable mode, the server answers the SSLv2, heartbeat and RSA-only handshakes the way an affected server would, without the real cryptographic stack.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: Heartbleed, SSLv2 and RSA padding oracle
go run main.go serve --vulnerable=true

# fixed: every connection is served by crypto/tls (TLS 1.2+, no SSLv2, no heartbeat, no RSA key exchange)
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
proxyaudit scan https://localhost:8080
# vulnerable: tls.heartbleed (cve.heartbleed), tls.drown (cve.drown), tls.robot (cve.robot)
vulnapi scan curl https://localhost:8080 -k --scans "tls.*"
# vulnerable: tls.heartbleed, tls.drown, tls.robot
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
