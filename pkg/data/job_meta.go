package data

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type JobMeta struct {
	Model

	Key   string `db:"key" json:"key"`
	Value string `db:"value" json:"value"`

	JobID uint `db:"job_id" json:"-"`
}

func GetJobMeta(ctx context.Context, job *Job, key string) (*JobMeta, error) {
	query := `
		SELECT * FROM job_meta
		WHERE
			deleted_at IS NULL
			AND job_id = $1
			AND key = $2
	`
	meta := JobMeta{}
	if err := db.GetContext(ctx, &meta, query, job.ID, key); err != nil {
		return nil, err
	}

	return &meta, nil
}

func GetJobMetas(ctx context.Context, jobID uint) ([]JobMeta, error) {
	query := `
		SELECT * FROM job_meta
		WHERE
			deleted_at IS NULL
			AND job_id = $1
	`
	metas := []JobMeta{}
	if err := db.SelectContext(ctx, &metas, query, jobID); err != nil {
		return nil, err
	}

	return metas, nil
}

func GetAllJobMetas(ctx context.Context, jobIDs []uint) (map[uint][]JobMeta, error) {
	query := `
		SELECT * FROM job_meta
		WHERE
			deleted_at IS NULL
			AND job_id IN (?)
	`
	metas := []JobMeta{}
	query, args, err := sqlx.In(query, jobIDs)
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	if err := db.SelectContext(ctx, &metas, query, args...); err != nil {
		return nil, err
	}

	jobMetas := map[uint][]JobMeta{}
	for _, meta := range metas {
		jobMetas[meta.JobID] = append(jobMetas[meta.JobID], meta)
	}

	return jobMetas, nil
}

func DeleteJobMetasTx(ctx context.Context, tx *sqlx.Tx, jobID uint) error {
	query := `
		UPDATE job_meta
		SET deleted_at = date('now')
		WHERE job_id = $1
	`
	if _, err := tx.ExecContext(ctx, query, jobID); err != nil {
		return err
	}

	return nil
}

func InsertJobMeta(ctx context.Context, meta *JobMeta) error {
	query := `
		INSERT INTO job_meta (job_id, key, value)
		VALUES (:job_id, :key, :value)
	`
	if _, err := db.NamedExecContext(ctx, query, meta); err != nil {
		return err
	}

	return nil
}

func InsertJobMetasTx(ctx context.Context, tx *sqlx.Tx, metas []JobMeta) error {
	query := `
		INSERT INTO job_meta (job_id, key, value)
		VALUES (:job_id, :key, :value)
	`
	if _, err := tx.NamedExecContext(ctx, query, metas); err != nil {
		return err
	}

	return nil
}
