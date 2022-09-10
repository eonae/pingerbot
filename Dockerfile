FROM golang:1.18-alpine as builder

WORKDIR /usr/src

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

RUN go mod download && go mod verify

RUN go build -v -o ./build/pingerbot ./cmd/main.go

FROM alpine:latest as exec

COPY --from=builder /usr/src/build/pingerbot /opt

CMD ["/opt/pingerbot"]
