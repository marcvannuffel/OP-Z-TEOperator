# ── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.21-bookworm AS builder

WORKDIR /app
COPY . .
RUN go build -o teoperator .

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM debian:bookworm-slim

# Install runtime dependencies (ffmpeg + audiowaveform for waveform images)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ffmpeg \
    audiowaveform \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/teoperator .
COPY --from=builder /app/src/server/templates ./src/server/templates
COPY --from=builder /app/src/server/static ./src/server/static

# Default port – override with -e PORT=... or docker run --port
EXPOSE 8080

CMD ["./teoperator", "server", "--port", "8080"]
