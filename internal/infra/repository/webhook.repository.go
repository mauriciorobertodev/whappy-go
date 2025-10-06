package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/repository/models"
)

type WebhookRepository struct {
	db *sqlx.DB
}

func NewWebhookRepository(db *sqlx.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) Insert(w *webhook.Webhook) error {
	sqlWebhook, err := models.FromWebhookEntity(w)
	if err != nil {
		return err
	}

	_, err = r.db.NamedExec(`
		INSERT INTO webhooks (
			id, secret, events, url, active, instance_id, created_at, updated_at
		) VALUES (
			:id, :secret, :events, :url, :active,:instance_id, :created_at, :updated_at
		)
	`, sqlWebhook)
	return err
}

func (r *WebhookRepository) InsertMany(webhooks []*webhook.Webhook) error {
	tx := r.db.MustBegin()

	sqlWebhooks := make([]*models.SQLWebhook, len(webhooks))
	for i, w := range webhooks {
		sqlWebhook, err := models.FromWebhookEntity(w)
		if err != nil {
			return err
		}
		sqlWebhooks[i] = sqlWebhook
	}

	_, err := tx.NamedExec(`
		INSERT INTO webhooks (
			id, secret, events, url, active, instance_id, created_at, updated_at
		) VALUES (
			:id, :secret, :events, :url, :active, :instance_id, :created_at, :updated_at
		)
	`, sqlWebhooks)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *WebhookRepository) Update(w *webhook.Webhook) error {
	sqlWebhook, err := models.FromWebhookEntity(w)
	if err != nil {
		return err
	}

	_, err = r.db.NamedExec(`
		UPDATE webhooks SET
			secret = :secret,
			events = :events,
			url = :url,
			active = :active,
			instance_id = :instance_id,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE id = :id
	`, sqlWebhook)
	return err
}

func (r *WebhookRepository) Get(opts ...webhook.WebhookQueryOption) (*webhook.Webhook, error) {
	queryOptions := &webhook.WebhookQueryOptions{
		OrderBy: "created_at",
		SortBy:  webhook.SortByDesc,
	}

	for _, opt := range opts {
		opt(queryOptions)
	}

	query := `SELECT * FROM webhooks WHERE 1=1`
	args := map[string]interface{}{}

	if queryOptions.ID != nil {
		query += " AND id = :id"
		args["id"] = *queryOptions.ID
	}
	if queryOptions.Active != nil {
		query += " AND active = :active"
		args["active"] = *queryOptions.Active
	}
	if queryOptions.URL != nil {
		query += " AND url = :url"
		args["url"] = *queryOptions.URL
	}
	if queryOptions.InstanceID != nil {
		query += " AND instance_id = :instance_id"
		args["instance_id"] = *queryOptions.InstanceID
	}

	query += " ORDER BY " + queryOptions.OrderBy + " " + string(queryOptions.SortBy)
	query += " LIMIT 1"

	var sqlWebhook models.SQLWebhook
	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	err = nstmt.Get(&sqlWebhook, args)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return sqlWebhook.ToEntity(), nil
}

func (r *WebhookRepository) List(opts ...webhook.WebhookQueryOption) ([]*webhook.Webhook, error) {
	queryOptions := &webhook.WebhookQueryOptions{
		OrderBy: "created_at",
		SortBy:  webhook.SortByDesc,
	}
	for _, opt := range opts {
		opt(queryOptions)
	}

	query := `SELECT * FROM webhooks WHERE 1=1`
	args := map[string]interface{}{}

	if queryOptions.ID != nil {
		query += " AND id = :id"
		args["id"] = *queryOptions.ID
	}
	if queryOptions.Active != nil {
		query += " AND active = :active"
		args["active"] = *queryOptions.Active
	}
	if queryOptions.URL != nil {
		query += " AND url = :url"
		args["url"] = *queryOptions.URL
	}
	if queryOptions.InstanceID != nil {
		query += " AND instance_id = :instance_id"
		args["instance_id"] = *queryOptions.InstanceID
	}

	query += " ORDER BY " + queryOptions.OrderBy + " " + string(queryOptions.SortBy)
	if queryOptions.Limit != nil {
		query += " LIMIT :limit"
		args["limit"] = *queryOptions.Limit
	}

	var sqlWebhooks []models.SQLWebhook
	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	err = nstmt.Select(&sqlWebhooks, args)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webhook.ErrNotFound
		}
		return nil, err
	}

	webhooks := make([]*webhook.Webhook, len(sqlWebhooks))
	for i, sqlWebhook := range sqlWebhooks {
		webhooks[i] = sqlWebhook.ToEntity()
	}

	return webhooks, nil
}

func (r *WebhookRepository) Delete(opts ...webhook.WebhookQueryOption) error {
	queryOptions := &webhook.WebhookQueryOptions{}
	for _, opt := range opts {
		opt(queryOptions)
	}

	query := `DELETE FROM webhooks WHERE 1=1`
	args := map[string]interface{}{}

	if queryOptions.ID != nil {
		query += " AND id = :id"
		args["id"] = *queryOptions.ID
	}
	if queryOptions.Active != nil {
		query += " AND active = :active"
		args["active"] = *queryOptions.Active
	}
	if queryOptions.URL != nil {
		query += " AND url = :url"
		args["url"] = *queryOptions.URL
	}
	if queryOptions.InstanceID != nil {
		query += " AND instance_id = :instance_id"
		args["instance_id"] = *queryOptions.InstanceID
	}

	_, err := r.db.NamedExec(query, args)
	return err
}

func (r *WebhookRepository) Count(opts ...webhook.WebhookQueryOption) (uint64, error) {
	queryOptions := &webhook.WebhookQueryOptions{}
	for _, opt := range opts {
		opt(queryOptions)
	}

	query := `SELECT COUNT(*) FROM webhooks WHERE 1=1`
	args := map[string]interface{}{}

	if queryOptions.ID != nil {
		query += " AND id = :id"
		args["id"] = *queryOptions.ID
	}
	if queryOptions.Active != nil {
		query += " AND active = :active"
		args["active"] = *queryOptions.Active
	}
	if queryOptions.URL != nil {
		query += " AND url = :url"
		args["url"] = *queryOptions.URL
	}
	if queryOptions.InstanceID != nil {
		query += " AND instance_id = :instance_id"
		args["instance_id"] = *queryOptions.InstanceID
	}

	var count uint64
	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return 0, err
	}
	err = nstmt.Get(&count, args)
	if err != nil {
		return 0, err
	}

	return count, nil
}
