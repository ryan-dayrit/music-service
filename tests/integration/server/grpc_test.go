//go:build integration

package server

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"music-service/gen/pb"
	handler "music-service/internal/handler/grpc"
	sqlxrepo "music-service/internal/repository/postgres/sqlx"
	"music-service/tests/integration/testutil"
)

func setupPostgresSqlx(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()
	ctx := context.Background()

	tmpFile, err := os.CreateTemp("", "integration_init_*.sql")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := tmpFile.WriteString(testutil.InitSQL); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to write init script: %v", err)
	}
	initScriptPath := tmpFile.Name()
	tmpFile.Close()
	t.Cleanup(func() { os.Remove(initScriptPath) })

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(initScriptPath),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping: %v", err)
	}

	cleanup := func() {
		db.Close()
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

func setupGrpcServer(t *testing.T, db *sqlx.DB) (string, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	repository := sqlxrepo.NewRepository(db)
	albumHandler := handler.NewAlbumHandler(repository)

	srv := grpc.NewServer()
	pb.RegisterMusicServiceServer(srv, albumHandler)

	go func() {
		if err := srv.Serve(listener); err != nil {
			t.Logf("grpc server error: %v", err)
		}
	}()

	// Brief delay for server to be ready
	time.Sleep(50 * time.Millisecond)

	addr := listener.Addr().String()
	cleanup := func() {
		srv.GracefulStop()
	}

	return addr, cleanup
}

func TestGrpcGetAlbumList_Empty(t *testing.T) {
	db, dbCleanup := setupPostgresSqlx(t)
	defer dbCleanup()

	addr, srvCleanup := setupGrpcServer(t, db)
	defer srvCleanup()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMusicServiceClient(conn)

	resp, err := client.GetAlbumList(context.Background(), &pb.GetAlbumsRequest{})
	if err != nil {
		t.Fatalf("GetAlbumList failed: %v", err)
	}

	if len(resp.Albums) != 0 {
		t.Errorf("expected 0 albums, got %d", len(resp.Albums))
	}
}

func TestGrpcGetAlbumList_WithData(t *testing.T) {
	db, dbCleanup := setupPostgresSqlx(t)
	defer dbCleanup()

	_, err := db.Exec("INSERT INTO music.albums (title, artist, price) VALUES ($1, $2, $3)",
		"Blue Train", "John Coltrane", 56.99)
	if err != nil {
		t.Fatalf("failed to insert: %v", err)
	}
	_, err = db.Exec("INSERT INTO music.albums (title, artist, price) VALUES ($1, $2, $3)",
		"Giant Steps", "John Coltrane", 63.99)
	if err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	addr, srvCleanup := setupGrpcServer(t, db)
	defer srvCleanup()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMusicServiceClient(conn)

	resp, err := client.GetAlbumList(context.Background(), &pb.GetAlbumsRequest{})
	if err != nil {
		t.Fatalf("GetAlbumList failed: %v", err)
	}

	if len(resp.Albums) != 2 {
		t.Fatalf("expected 2 albums, got %d", len(resp.Albums))
	}
	if resp.Albums[0].Title != "Blue Train" || resp.Albums[0].Artist != "John Coltrane" {
		t.Errorf("expected first album Blue Train by John Coltrane, got %q by %q",
			resp.Albums[0].Title, resp.Albums[0].Artist)
	}
	if resp.Albums[1].Title != "Giant Steps" || resp.Albums[1].Artist != "John Coltrane" {
		t.Errorf("expected second album Giant Steps by John Coltrane, got %q by %q",
			resp.Albums[1].Title, resp.Albums[1].Artist)
	}
}
