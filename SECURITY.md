# Security Policy

## Reporting

Report suspected vulnerabilities through a
[private GitHub security advisory](https://github.com/faustbrian/go-wsdl/security/advisories/new).
Do not open a public issue containing exploit details, credentials, private
fixtures, or affected deployment information.

Include the affected module and version, impact, reproduction, preconditions,
and any suggested mitigation. Reports are acknowledged as soon as practical;
timelines depend on severity and verification.

## Supported Versions

The latest stable release is
[`github.com/faustbrian/go-wsdl/v2` v2.0.0](https://github.com/faustbrian/go-wsdl/releases/tag/v2.0.0).
Support windows are documented per module and in
[`COMPATIBILITY.md`](COMPATIBILITY.md).

The v2 release includes categorical default errors and the owned cancellation
checks described in the [versioned threat model](docs/security-threat-model.md).
Its publication does not establish that historical v1 versions contain those
fixes. Reports concerning v1 remain subject to assessment through the private
reporting route above. Adoption of the separate XSD major remains pending.

## Security Gates

Releases require isolated tests, race and hostile-input checks, exact coverage
and mutation results, `govulncheck`, secret scanning, license verification,
SBOM generation, provenance validation, and clean-consumer resolution. A
missing scanner or unavailable service is a failed gate, not a warning.

Security fixes MUST include a regression test that does not publish weaponized
details or real secrets. Credentials MUST be redacted from logs and evidence.

## Repository Assurance

The repository [safety and concurrency policy](AGENTS.md#safety-and-concurrency)
and [supply-chain policy](AGENTS.md#dependencies-and-supply-chain) define shared
trust boundaries and release requirements. Package-specific security guidance
refines those rules for its owned boundary.
