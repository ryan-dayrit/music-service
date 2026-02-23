# music-service
practice Golang backend application using
  * gRPC
  * REST
  * Kafka
  * PostgreSQL
  * Protobuf

UI is written in Svelte
  * It shows abums in Postgres by calling the REST API
  * It adds an album enterred by sending it to REST API

UI e2e tests are written in Cypress

Backend integration tests use testcontainers-go to spin up a PostgreSQL container
  * They are excluded from normal test runs via the integration build tag

# Components 
1. PostgreSQL database in minikube with DDL/DML initialization scripts
2. Kafka in docker with a topic and a consumer group
3. gRPC API (internal) to service the data in the Postgres database 
4. REST API (external) using fiber library to receive JSON payloads and produce a Protobuf message to a kakfa topic
5. Kafka consumer using sarama library to process Protobuf messages from a Kafka topic
6. Kafka consumer using confluent library to process Protobuf messages from a Kafka topic
7. UI using svelte which calls the REST API to show the results and add an album
8. Cypress acceptance tests for the UI

# Processing Flow
Writes 
1. REST API POST/PUT receiver for json payloads
2. REST API publishes proto to Kafka
3. Kafka consumers consume the proto from the kafka topic and creates/updates PostgreSQL

Reads 
1. gRPC API (internal) which reads from the PostgreSQL database using Sqlx library and returns protos in json format
2. REST API (exteranl) which reads from the PostgreSQL database using ORM library and returns protos in json format
3. Svelte UI (external) which calls the REST API (external)

# CLI Testers
1. REST API client which sends POST/PUT requests
2. gRPC API client which calls gRPC API and prints response
3. Kakfa producer using sarama library which sends marshalled protos to the Kafka topic 
4. Kakfa producer using confluent library which sends marshalled protos to the Kafka topic
5. PostgreSQL client which reads gets all albums from the PostgreSQL database direcly using Sqlx library 
6. PostgreSQL client which reads gets album by id from the PostgreSQL database direcly using ORM library 
7. PostgreSQL client which inserts PostgreSQL database direcly using ORM library 
8. REST API client which sends put/post requests to REST API album endpoint
9. REST API client which sends put/post requests to REST API albums endpoint

# Unit Tests
Run unit tests:
```bash
go test ./...
```

# Integration Tests
Run integration tests (requires Docker):
```bash
go test -tags=integration ./tests/integration/...
```

Tests cover REST API, gRPC server, Sqlx repository, and ORM repository

# E2E Tests
```bash
cd ui
npm run dev           # start dev server in one terminal
npm run test:e2e      # run tests headlessly
npm run test:e2e:open # open Cypress interactive UI
```

