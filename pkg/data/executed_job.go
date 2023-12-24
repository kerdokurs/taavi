package data

import (
	"context"
	"time"
)

type ExecutedJob struct {
	Model

	RanAt   time.Time `db:"ran_at" json:"ran_at"`
	Content string    `db:"content" json:"content"`

	Error string `db:"error" json:"error"`

	JobID uint `db:"job_id" json:"-"`
}

func InsertExecutedJob(ctx context.Context, executedJob ExecutedJob) error {
	query := `
		INSERT INTO executed_jobs (ran_at, content, error, job_id)
		VALUES (:ran_at, :content, :error, :job_id)
	`
	if _, err := db.NamedExecContext(ctx, query, executedJob); err != nil {
		return err
	}

	return nil
}
