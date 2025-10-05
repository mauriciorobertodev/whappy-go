package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/repository/models"
)

type InstanceRepository struct {
	db *sqlx.DB
}

func NewInstanceRepository(db *sqlx.DB) *InstanceRepository {
	return &InstanceRepository{db: db}
}

func (r *InstanceRepository) Insert(inst *instance.Instance) error {
	sqlInst, err := models.FromInstanceEntity(inst)
	if err != nil {
		return err
	}

	_, err = r.db.NamedExec(`
		INSERT INTO instances (
			id, name, phone, jid, lid, device, status,
			last_login_at, last_connected_at, banned_at, ban_expires_at, created_at, updated_at
		) VALUES (
			:id, :name, :phone, :jid, :lid, :device, :status,
			:last_login_at, :last_connected_at, :banned_at, :ban_expires_at, :created_at, :updated_at
		)
	`, sqlInst)
	return err
}

func (r *InstanceRepository) InsertMany(insts []*instance.Instance) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(`
		INSERT INTO instances (
			id, name, phone, jid, lid, device, status,
			last_login_at, last_connected_at, banned_at, ban_expires_at, created_at, updated_at
		) VALUES (
			:id, :name, :phone, :jid, :lid, :device, :status,
			:last_login_at, :last_connected_at, :banned_at, :ban_expires_at, :created_at, :updated_at
		)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, inst := range insts {
		sqlInst, err := models.FromInstanceEntity(inst)
		if err != nil {
			tx.Rollback()
			return err
		}

		if _, err := stmt.Exec(sqlInst); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *InstanceRepository) Update(inst *instance.Instance) error {
	sqlInst, err := models.FromInstanceEntity(inst)
	if err != nil {
		return err
	}

	_, err = r.db.NamedExec(`
		UPDATE instances SET
			name = :name,
			phone = :phone,
			jid = :jid,
			lid = :lid,
			device = :device,
			status = :status,

			last_login_at = :last_login_at,
			last_connected_at = :last_connected_at,
			banned_at = :banned_at,
			ban_expires_at = :ban_expires_at,
			updated_at = :updated_at
		WHERE id = :id
	`, sqlInst)
	return err
}

func (r *InstanceRepository) Get(opts ...instance.InstanceQueryOption) (*instance.Instance, error) {
	params := &instance.InstanceQueryOptions{}
	for _, opt := range opts {
		opt(params)
	}

	query := "SELECT * FROM instances WHERE 1=1"

	if params.ID != nil {
		query += " AND id = :id"
	}

	if params.Phone != nil {
		query += " AND phone = :phone"
	}

	if params.JID != nil {
		query += " AND jid = :jid"
	}

	if params.LID != nil {
		query += " AND lid = :lid"
	}

	if params.Name != nil {
		query += " AND name = :name"
	}

	if params.Device != nil {
		query += " AND device = :device"
	}

	if params.Status != nil {
		query += " AND status = :status"
	}

	if params.OrderBy == "" {
		params.OrderBy = "created_at"
	}

	if params.SortBy == "" {
		params.SortBy = "DESC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", params.OrderBy, params.SortBy)

	query += " LIMIT 1"

	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer nstmt.Close()

	var row models.SQLInstance
	err = nstmt.Get(&row, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, instance.ErrInstanceNotFound
		}
		return nil, err
	}

	i, err := row.ToEntity()
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (r *InstanceRepository) List(opts ...instance.InstanceQueryOption) ([]*instance.Instance, error) {
	params := &instance.InstanceQueryOptions{}
	for _, opt := range opts {
		opt(params)
	}

	query := "SELECT * FROM instances WHERE 1=1"

	if params.ID != nil {
		query += " AND id = :id"
	}

	if params.Phone != nil {
		query += " AND phone = :phone"
	}

	if params.Status != nil {
		query += " AND status = :status"
	}

	if params.OrderBy == "" {
		params.OrderBy = "created_at"
	}

	if params.SortBy == "" {
		params.SortBy = "DESC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", params.OrderBy, params.SortBy)

	if params.Limit != nil {
		query += " LIMIT :limit"
	}

	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer nstmt.Close()

	var rows []models.SQLInstance
	err = nstmt.Select(&rows, params)
	if err != nil {
		return nil, err
	}

	var instances []*instance.Instance
	for _, row := range rows {
		i, err := row.ToEntity()
		if err != nil {
			return nil, err
		}
		instances = append(instances, i)
	}

	return instances, nil
}

func (r *InstanceRepository) Delete(opts ...instance.InstanceQueryOption) error {
	params := &instance.InstanceQueryOptions{}
	for _, opt := range opts {
		opt(params)
	}

	query := "DELETE FROM instances WHERE 1=1"

	if params.ID != nil {
		query += " AND id = :id"
	}

	if params.Phone != nil {
		query += " AND phone = :phone"
	}

	if params.Status != nil {
		query += " AND status = :status"
	}

	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	_, err = nstmt.Exec(params)
	return err
}
