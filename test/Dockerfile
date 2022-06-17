FROM golang:1.16

RUN apt update && apt install awscli -y
RUN go install github.com/golang/mock/mockgen@v1.6.0

COPY . /src/
WORKDIR /src
RUN go mod download

RUN make worker-pool
RUN make bin/e2e-test
ENTRYPOINT []
CMD []