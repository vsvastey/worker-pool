BINDIR := bin

.PHONY: worker-pool
worker-pool: $(BINDIR)/worker-pool

.PHONY: ${BINDIR}/worker-pool
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

.PHONY: $(BINDIR)/e2e-test
$(BINDIR)/e2e-test:
	@go build -race -o $@ ./test/e2e/

.PHONY: e2e-test
e2e-test: docker-worker-pool
	docker-compose -f deployment/e2e-test/docker-compose.yaml up --exit-code-from test --abort-on-container-exit --build
	-docker-compose -f deployment/e2e-test/docker-compose.yaml down
