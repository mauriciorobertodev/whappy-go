package service

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/storage"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

type UploadService struct {
	fileService *FileService
	fileRepo    file.FileRepository
	storage     storage.Storage
	bus         events.EventBus
}

func NewUploadService(fileService *FileService, fileRepo file.FileRepository, storage storage.Storage, bus events.EventBus) *UploadService {
	return &UploadService{
		fileService: fileService,
		fileRepo:    fileRepo,
		storage:     storage,
		bus:         bus,
	}
}

func (s *UploadService) UploadWithStream(ctx context.Context, inst *instance.Instance, inp input.UploadFile) (*file.File, *app.AppError) {
	l := app.GetUploadServiceLogger()

	if s.storage == nil {
		l.Error("Global storage is not configured")
		return nil, app.NewAppError("upload service", app.GLOBAL_STORAGE_UNAVAILABLE, storage.ErrStorageNotConfigured)
	}

	l.Debug("Uploading file to storage")

	f, err := s.fileService.SaveStream(ctx, inp.Stream, inp.Metadata.Mime)
	if err != nil {
		l.Error("Error uploading file to storage", "error", err.Error())
		return nil, app.TranslateError("upload service", err)
	}

	f.UpdateMeta(inp.Metadata)
	f.InstanceID = &inst.ID

	if inp.ThumbnailID != nil && *inp.ThumbnailID != "" {
		thumbFile, err := s.fileRepo.Get(file.WhereID(*inp.ThumbnailID))
		if err != nil {
			if errors.Is(err, file.ErrFileNotFound) {
				l.Error("Thumbnail file not found", "thumbnail_id", *inp.ThumbnailID)
				return nil, app.NewAppError("upload service", app.CodeFileNotFound, file.ErrFileNotFound)
			}

			l.Error("Error getting thumbnail file from database", "error", err.Error())
			return nil, app.NewDatabaseError("upload service", err)
		}

		thumb, err := thumbFile.ToImageFile()
		if err != nil {
			l.Error("Error converting thumbnail file to image", "error", err.Error())
			return nil, app.NewAppError("upload service", app.CodeInvalidImage, file.ErrFileCannotBeImage)
		}

		f.Thumbnail = thumb
	}

	if f.HasThumbnail() {
		l.Debug("File has thumbnail", "thumbnail_id", f.Thumbnail.ID)
	}

	l.Info("File uploaded to storage successfully", "name", f.Name)

	err = s.fileRepo.Insert(f)
	if err != nil {
		l.Error("Error saving file to database", "error", err.Error())
		return nil, app.TranslateError("upload service", err)
	}

	go s.bus.Publish(f.EventUploaded(&inst.ID))

	l.Info("File saved to database successfully", "name", f.Name)

	return f, nil
}

func (s *UploadService) ListUploads(ctx context.Context, inst *instance.Instance, inp input.ListUploads) ([]*file.File, *string, *app.AppError) {
	l := app.GetUploadServiceLogger()

	l.Info("Listing files from database", "instance", inst.ID, "limit", inp.Limit, "cursor", inp.Cursor)

	inp.Normalize()

	var cursor *time.Time

	if inp.Cursor != nil {
		decodedBytes, err := base64.URLEncoding.DecodeString(*inp.Cursor)
		if err != nil {
			l.Error("Error decoding cursor", "error", err.Error())
			return nil, nil, app.NewAppError("upload service", app.CodeInvalidCursor, err)
		}

		decodedTimeString := string(decodedBytes)
		cursor = new(time.Time)
		*cursor, err = time.Parse(time.RFC3339, decodedTimeString)
		if err != nil {
			l.Error("Error parsing cursor", "error", err.Error())
			return nil, nil, app.NewAppError("upload service", app.CodeInvalidCursor, err)
		}
	}

	files, err := s.fileRepo.List(file.WithCursor(cursor, inp.Limit+1))
	if err != nil {
		l.Error("Error listing files from database", "error", err.Error())
		return nil, nil, app.TranslateError("upload service", err)
	}

	l.Info("Files listed from database successfully", "instance", inst.ID, "found", len(files))

	var nextCursorEncoded *string

	if len(files) > inp.Limit {
		lastFile := files[len(files)-1]
		cursorTimeString := lastFile.CreatedAt.Format(time.RFC3339)
		nextCursorBase64 := base64.URLEncoding.EncodeToString([]byte(cursorTimeString))
		nextCursorEncoded = &nextCursorBase64
		files = files[:len(files)-1]
	}

	return files, nextCursorEncoded, nil
}

