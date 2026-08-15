# syntax=docker/dockerfile:1
# Batam MedHub core backend image.
#
# Builds BOTH the API server and the migration CLI so a single image can:
#   - apply core_db migrations  (command: ["migrate"])
#   - serve the HTTP API        (command: ["api"], the default)
#
# Migrations are copied into /app/migrations because cmd/migrate resolves them
# relative to the working directory via `file://migrations`.

FROM golang:1.26.5-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

# CGO disabled => static binaries that run on scratch/alpine without libc.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.22

RUN addgroup -S medhub && adduser -S medhub -G medhub

WORKDIR /app
COPY --from=build /out/api /usr/local/bin/api
COPY --from=build /out/migrate /usr/local/bin/migrate
COPY migrations ./migrations

USER medhub
EXPOSE 8080

CMD ["api"]
