BINDIR := bin

.PHONY: ${BINDIR}/worker-pool worker-pool

worker-pool: $(BINDIR)/worker-pool

$(BINDIR)/worker-pool:
	go build -o $@ ./cmd/worker-pool

.PHONY: test
test:
	go test ./...