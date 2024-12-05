FROM golang:1.19.3 AS builder

WORKDIR /code
COPY . /code
RUN go build -v -o rmq-tools main.go

FROM centos:8
COPY --from=builder /code/rmq-tools /usr/local/bin/rmq-tools
RUN chmod +x /usr/local/bin/rmq-tools
ENV ROCKETMQ_GO_LOG_LEVEL=error
CMD rmq-tools --role=receiver