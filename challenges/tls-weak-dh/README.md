# TLS Weak DH

This challenge demonstrates an HTTPS server that only negotiates finite-field Diffie-Hellman (DHE) cipher suites with a 1024-bit group. A 1024-bit DH group is within reach of well-funded attackers (Logjam, CVE-2015-4000): recorded traffic can be decrypted once the group is broken, which defeats the point of forward secrecy.

## How to run it

```bash
docker build -f Dockerfile -t tls-weak-dh ../..
docker run --rm -p 8080:8080 tls-weak-dh
```

## Modes

This challenge is a real nginx server built against OpenSSL, since the Go standard library has no DHE suites. The mode is toggled with the `VULNERABLE` environment variable at container startup (defaults to `true`):

```bash
# vulnerable: DHE suites only, with a 1024-bit group
docker run --rm -p 8080:8080 -e VULNERABLE=true tls-weak-dh

# fixed: ECDHE suites only
docker run --rm -p 8080:8080 -e VULNERABLE=false tls-weak-dh
```

## How to exploit it

```bash
openssl s_client -connect localhost:8080 -cipher 'DHE:@SECLEVEL=0' </dev/null 2>/dev/null | grep 'Server Temp Key'
# vulnerable -> Server Temp Key: DH, 1024 bits

proxyaudit scan https://localhost:8080
# vulnerable: tls.cipher_strength with code cipher.weak_dh
vulnapi scan curl https://localhost:8080 -k
# vulnerable: tls.forward_secrecy_missing
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