func (s *UploadService) GetUpload(ctx context.Context, inst *instance.Instance, inp input.GetUpload) (*file.File, *app.AppError) {
	l := app.GetUploadServiceLogger()

	l.Info("Getting file from database", "file", inp.FileID)

	f, err := s.fileRepo.Get(file.WhereID(inp.FileID))
	if err != nil {
		if errors.Is(err, file.ErrFileNotFound) {
			l.Error("File not found", "file", inp.FileID)
			return nil, app.NewAppError("upload service", app.CodeFileNotFound, file.ErrFileNotFound)
		}

		l.Error("Error getting file from database", "error", err.Error())
		return nil, app.TranslateError("upload service", err)
	}

	l.Info("File retrieved from database successfully", "file", inp.FileID)

	return f, nil
}

func (s *UploadService) DeleteUpload(ctx context.Context, inst *instance.Instance, inp input.DeleteUpload) *app.AppError {
	l := app.GetUploadServiceLogger()

	if s.storage == nil {
		l.Error("Global storage is not configured")
		return app.NewAppError("upload service", app.GLOBAL_STORAGE_UNAVAILABLE, storage.ErrStorageNotConfigured)
	}

	l.Info("Deleting file from database", "file", inp.FileID)

	f, err := s.fileRepo.Get(file.WhereID(inp.FileID))
	if err != nil {
		if errors.Is(err, file.ErrFileNotFound) {
			l.Error("File not found", "file", inp.FileID)
			return app.NewAppError("upload service", app.CodeFileNotFound, file.ErrFileNotFound)
		}

		l.Error("Error getting file from database", "error", err.Error())
		return app.TranslateError("upload service", err)
	}

	if f.InstanceID == nil || *f.InstanceID != inst.ID {
		l.Error("File does not belong to this instance", "file", inp.FileID, "instance", inst.ID)
		return app.NewAppError("upload service", app.CodeFileNotFound, file.ErrFileNotFound)
	}

	err = s.fileRepo.Delete(file.WhereID(inp.FileID))
	if err != nil {
		l.Error("Error deleting file from database", "error", err.Error())
		return app.TranslateError("upload service", err)
	}

	l.Info("File deleted from database successfully", "file", inp.FileID)

	err = s.storage.Delete(ctx, f.Path)
	if err != nil {
		l.Error("Error deleting file from storage", "error", err.Error())
		return app.TranslateError("upload service", err)
	}

	l.Info("File deleted from storage successfully", "file", inp.FileID)

	go s.bus.Publish(f.EventDeleted(&inst.ID))

	return nil
}

func (s *UploadService) UpdateFileMetadata(ctx context.Context, inst *instance.Instance, inp input.UpdateUploadMetadata) (*file.File, *app.AppError) {
	l := app.GetUploadServiceLogger()

	l.Info("Updating file metadata in database", "file", inp.FileID)

	f, err := s.fileRepo.Get(file.WhereID(inp.FileID))
	if err != nil {
		l.Error("Error getting file from database", "error", err.Error())
		return nil, app.TranslateError("upload service", err)
	}

	if f == nil {
		l.Error("File not found", "file", inp.FileID)
		return nil, app.NewAppError("upload service", app.CodeFileNotFound, file.ErrFileNotFound)
	}

	if f.InstanceID == nil || *f.InstanceID != inst.ID {
		l.Error("File does not belong to this instance", "file", inp.FileID, "instance", inst.ID)
		return nil, app.NewAppError("upload service", app.CodeFileNotFound, file.ErrFileNotFound)
	}

	f.UpdateMeta(inp.Metadata)

	err = s.fileRepo.Update(f)
	if err != nil {
		l.Error("Error updating file in database", "error", err.Error())
		return nil, app.TranslateError("upload service", err)
	}

	go s.bus.Publish(f.EventUpdated(&inst.ID))
	l.Info("File metadata updated in database successfully", "file", inp.FileID)

	return f, nil
}

// TODO: download
