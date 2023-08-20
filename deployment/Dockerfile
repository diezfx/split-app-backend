# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder
ARG APP_NAME
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
# Build
RUN go build -o /app/out/$APP_NAME /app/cmd/$APP_NAME/main.go


FROM debian:bullseye-slim as runner
COPY --from=builder /app/out/${APP_NAME} /app/app
# Run
CMD ["/app/app"]