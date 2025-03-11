# syntax=docker.io/docker/dockerfile-upstream:1.14.1

# Build the manager binary
FROM golang:1.23.2 AS builder
WORKDIR /vcluster

# Copy the Go Modules manifests
COPY go.mod go.sum /vcluster/

# Install dependencies
RUN go mod download

# Copy the sources
COPY --parents pkg main.go /vcluster/

# Build plugin
RUN CGO_ENABLED=0 go build -o /plugin main.go

# we use alpine for easier debugging
FROM alpine
WORKDIR /plugin
COPY --from=builder /plugin .
