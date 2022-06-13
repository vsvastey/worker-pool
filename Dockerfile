FROM golang:1.18.3

RUN go install github.com/golang/mock/mockgen@v1.6.0

COPY . /src/
WORKDIR /src
RUN go mod download
RUN make worker-pool

ENTRYPOINT ["/src/bin/worker-pool"]
CMD ["--config",  "/src/config/sleep_only.yaml"]