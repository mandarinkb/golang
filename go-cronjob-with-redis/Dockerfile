# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.23-bookworm AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/api .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/worker ./worker

# ---- Runtime stage ----
FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends bash ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /bin/api    /bin/api
COPY --from=build /bin/worker /bin/worker
COPY script/                  /app/script/

RUN chmod +x /app/script/*

EXPOSE 3000

# Default: run API server
# สามารถ override ได้ใน docker-compose ด้วย command: ["/bin/worker"]
CMD ["/bin/api"]
