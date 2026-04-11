FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY *.go ./
RUN go build -o genChallResult .

FROM alpine:3.20

WORKDIR /work
COPY --from=builder /app/genChallResult /usr/local/bin/genChallResult

ENTRYPOINT ["genChallResult"]
