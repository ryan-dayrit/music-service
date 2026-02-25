//go:build integration

package kafka

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-pg/pg/v10"
	"github.com/shopspring/decimal"
	"github.com/testcontainers/testcontainers-go"
	tckafka "github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"google.golang.org/protobuf/proto"

	"music-service/gen/pb"
	confluentconsumer "music-service/internal/handler/kafka/confluent/consumer"
	saramaconsumer "music-service/internal/handler/kafka/sarama/consumer"
	"music-service/internal/models"
	"music-service/internal/repository/postgres/orm"
	appkafka "music-service/pkg/kafka"
	"music-service/tests/integration/testutil"
)

type consumerFactory func(cfg appkafka.Config, repository orm.Repository) (appkafka.ConsumerHandler, error)

func TestKafkaConsumers_ConsumeAlbumMessages_Integration(t *testing.T) {
	tests := []struct {
		name     string
		assignor string
		factory  consumerFactory
	}{
		{
			name:     "sarama",
			assignor: "sticky",
			factory:  saramaconsumer.NewConsumerHandler,
		},
		{
			name:     "confluent",
			assignor: "cooperative-sticky",
			factory:  confluentconsumer.NewConsumerHandler,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, broker, cleanup := setupPostgresAndKafka(t)
			defer cleanup()

			topic := fmt.Sprintf("albums-%s-%d", tt.name, time.Now().UnixNano())
			createTopic(t, broker, topic)

			repository := orm.NewRepository(db)
			handler, err := tt.factory(appkafka.Config{
				Brokers:       broker,
				Topics:        topic,
				ConsumerGroup: fmt.Sprintf("integration-%s-%d", tt.name, time.Now().UnixNano()),
				Assignor:      tt.assignor,
				Oldest:        true,
			}, repository)
			if err != nil {
				t.Fatalf("failed to create %s consumer: %v", tt.name, err)
			}

			consumeCtx, cancelConsume := context.WithCancel(context.Background())
			consumeErrCh := make(chan error, 1)
			go func() {
				consumeErrCh <- handler.Consume(consumeCtx)
			}()

			defer func() {
				cancelConsume()
				select {
				case consumeErr := <-consumeErrCh:
					if consumeErr != nil && !errors.Is(consumeErr, context.Canceled) {
						t.Errorf("%s consumer exited with error: %v", tt.name, consumeErr)
					}
				case <-time.After(15 * time.Second):
					t.Errorf("%s consumer did not stop before timeout", tt.name)
				}
			}()

			createAlbum := &pb.Album{
				Id:     0,
				Title:  fmt.Sprintf("integration-created-%s-%d", tt.name, time.Now().UnixNano()),
				Artist: "integration-artist",
				Price:  19.50,
			}
			produceAlbumMessage(t, broker, topic, createAlbum)

			createdAlbum := waitForAlbumByTitle(t, db, createAlbum.Title, 20*time.Second)
			if createdAlbum.Id <= 0 {
				t.Fatalf("expected created album id to be > 0, got %d", createdAlbum.Id)
			}
			if createdAlbum.Title != createAlbum.Title {
				t.Fatalf("expected created title %q, got %q", createAlbum.Title, createdAlbum.Title)
			}
			if createdAlbum.Artist != createAlbum.Artist {
				t.Fatalf("expected created artist %q, got %q", createAlbum.Artist, createdAlbum.Artist)
			}
			if !createdAlbum.Price.Equals(decimal.RequireFromString("19.50")) {
				t.Fatalf("expected created price 19.50, got %s", createdAlbum.Price.String())
			}

			const existingAlbumID = 4242
			_, err = db.Exec(
				"INSERT INTO music.albums (id, title, artist, price) OVERRIDING SYSTEM VALUE VALUES (?, ?, ?, ?)",
				existingAlbumID,
				"integration-before-update",
				"integration-before-artist",
				10.00,
			)
			if err != nil {
				t.Fatalf("failed to seed album for update path: %v", err)
			}

			updateAlbum := &pb.Album{
				Id:     existingAlbumID,
				Title:  "integration-after-update",
				Artist: "integration-after-artist",
				Price:  33.75,
			}
			produceAlbumMessage(t, broker, topic, updateAlbum)

			expectedUpdatedPrice := decimal.RequireFromString("33.75")
			updatedAlbum := waitForAlbumByID(
				t,
				db,
				existingAlbumID,
				updateAlbum.Title,
				updateAlbum.Artist,
				expectedUpdatedPrice,
				20*time.Second,
			)
			if updatedAlbum.Title != updateAlbum.Title {
				t.Fatalf("expected updated title %q, got %q", updateAlbum.Title, updatedAlbum.Title)
			}
			if updatedAlbum.Artist != updateAlbum.Artist {
				t.Fatalf("expected updated artist %q, got %q", updateAlbum.Artist, updatedAlbum.Artist)
			}
			if !updatedAlbum.Price.Equals(expectedUpdatedPrice) {
				t.Fatalf("expected updated price 33.75, got %s", updatedAlbum.Price.String())
			}
		})
	}
}

