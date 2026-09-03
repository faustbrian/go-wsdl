# Module identity reviews

## 2026-09-03 public proxy reconciliation

The WSDL module retains its existing `go-wire` and `go-xsd` v1.0.0
requirements. This review changes only their recorded module-content checksums
from the historical bootstrap archive identities to the immutable identities
served by the public Go proxy and authenticated by the Go checksum database.
The corresponding `go.mod` checksums were already public and remain unchanged.

| Module | Previous module checksum | Public module checksum | Tag object | Peeled commit |
| --- | --- | --- | --- | --- |
| `github.com/faustbrian/go-wire@v1.0.0` | `h1:m0Rcz+vdhobAKzaySFmMFgwq4qTCzeKhJG1wQ0n6dKM=` | `h1:vVniAQ0+dtWPlYvLzJc4xM74oyZ4ZGo5IvQKBTO6UyM=` | `3f6e9e330bd76ff0667722c52e540e5ff895ef51` | `6fe443d8da77446fe6eb0f1ccbc7fb869ba97c98` |
| `github.com/faustbrian/go-xsd@v1.0.0` | `h1:/CCCiNo7M317iF0TWk4lkpJhT9vg+vhyCcWktPFcq50=` | `h1:BT9PtbTpTGt8/s0ACq4aFLsEn2ETAjPgQtSu3nSSrkI=` | `2156538376045e1e3357e41406375bfb3e457911` | `8f93e6ebd80508deeb73de7fc91adeb24d2c4b8c` |

The reviewed authorities are the public proxy archives and checksum-database
records for
[`go-wire@v1.0.0`](https://proxy.golang.org/github.com/faustbrian/go-wire/@v/v1.0.0.zip)
and
[`go-xsd@v1.0.0`](https://proxy.golang.org/github.com/faustbrian/go-xsd/@v/v1.0.0.zip).
The checksum records are available from the checksum database for
[`go-wire@v1.0.0`](https://sum.golang.org/lookup/github.com/faustbrian/go-wire@v1.0.0)
and
[`go-xsd@v1.0.0`](https://sum.golang.org/lookup/github.com/faustbrian/go-xsd@v1.0.0).

The v1.4.0 shared workflow configures the verified bootstrap archive after the
public proxy in the comma-separated proxy chain. Publicly available versions
therefore cannot be shadowed by historical bootstrap bytes. This identity
reconciliation changes no dependency version, public WSDL API, specification
decision, conformance binding, or runtime policy.
