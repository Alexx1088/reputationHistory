package consumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"

	kafkamodel "github.com/Alexx1088/reputationhistory/internal/kafka"
	"github.com/Alexx1088/reputationhistory/internal/repository"
)

type ReputationConsumer struct {
	reader *kafka.Reader
	repo   *repository.Repo
}

func NewReputationConsumer(db *sql.DB) *ReputationConsumer {
	brokersEnv := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	if brokersEnv == "" {
		// локально НЕ надо, но оставим как защиту
		brokersEnv = "kafka:9092"
	}
	brokers := strings.Split(brokersEnv, ",")
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	log.Printf("Kafka config: brokers=%v topic=%s group=%s", brokers, topic, groupID)

	return &ReputationConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 1,
			MaxBytes: 10e6,
		}),
		repo: &repository.Repo{DB: db},
	}
}

func (c *ReputationConsumer) Run(ctx context.Context) error {
	log.Println("Kafka consumer started...")
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var ev kafkamodel.ReputationEntryEvent
		if err := json.Unmarshal(m.Value, &ev); err != nil {
			log.Printf("json decode error: %v", err)
			continue
		}

		// применяем к БД
		ctxDB, cancel := context.WithTimeout(ctx, 5*time.Second)
		if err := c.repo.ApplyReputation(ctxDB, ev); err != nil {
			log.Printf("failed to apply reputation: %v", err)
		} else {
			log.Printf("applied reputation event %s for user %s (Δ=%d)", ev.EventID, ev.UserID, ev.Delta)
		}
		cancel()
	}
}
