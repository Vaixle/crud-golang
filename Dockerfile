ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR ${GOPATH}/src/empha-soft

COPY go.mod ${GOPATH}/src/empha-soft
COPY go.sum ${GOPATH}/src/empha-soft
RUN go mod download

COPY ./cmd/app/main.go ${GOPATH}/src/empha-soft
COPY config ${GOPATH}/src/empha-soft/config
COPY internal ${GOPATH}/src/empha-soft/internal
COPY migration ${GOPATH}/src/empha-soft/migration
COPY docs ${GOPATH}/src/empha-soft/docs
COPY pkg ${GOPATH}/src/empha-soft/pkg

RUN go build -o /app main.go

FROM alpine:latest

WORKDIR /api
COPY --from=builder /app .
COPY --from=builder /go/src/empha-soft/config ./config

EXPOSE 8088

CMD ["/api/app"]