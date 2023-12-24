package data

import "context"

type JobType string

const (
	Simple JobType = "SIMPLE"
	Random JobType = "RANDOM"
)

type Job struct {
	Model

	Type JobType `db:"job_type" json:"type"`

	StreamID string `db:"stream_id" json:"stream_id"`
	TopicID  string `db:"topic_id" json:"topic_id"`
	Content  string `db:"content" json:"content"`

	CronTime string `db:"cron_time" json:"cron_time"`

	Enabled bool `db:"enabled" json:"enabled"`

	Meta []JobMeta `json:"meta"`
}

func GetAllJobs(ctx context.Context) ([]Job, error) {
	query := `
		SELECT * FROM jobs
		WHERE deleted_at IS NULL
	`
	jobs := []Job{}
	if err := db.SelectContext(ctx, &jobs, query); err != nil {
		return nil, err
	}

	// TODO: Support job_meta

	return jobs, nil
}

func GetJobs(ctx context.Context) ([]Job, error) {
	query := `
		SELECT * FROM jobs
		WHERE
			deleted_at IS NULL
			AND enabled = true
			AND job_type = 'SIMPLE'
	`
	jobs := []Job{}
	if err := db.SelectContext(ctx, &jobs, query); err != nil {
		return nil, err
	}

	return jobs, nil
}

func GetRandomJobs(ctx context.Context) ([]Job, error) {
	query := `
		SELECT * FROM jobs
		WHERE
			deleted_at IS NULL
			AND enabled = true
			AND job_type = 'RANDOM'
	`
	jobs := []Job{}
	if _, err := db.NamedQueryContext(ctx, query, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func GetJob(ctx context.Context, jobID uint) (*Job, error) {
	query := `
		SELECT * FROM jobs
		WHERE
			deleted_at IS NULL
			AND id = $1
	`
	job := Job{}
	if err := db.GetContext(ctx, &job, query, jobID); err != nil {
		return nil, err
	}

	return &job, nil
}

func DeleteJob(ctx context.Context, jobID uint) error {
	var err error
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `
		UPDATE jobs
		SET deleted_at = date('now')
		WHERE id = $1
	`
	if _, err := db.ExecContext(ctx, query, jobID); err != nil {
		return err
	}
	if err := DeleteJobMetasTx(ctx, tx, jobID); err != nil {
		return err
	}

	return nil
}

func ToggleJobEnabled(ctx context.Context, jobID uint) (bool, error) {
	query := `
		UPDATE jobs
		SET enabled = NOT enabled
		WHERE id = $1
		RETURNING enabled
	`
	var enabled bool
	if err := db.GetContext(ctx, &enabled, query, jobID); err != nil {
		return false, err
	}

	return enabled, nil
}

func InsertJobWithMeta(ctx context.Context, job Job, metas []JobMeta) error {
	var err error

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `
		INSERT INTO jobs (job_type, stream_id, topic_id, content, cron_time)
		VALUES (:job_type, :stream_id, :topic_id, :content, :cron_time)
	`
	var jobID int64
	if res, err := tx.NamedExecContext(ctx, query, job); err != nil {
		return err
	} else {
		if jobID, err = res.LastInsertId(); err != nil {
			return err
		}
	}

	for i := range metas {
		metas[i].JobID = uint(jobID)
	}

	if err := InsertJobMetasTx(ctx, tx, metas); err != nil {
		return err
	}

	return nil
}