func setupPostgresAndKafka(t *testing.T) (*pg.DB, string, func()) {
	t.Helper()

	ctx := context.Background()

	tmpFile, err := os.CreateTemp("", "integration_init_*.sql")
	if err != nil {
		t.Fatalf("failed to create init temp file: %v", err)
	}
	if _, err := tmpFile.WriteString(testutil.InitSQL); err != nil {
		tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
		t.Fatalf("failed to write init script: %v", err)
	}
	initScriptPath := tmpFile.Name()
	_ = tmpFile.Close()
	t.Cleanup(func() { _ = os.Remove(initScriptPath) })

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

	kafkaContainer, err := tckafka.RunContainer(ctx)
	if err != nil {
		if terminateErr := testcontainers.TerminateContainer(postgresContainer); terminateErr != nil {
			t.Logf("failed to terminate postgres after kafka startup error: %v", terminateErr)
		}
		t.Fatalf("failed to start kafka container: %v", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		if terminateErr := testcontainers.TerminateContainer(kafkaContainer); terminateErr != nil {
			t.Logf("failed to terminate kafka after postgres connection string error: %v", terminateErr)
		}
		if terminateErr := testcontainers.TerminateContainer(postgresContainer); terminateErr != nil {
			t.Logf("failed to terminate postgres after connection string error: %v", terminateErr)
		}
		t.Fatalf("failed to get postgres connection string: %v", err)
	}

	opts, err := pg.ParseURL(connStr)
	if err != nil {
		if terminateErr := testcontainers.TerminateContainer(kafkaContainer); terminateErr != nil {
			t.Logf("failed to terminate kafka after parse url error: %v", terminateErr)
		}
		if terminateErr := testcontainers.TerminateContainer(postgresContainer); terminateErr != nil {
			t.Logf("failed to terminate postgres after parse url error: %v", terminateErr)
		}
		t.Fatalf("failed to parse postgres connection string: %v", err)
	}

	db := pg.Connect(opts)
	if err := db.Ping(ctx); err != nil {
		db.Close()
		if terminateErr := testcontainers.TerminateContainer(kafkaContainer); terminateErr != nil {
			t.Logf("failed to terminate kafka after ping error: %v", terminateErr)
		}
		if terminateErr := testcontainers.TerminateContainer(postgresContainer); terminateErr != nil {
			t.Logf("failed to terminate postgres after ping error: %v", terminateErr)
		}
		t.Fatalf("failed to ping postgres: %v", err)
	}

	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		db.Close()
		if terminateErr := testcontainers.TerminateContainer(kafkaContainer); terminateErr != nil {
			t.Logf("failed to terminate kafka after broker lookup error: %v", terminateErr)
		}
		if terminateErr := testcontainers.TerminateContainer(postgresContainer); terminateErr != nil {
			t.Logf("failed to terminate postgres after broker lookup error: %v", terminateErr)
		}
		t.Fatalf("failed to get kafka brokers: %v", err)
	}
	if len(brokers) == 0 {
		db.Close()
		if terminateErr := testcontainers.TerminateContainer(kafkaContainer); terminateErr != nil {
			t.Logf("failed to terminate kafka after empty broker list: %v", terminateErr)
		}
		if terminateErr := testcontainers.TerminateContainer(postgresContainer); terminateErr != nil {
			t.Logf("failed to terminate postgres after empty broker list: %v", terminateErr)
		}
		t.Fatalf("kafka broker list is empty")
	}

	cleanup := func() {
		db.Close()
		if err := testcontainers.TerminateContainer(kafkaContainer); err != nil {
			t.Logf("failed to terminate kafka container: %v", err)
		}
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	}

	return db, brokers[0], cleanup
}

