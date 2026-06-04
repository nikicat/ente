FROM --platform=$BUILDPLATFORM golang:1.23.12-alpine3.22 AS builder
ENV CGO_ENABLED=0

# Populated automatically by buildx; used to cross-compile natively on the
# build host (no QEMU emulation) since this is a pure-Go, CGO-free build.
ARG TARGETOS
ARG TARGETARCH

WORKDIR /etc/ente/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
  GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -o museum cmd/museum/main.go

FROM alpine:3.22
COPY --from=builder /etc/ente/museum .
COPY configurations configurations
COPY migrations migrations
COPY mail-templates mail-templates
COPY web-templates web-templates

ARG GIT_COMMIT
ENV GIT_COMMIT=$GIT_COMMIT

CMD ["./museum"]
