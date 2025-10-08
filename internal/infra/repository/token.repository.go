package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/token"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/repository/models"
)

type TokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) FindByID(id string) (*token.Token, error) {
	var sqlToken models.SQLToken
	err := r.db.Get(&sqlToken, `SELECT * FROM tokens WHERE id=$1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return sqlToken.ToEntity(), nil
}
func (r *TokenRepository) FindByInstanceID(instanceID string) ([]*token.Token, error) {
	var sqlTokens []models.SQLToken
	err := r.db.Select(&sqlTokens, `SELECT * FROM tokens WHERE instance_id=$1`, instanceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*token.Token{}, nil
		}
		return nil, err
	}

	tokens := make([]*token.Token, len(sqlTokens))
	for i, sqlToken := range sqlTokens {
		tokens[i] = sqlToken.ToEntity()
	}

	return tokens, nil
}

func (r *TokenRepository) Insert(t *token.Token) error {
	sqlToken, err := models.FromTokenEntity(t)
	if err != nil {
		return err
	}

	_, err = r.db.NamedExec(`
		INSERT INTO tokens (
			id, instance_id, token_hash, created_at, updated_at
		) VALUES (
			:id, :instance_id, :token_hash, :created_at, :updated_at
		)
	`, sqlToken)
	return err
}

func (r *TokenRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM tokens WHERE id=$1`, id)
	return err
}

func (r *TokenRepository) Count(opts ...token.TokenQueryOption) int {
	params := &token.TokenQueryOptions{
		OrderBy: "created_at",
		SortBy:  "DESC",
	}
	for _, o := range opts {
		o(params)
	}

	query := "SELECT COUNT(*) FROM tokens WHERE 1=1"

	if params.ID != nil {
		query += " AND id = :id"
	}

	if params.InstanceID != nil {
		query += " AND instance_id = :instance_id"
	}

	var count int
	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return 0
	}
	defer nstmt.Close()

	err = nstmt.Get(&count, params)
	if err != nil {
		return 0
	}

	return count
}
