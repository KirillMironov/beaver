FROM golang:1.20.3-alpine3.17 AS builder
WORKDIR /usr/src/beaver
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go test -v ./...
RUN go build -o /usr/local/bin/beaver ./cmd/beaver

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/bin/beaver .
ENTRYPOINT ["./beaver"]
