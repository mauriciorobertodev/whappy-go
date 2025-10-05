package app

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/token"
)

// we can make this in http layer, bur i like to keep it here, em app layer to propagare another presenters if needed

var errorCodeMap = map[error]AppCode{
	// Instance errors
	instance.ErrInstanceNotFound:            CodeInstanceNotFound,
	instance.ErrInstanceBanned:              CodeInstanceBanned,
	instance.ErrInstanceAlreadyLoggedIn:     CodeInstanceAlreadyLoggedIn,
	instance.ErrInstanceIsPairing:           CodeInstanceIsPairing,
	instance.ErrInstanceAlreadyConnected:    CodeInstanceAlreadyConnected,
	instance.ErrInstanceIsConnecting:        CodeInstanceIsConnecting,
	instance.ErrInstanceNotLoggedIn:         CodeInstanceNotLoggedIn,
	instance.ErrInstanceNotConnected:        CodeInstanceNotConnected,
	instance.ErrInstanceNotPairing:          CodeInstanceNotPairing,
	instance.ErrInstanceAlreadyLoggedOut:    CodeInstanceAlreadyLoggedOut,
	instance.ErrInstanceAlreadyDisconnected: CodeInstanceAlreadyDisconnected,
	instance.ErrInstanceNoQrCode:            CodeInstanceNoQrCode,
	instance.ErrInstanceAlreadyPaired:       CodeInstanceAlreadyPaired,

	// Whatsapp errors
	whatsapp.ErrClientOutdated:            CodeFailPairingClientOutdated,
	whatsapp.ErrScannedWithoutMultiDevice: CodeFailPairingWithoutMultidevice,
	whatsapp.ErrTimeout:                   CodeFailPairingTimeout,

	// Token errors
	token.ErrInvalidToken: CodeInvalidToken,

	// File errors
	file.ErrFileNotFound:        CodeFileNotFound,
	file.ErrInvalidFile:         CodeInvalidFile,
	file.ErrFileTooLarge:        CodeFileTooLarge,
	file.ErrUnsupported:         CodeUnsupportedFileType,
	file.ErrCorruptedFile:       CodeCorruptedFile,
	file.ErrUploadFailed:        CodeFileUploadFailed,
	file.ErrDownloadFailed:      CodeFileDownloadFailed,
	file.ErrStorageFailed:       CodeFileStorageFailed,
	file.ErrFileUnreachable:     CodeFileUnreachable,
	file.ErrFileCannotBeImage:   CodeInvalidImage,
	file.ErrFileCannotBeVideo:   CodeInvalidVideo,
	file.ErrFileCannotBeAudio:   CodeInvalidAudio,
	file.ErrFileCannotBeVoice:   CodeInvalidVoice,
	file.ErrFileCannotBeDeleted: CodeFileCannotBeDeleted,
	file.ErrFileSourceEmpty:     CodeFileSourceEmpty,
}

func TranslateError(location string, err error) *AppError {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(*AppError); ok {
		return WrapLoc(location, appErr)
	}

	if code, exists := errorCodeMap[err]; exists {
		return NewAppError(location, code, err)
	}

	return NewUnknownError(location, err)
}
