FROM golang:1.26.5-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

ARG SERVICE
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/provider ./cmd/${SERVICE}

FROM alpine:3.22

RUN addgroup -S provider && adduser -S provider -G provider
COPY --from=build /out/provider /usr/local/bin/provider

USER provider
EXPOSE 8080
ENTRYPOINT ["provider"]
