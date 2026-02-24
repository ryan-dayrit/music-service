# AGENTS.md

## Cursor Cloud specific instructions

### Overview

This is a Go + Svelte music-service application. It uses REST (Fiber), gRPC, Kafka, and PostgreSQL. See `README.md` for architecture details.

### Infrastructure

Local dev infrastructure is defined in `docker-compose.dev.yaml`: Postgres 16, Kafka (KRaft/native image), and Kafka UI.

```bash
docker compose -f docker-compose.dev.yaml up -d
```

**Gotcha:** The Kafka native image (`apache/kafka-native`) does not include shell scripts like `kafka-broker-api-versions.sh`, so the healthcheck in `docker-compose.dev.yaml` will always report `unhealthy`. Kafka is actually running fine — start `kafka-ui` manually with `docker start workspace-kafka-ui-1` if it failed to start due to the healthcheck dependency.

**Gotcha:** `config.yaml` postgres password must match the docker-compose.dev.yaml password (`postgres`). If they diverge, the REST server will fail to connect to Postgres.

### Running services

- **REST server:** `go run . rest-server` (port 3000). Must be run from the repo root (reads `config.yaml` from cwd).
- **Kafka consumer:** `go run . kafka-consumer-sarama` — processes messages from Kafka into Postgres. Required for write-path to complete.
- **Svelte UI:** `cd ui && npm run dev` (port 5173). Hardcodes REST API at `http://localhost:3000`.

### Commands

| Task | Command |
|------|---------|
| Unit tests | `go test ./...` |
| Integration tests | `go test -tags=integration -v -count=1 ./tests/integration/...` (requires Docker) |
| Lint | `go vet ./...` |
| Build | `go build -o bin/music-service .` |
| Protobuf codegen | `make gen` (requires `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`) |
| UI dev server | `cd ui && npm run dev` |
| UI e2e tests | `cd ui && npm run test:e2e` (requires REST server + Kafka consumer running) |

### System dependencies

These must be installed on the VM for full development support:
- Docker (for docker-compose infrastructure and integration tests via testcontainers)
- `protobuf-compiler` (for `make gen`)
- `librdkafka-dev` (for `confluent-kafka-go` CGO dependency)
- `protoc-gen-go` and `protoc-gen-go-grpc` (Go protobuf plugins)

### Write path flow

POST/PUT to REST API → Kafka produce → Kafka consumer → Postgres insert/update. The GET endpoint reads directly from Postgres. Without a running Kafka consumer, POSTed albums won't appear in GET responses.
