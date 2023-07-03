# syntax=docker/dockerfile:1.2

# STAGE 1: build the executable
FROM golang:1.20-alpine AS builder

RUN apk add --no-cache git openssh ca-certificates
WORKDIR /src

ARG TARGETARCH

# IF statement executed due to incosistent package names @ https://github.com/moparisthebest/static-curl/issues/8
RUN if [ "$TARGETARCH" = "arm64" ]; then \
        TARGETARCH=aarch64 ; \
        fi; \
    wget -O /usr/bin/curl https://github.com/moparisthebest/static-curl/releases/download/v8.0.1/curl-$TARGETARCH \
        && chmod +x /usr/bin/curl

RUN git config --global --add url."ssh://git@github.com/twlabs/".insteadOf "https://github.com/twlabs/"

RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

COPY ./go.mod ./go.sum ./

ENV GOPRIVATE=github.com/twlabs
RUN --mount=type=ssh go mod download

COPY ./ ./

RUN --mount=type=cache,target=/root/.cache/go-build/ CGO_ENABLED=0 go build -o /app main.go

# STAGE 2: build the container
FROM gcr.io/distroless/static:nonroot AS final

COPY --from=builder /usr/bin/curl /usr/bin/curl

USER 10000:10000

COPY --from=builder --chown=nonroot:nonroot /app /app

ENTRYPOINT ["/app"]
