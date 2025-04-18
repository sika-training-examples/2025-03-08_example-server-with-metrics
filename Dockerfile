FROM golang:1.24.2 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN \
  --mount=type=cache,target=/go \
  go mod download
COPY . .
ENV CGO_ENABLED=0
ARG VERSION=master
RUN \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go \
  go build -ldflags="-s -w -X example/version.Version=${VERSION}"

FROM debian:12-slim
COPY --from=build /build/example /usr/local/bin/
CMD ["/usr/local/bin/example"]
EXPOSE 8000
