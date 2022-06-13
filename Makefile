BINDIR := bin

.PHONY: ${BINDIR}/worker-pool worker-pool

worker-pool: $(BINDIR)/worker-pool

$(BINDIR)/worker-pool:
	go build -o $@ ./cmd/worker-pool

.PHONY: test
test: mocks
	go test ./...

.PHONY: mocks
mocks: task-mock

.PHONY: task-mock
task-mock:
	mockgen -package=task -source=internal/task/task.go -destination=internal/task/task_mock.go

.PHONE: docker-worker-pool
docker-worker-pool:
	docker build -f ./Dockerfile . -t worker-pool:latest

.PHONE: demo
demo: docker-worker-pool
	docker run -it worker-pool:latest