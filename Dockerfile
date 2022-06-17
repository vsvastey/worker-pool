FROM golang:1.18.3

RUN go install github.com/golang/mock/mockgen@v1.6.0
RUN apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

COPY . /src/
WORKDIR /src
RUN go mod download
RUN make worker-pool

ENTRYPOINT ["/src/bin/worker-pool"]
CMD ["--config",  "/src/config/sleep_only.yaml"]