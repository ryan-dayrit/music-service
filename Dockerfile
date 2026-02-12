# syntax=docker/dockerfile:1.6

FROM golang:1.25.5-bookworm AS build

WORKDIR /src

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        make \
        protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

ENV PATH="/root/go/bin:${PATH}"
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=build /src/bin/music-service /app/music-service
COPY --from=build /src/config.yaml /app/config.yaml

EXPOSE 50051

USER nonroot:nonroot

ENTRYPOINT ["/app/music-service", "server"]
