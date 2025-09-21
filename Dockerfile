# build stage
ARG GO_VERSION=1.24.1
FROM golang:${GO_VERSION} AS builder
ENV CGO_ENABLED=0 GOOS=linux GOTOOLCHAIN=auto

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o reputationhistory ./cmd/main.go

# final stage
FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/reputationhistory /app/reputationhistory
ENTRYPOINT ["/app/reputationhistory"]
