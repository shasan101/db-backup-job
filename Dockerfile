# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux make build



# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/bin/db-backup /db-backup

USER nonroot:nonroot

ENTRYPOINT ["/db-backup"]
