FROM golang:1.24-bullseye AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -a -o operator .

FROM debian:bullseye-slim

RUN set -eux; \
    if ! getent group operator >/dev/null; then \
      addgroup --system operator; \
    fi; \
    if ! id operator >/dev/null 2>&1; then \
      adduser --system --ingroup operator --no-create-home operator; \
    fi

COPY --from=builder /src/operator /usr/local/bin/operator

USER operator

ENTRYPOINT ["/usr/local/bin/operator"]
