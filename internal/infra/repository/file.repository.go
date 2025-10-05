package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/repository/models"
)

type FileRepository struct {
	db *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Insert(f *file.File) error {
	sqlFile, err := models.FromFileEntity(f)
	if err != nil {
		return err
	}

	_, err = r.db.NamedExec(`
		INSERT INTO files (
			id, name, mime, size, sha256, extension, 
			path, url,
			width, height, duration, pages,
			created_at, updated_at, instance_id, thumbnail_id
		) VALUES (
			:id, :name, :mime, :size, :sha256, :extension, 
			:path, :url, 
			:width, :height, :duration, :pages,
			:created_at, :updated_at, :instance_id, :thumbnail_id
		)
	`, sqlFile)
	return err
}

func (r *FileRepository) InsertMany(files []*file.File) error {
	tx := r.db.MustBegin()

	for _, f := range files {
		sqlFile, err := models.FromFileEntity(f)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.NamedExec(`
			INSERT INTO files (
				id, name, mime, size, sha256, extension, 
				path, url,
				width, height, duration, pages,
				created_at, updated_at, instance_id, thumbnail_id
			) VALUES (
				:id, :name, :mime, :size, :sha256, :extension, 
				:path, :url, 
				:width, :height, :duration, :pages,
				:created_at, :updated_at, :instance_id, :thumbnail_id
			)
		`, sqlFile)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *FileRepository) Update(f *file.File) error {
	sqlFile, err := models.FromFileEntity(f)
	if err != nil {
		return err
	}

	_, err = r.db.NamedExec(`
		UPDATE files SET
			name = :name,
			mime = :mime,
			extension = :extension,

			size = :size,
			sha256 = :sha256,
			
			path = :path,
			url = :url,
			
			width = :width,
			height = :height,
			duration = :duration,
			pages = :pages,

			updated_at = :updated_at,

			thumbnail_id = :thumbnail_id
		WHERE id = :id
	`, sqlFile)
	return err
}

func (r *FileRepository) Get(opts ...file.FileQueryOption) (*file.File, error) {
	params := &file.FileQueryOptions{}
	for _, opt := range opts {
		opt(params)
	}

	query := "SELECT f.*"

	if params.WithThumbnail {
		query += `,
			t.id           AS "thumbnail.id",
			t.name         AS "thumbnail.name",
			t.mime         AS "thumbnail.mime",
			t.extension    AS "thumbnail.extension",
			t.size         AS "thumbnail.size",
			t.sha256       AS "thumbnail.sha256",
			t.path         AS "thumbnail.path",
			t.url          AS "thumbnail.url",
			t.width        AS "thumbnail.width",
			t.height       AS "thumbnail.height",
			t.duration     AS "thumbnail.duration",
			t.pages        AS "thumbnail.pages",
			t.created_at   AS "thumbnail.created_at",
			t.updated_at   AS "thumbnail.updated_at",
			t.instance_id  AS "thumbnail.instance_id",
			t.thumbnail_id AS "thumbnail.thumbnail_id"
		`
	}

	query += " FROM files f"

	if params.WithThumbnail {
		query += " LEFT JOIN files t ON f.thumbnail_id = t.id"
	}

	query += " WHERE 1=1"

	if params.InstanceID != nil {
		query += " AND f.instance_id = :instance_id"
	}

	if params.ID != nil {
		query += " AND f.id = :id"
	}

	if params.Sha256 != nil {
		query += " AND f.sha256 = :sha256"
	}

	if params.HasThumbnail != nil {
		if *params.HasThumbnail {
			query += " AND thumbnail_id IS NOT NULL"
		} else {
			query += " AND thumbnail_id IS NULL"
		}
	}

	if params.OrderBy == "" {
		params.OrderBy = "f.created_at"
	}

	if params.SortBy == "" {
		params.SortBy = "DESC"
	}

	query += " LIMIT 1"

	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer nstmt.Close()

	var sqlFile models.SQLFile
	err = nstmt.Get(&sqlFile, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, file.ErrFileNotFound
		}
		return nil, err
	}

	f, err := sqlFile.ToEntity()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *FileRepository) List(opts ...file.FileQueryOption) ([]*file.File, error) {
	params := &file.FileQueryOptions{}
	for _, opt := range opts {
		opt(params)
	}

	query := "SELECT f.*"

	if params.WithThumbnail {
		query += `,
			t.id           AS "thumbnail.id",
			t.name         AS "thumbnail.name",
			t.mime         AS "thumbnail.mime",
			t.extension    AS "thumbnail.extension",
			t.size         AS "thumbnail.size",
			t.sha256       AS "thumbnail.sha256",
			t.path         AS "thumbnail.path",
			t.url          AS "thumbnail.url",
			t.width        AS "thumbnail.width",
			t.height       AS "thumbnail.height",
			t.duration     AS "thumbnail.duration",
			t.pages        AS "thumbnail.pages",
			t.created_at   AS "thumbnail.created_at",
			t.updated_at   AS "thumbnail.updated_at",
			t.instance_id  AS "thumbnail.instance_id",
			t.thumbnail_id AS "thumbnail.thumbnail_id"
		`
	}

	query += " FROM files f"

	if params.WithThumbnail {
		query += " LEFT JOIN files t ON f.thumbnail_id = t.id"
	}

	query += " WHERE 1=1"

	if params.InstanceID != nil {
		query += " AND f.instance_id = :instance_id"
	}

	if params.Cursor != nil {
		query += " AND f.created_at <= :cursor"
	}

	if params.ID != nil {
		query += " AND f.id = :id"
	}

	if params.HasThumbnail != nil {
		if *params.HasThumbnail {
			query += " AND f.thumbnail_id IS NOT NULL"
		} else {
			query += " AND f.thumbnail_id IS NULL"
		}
	}

	query += " ORDER BY f.created_at DESC"

	if params.Limit != nil {
		query += " LIMIT :limit"
	}

	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer nstmt.Close()

	var sqlFiles []models.SQLFile
	err = nstmt.Select(&sqlFiles, params)
	if err != nil {
		return nil, err
	}

	files := make([]*file.File, 0, len(sqlFiles))
	for _, sf := range sqlFiles {
		f, err := sf.ToEntity()
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (r *FileRepository) Delete(opts ...file.FileQueryOption) error {
	params := &file.FileQueryOptions{}
	for _, opt := range opts {
		opt(params)
	}

	query := "DELETE FROM files WHERE 1=1"

	if params.ID != nil {
		query += " AND id = :id"
	}

	if params.Sha256 != nil {
		query += " AND sha256 = :sha256"
	}

	if params.InstanceID != nil {
		query += " AND instance_id = :instance_id"
	}

	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	_, err = nstmt.Exec(params)
	return err
}

func (r *FileRepository) Count(opts ...file.FileQueryOption) (uint64, error) {
	params := &file.FileQueryOptions{}
	for _, opt := range opts {
		opt(params)
	}

	query := "SELECT COUNT(*) FROM files WHERE 1=1"

	if params.InstanceID != nil {
		query += " AND instance_id = :instance_id"
	}

	if params.HasThumbnail != nil {
		if *params.HasThumbnail {
			query += " AND thumbnail_id IS NOT NULL"
		} else {
			query += " AND thumbnail_id IS NULL"
		}
	}

	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return 0, err
	}
	defer nstmt.Close()

	var count uint64
	err = nstmt.Get(&count, params)
	if err != nil {
		return 0, err
	}

	return count, nil
}
