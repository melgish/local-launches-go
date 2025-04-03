FROM mcr.microsoft.com/devcontainers/go:1-1.24-bookworm AS base
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

FROM base AS build
COPY *.go ./
COPY www/* ./www/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./local-launches

FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/local-launches .
ENTRYPOINT ["/app/local-launches"]