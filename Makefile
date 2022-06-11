BINDIR := bin

.PHONY: ${BINDIR}/worker-pool worker-pool

worker-pool: $(BINDIR)/worker-pool

$(BINDIR)/worker-pool:
	go build -o $@ ./cmd/worker-pool

.PHONY: test
test:
	go test ./...

.PHONY: mocks
mocks: task-mock

.PHONY: task-mock
task-mock:
	mockgen -package=task -source=internal/task/task.go -destination=internal/task/task_mock.go
