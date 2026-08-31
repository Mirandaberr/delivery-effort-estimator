.PHONY: fmt vet test build sync-plugin-engine check-plugin-engine-sync test-plugin

fmt:
	gofmt -l .

vet:
	go vet ./...

test: check-plugin-engine-sync
	go test ./...

build:
	go build ./...

sync-plugin-engine:
	./scripts/sync-plugin-engine.sh

check-plugin-engine-sync:
	./scripts/check-plugin-engine-sync.sh

test-plugin:
	for t in plugin/tests/*.sh; do sh "$$t" || exit 1; done
