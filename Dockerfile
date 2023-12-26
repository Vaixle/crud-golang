ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR ${GOPATH}/src/crud-golang

COPY go.mod ${GOPATH}/src/crud-golang
COPY go.sum ${GOPATH}/src/crud-golang
RUN go mod download

COPY ./cmd/app/main.go ${GOPATH}/src/crud-golang
COPY config ${GOPATH}/src/crud-golang/config
COPY internal ${GOPATH}/src/crud-golang/internal
COPY migration ${GOPATH}/src/crud-golang/migration
COPY docs ${GOPATH}/src/crud-golang/docs
COPY pkg ${GOPATH}/src/crud-golang/pkg

RUN go build -o /app main.go

FROM alpine:latest

WORKDIR /api
COPY --from=builder /app .
COPY --from=builder /go/src/crud-golang/config ./config

EXPOSE 8088

CMD ["/api/app"]