func createTopic(t *testing.T, broker, topic string) {
	t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	var lastErr error

	for time.Now().Before(deadline) {
		cfg, err := newSaramaConfig()
		if err != nil {
			lastErr = err
			time.Sleep(500 * time.Millisecond)
			continue
		}

		admin, err := sarama.NewClusterAdmin([]string{broker}, cfg)
		if err != nil {
			lastErr = err
			time.Sleep(500 * time.Millisecond)
			continue
		}

		err = admin.CreateTopic(topic, &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		closeErr := admin.Close()
		if err == nil && closeErr == nil {
			return
		}

		if err == nil {
			lastErr = closeErr
		} else {
			lastErr = err
		}
		time.Sleep(500 * time.Millisecond)
	}

	t.Fatalf("failed to create kafka topic %q: %v", topic, lastErr)
}

func produceAlbumMessage(t *testing.T, broker, topic string, album *pb.Album) {
	t.Helper()

	cfg, err := newSaramaConfig()
	if err != nil {
		t.Fatalf("failed to build sarama config: %v", err)
	}

	producer, err := sarama.NewSyncProducer([]string{broker}, cfg)
	if err != nil {
		t.Fatalf("failed to create sarama producer: %v", err)
	}
	defer func() {
		if closeErr := producer.Close(); closeErr != nil {
			t.Logf("failed to close sarama producer: %v", closeErr)
		}
	}()

	messageValue, err := proto.Marshal(album)
	if err != nil {
		t.Fatalf("failed to marshal album message: %v", err)
	}

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageValue),
	})
	if err != nil {
		t.Fatalf("failed to send kafka message: %v", err)
	}
}

func waitForAlbumByTitle(t *testing.T, db *pg.DB, title string, timeout time.Duration) *models.Album {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		album := &models.Album{}
		err := db.Model(album).
			Where("title = ?", title).
			Order("id ASC").
			Limit(1).
			Select()
		if err == nil {
			return album
		}
		if !errors.Is(err, pg.ErrNoRows) {
			t.Fatalf("failed while checking created album by title: %v", err)
		}

		time.Sleep(200 * time.Millisecond)
	}

	t.Fatalf("timed out waiting for album with title %q", title)
	return nil
}

func waitForAlbumByID(
	t *testing.T,
	db *pg.DB,
	id int,
	expectedTitle string,
	expectedArtist string,
	expectedPrice decimal.Decimal,
	timeout time.Duration,
) *models.Album {
	t.Helper()

	deadline := time.Now().Add(timeout)
	var lastAlbum *models.Album

	for time.Now().Before(deadline) {
		album := &models.Album{Id: id}
		err := db.Model(album).WherePK().Select()
		if err == nil && album.Title == expectedTitle && album.Artist == expectedArtist && album.Price.Equals(expectedPrice) {
			return album
		}
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			t.Fatalf("failed while checking album by id: %v", err)
		}

		lastAlbum = album
		time.Sleep(200 * time.Millisecond)
	}

	if lastAlbum != nil {
		t.Fatalf(
			"timed out waiting for album id %d to match title=%q artist=%q price=%s; last seen title=%q artist=%q price=%s",
			id,
			expectedTitle,
			expectedArtist,
			expectedPrice.String(),
			lastAlbum.Title,
			lastAlbum.Artist,
			lastAlbum.Price.String(),
		)
	}

	t.Fatalf("timed out waiting for album with id %d", id)
	return nil
}

func newSaramaConfig() (*sarama.Config, error) {
	cfg := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion(sarama.DefaultVersion.String())
	if err != nil {
		return nil, err
	}
	cfg.Version = version

	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5

	return cfg, nil
}
