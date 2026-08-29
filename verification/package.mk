GO ?= go
FUZZ_TIME ?= 5s
BENCH_TIME ?= 100ms

.PHONY: benchmark conformance fuzz interoperability

benchmark:
	./scripts/check-benchmarks.sh "$(BENCH_TIME)"

conformance:
	./scripts/check-provenance.sh
	$(GO) test ./... -count=1

fuzz:
	./scripts/check-fuzz.sh "$(FUZZ_TIME)"

interoperability:
	./scripts/check-woden.sh
