package file

import "errors"

var (
	ErrFileNotFound    = errors.New("file not found")
	ErrInvalidFile     = errors.New("invalid file")
	ErrFileTooLarge    = errors.New("file too large")
	ErrUnsupported     = errors.New("unsupported file type")
	ErrCorruptedFile   = errors.New("corrupted file")
	ErrUploadFailed    = errors.New("file upload failed")
	ErrDownloadFailed  = errors.New("file download failed")
	ErrStorageFailed   = errors.New("file storage failed")
	ErrFileUnreachable = errors.New("file unreachable")

	ErrFileCannotBeImage   = errors.New("file cannot be an image")
	ErrFileCannotBeVideo   = errors.New("file cannot be a video")
	ErrFileCannotBeAudio   = errors.New("file cannot be an audio")
	ErrFileCannotBeVoice   = errors.New("file cannot be a voice")
	ErrFileCannotBeDeleted = errors.New("file cannot be deleted")

	ErrFileSourceEmpty = errors.New("file source is empty")
	ErrInvalidFileID   = errors.New("invalid file ID")
)
