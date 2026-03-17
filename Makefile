BINARY      = devtool
INSTALL_DIR = /usr/local/bin
GO_IMAGE    = golang:1.26.1

# ── Inside devcontainer ──────────────────────────────────────────────────────

.PHONY: build run lint clean

build:
	go build -o $(BINARY) .

run:
	go run . $(ARGS)

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

# ── From host (no devcontainer needed) ───────────────────────────────────────
# Uses a named volume to cache Go modules between builds.

.PHONY: install uninstall

install:
	docker run --rm \
		-v $(shell pwd):/workspace \
		-v devtool-gomodcache:/go/pkg/mod \
		-w /workspace \
		$(GO_IMAGE) \
		go build -buildvcs=false -o $(BINARY) .
	sudo install -m 755 $(BINARY) $(INSTALL_DIR)/$(BINARY)
	rm -f $(BINARY)

uninstall:
	sudo rm -f $(INSTALL_DIR)/$(BINARY)
