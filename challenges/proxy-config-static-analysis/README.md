# Config-Level Static Analysis

This challenge is not a running server: it is a set of sample reverse-proxy configuration files meant to exercise a config-level static analysis / lint subcommand (e.g. `proxyaudit lint <config>`) against known misconfiguration patterns for Nginx, Traefik, Envoy and Caddy, each with a vulnerable and a fixed variant.

## Layout

```
configs/
  nginx/vulnerable.conf     nginx/fixed.conf
  traefik/vulnerable.yml    traefik/fixed.yml
  envoy/vulnerable.yaml     envoy/fixed.yaml
  caddy/vulnerable.Caddyfile  caddy/fixed.Caddyfile
```

## Findings covered

- **Nginx**: missing root-location fallback exposing the filesystem, off-by-slash alias traversal, DNS resolver pointed at an untrusted/public server, blind `Host` forwarding upstream.
- **Traefik**: dashboard/API exposed without authentication (`api.insecure`), forwarded headers trusted from any client, wildcard CORS origin combined with credentials.
- **Envoy**: `xff_num_trusted_hops` misconfigured with `use_remote_address: false` (client-controlled trust boundary), admin interface bound to all interfaces instead of loopback.
- **Caddy**: `templates` directive enabling server-side template injection via reflected headers, blind `X-Forwarded-Host` forwarding, directory listing enabled.

## How to use it

Point your static analyzer at either variant, e.g.:

```bash
proxyaudit lint configs/nginx/vulnerable.conf   # expected: findings reported
proxyaudit lint configs/nginx/fixed.conf        # expected: no findings
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
