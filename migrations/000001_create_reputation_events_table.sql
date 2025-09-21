CREATE TABLE IF NOT EXISTS reputation_events (
                                                 event_id    UUID PRIMARY KEY,
                                                 user_id     UUID NOT NULL,
                                                 delta       INTEGER NOT NULL,
                                                 reason      TEXT NOT NULL,
                                                 occurred_at TIMESTAMPTZ NOT NULL,
                                                 created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_rep_events_user_time
    ON reputation_events (user_id, occurred_at DESC);