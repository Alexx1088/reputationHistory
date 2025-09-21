package repository

import (
	"context"
	"database/sql"

	kafkamodel "github.com/Alexx1088/reputationhistory/internal/kafka"
	"github.com/jackc/pgconn"
)

type Repo struct {
	DB *sql.DB
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}

func (r *Repo) ApplyReputation(ctx context.Context, ev kafkamodel.ReputationEntryEvent) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1) insert raw event (dedupe by PK event_id)
	_, err = tx.ExecContext(ctx,
		`INSERT INTO reputation_events(event_id, user_id, delta, reason, occurred_at)
         VALUES ($1,$2,$3,$4,$5)`,
		ev.EventID, ev.UserID, ev.Delta, ev.Reason, ev.OccurredAt)
	if isUniqueViolation(err) {
		return tx.Commit() // already processed
	}
	if err != nil {
		return err
	}

	// 2) (optional now) keep a running total
	_, err = tx.ExecContext(ctx,
		`INSERT INTO user_reputation(user_id, score)
           VALUES ($1,$2)
         ON CONFLICT (user_id)
           DO UPDATE SET score = user_reputation.score + EXCLUDED.score`,
		ev.UserID, ev.Delta)
	if err != nil {
		return err
	}

	return tx.Commit()
}
