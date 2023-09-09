# syntax=docker/dockerfile:1.2

FROM golang:1.21 AS builder
ARG APP_NAME
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
# Build
RUN go build -o /app/out/app /app/cmd/$APP_NAME/main.go


FROM debian:bookworm-slim as runner
COPY --from=builder /app/out/app /app
COPY ./db/migrations /migrations 
# Run
CMD ["/app"]