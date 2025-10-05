package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/storage"
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/message"
	c "github.com/mauriciorobertodev/whappy-go/internal/infra/cache"
)

type MessageService struct {
	whatsapp           whatsapp.WhatsAppGateway
	storage            storage.Storage
	fileService        *FileService
	cache              cache.Cache
	cacheFileUploadTTL time.Duration
}

func NewMessageService(whatsapp whatsapp.WhatsAppGateway, storage storage.Storage, fileService *FileService, cache cache.Cache, cacheFileUploadTTL time.Duration) *MessageService {
	return &MessageService{
		whatsapp,
		storage,
		fileService,
		cache,
		cacheFileUploadTTL,
	}
}

func (s *MessageService) GetMessageIDs(ctx context.Context, inst *instance.Instance, inp input.GenerateMessageIDs) ([]string, error) {
	l := app.GetMessageServiceLogger()

	l.Debug("Generating message ID", "instance", inst.ID)

	ids := make([]string, inp.Quantity)

	for i := 0; i < inp.Quantity; i++ {
		id, err := s.whatsapp.GenerateMessageID(ctx, inst)
		if err != nil {
			l.Error("Error generating message ID", "error", err)
			return nil, app.TranslateError("message service", err)
		}
		ids[i] = id
	}

	l.Info("Message IDs generated successfully", "quantity", inp.Quantity, "instance", inst.ID)
	return ids, nil
}

func (s *MessageService) SendTextMessage(ctx context.Context, inst *instance.Instance, inp input.SendTextMessageInput) (*message.Message, *app.AppError) {
	l := app.GetMessageServiceLogger()

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("message service", err)
	}

	l.Debug("Sending text message", "instance", inst.ID, "phone", inst.Phone, "chat", inp.To)

	content := message.NewTextContent(inp.Text, inp.Mentions)
	message := message.NewMessage(inp.ID, inst.JID, inp.To, content, &inst.ID, inp.Expiration, true)
	message, err := s.whatsapp.SendTextMessage(ctx, inst, message)
	if err != nil {
		l.Error("Error sending text message", "error", err)
		return nil, app.TranslateError("message service", err)
	}

	return message, nil
}

func (s *MessageService) SendImageMessage(ctx context.Context, inst *instance.Instance, inp input.SendImageMessageInput) (*message.Message, error) {
	l := app.GetMessageServiceLogger()

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("message service", err)
	}

	l.Debug("Sending image message", "instance", inst.ID, "phone", inst.Phone, "chat", inp.To)

	var imageFile *file.ImageFile

	useCache := inp.Cache == nil || *inp.Cache
	source256 := sha256.Sum256([]byte(inp.Image))
	cacheKey := cache.CacheKeyFileUploadPrefix + hex.EncodeToString(source256[:])

	if useCache {
		cachedFile := s.getFileFromCache(ctx, cacheKey)
		if cachedFile != nil {
			imageFile, _ = cachedFile.ToImageFile()
		}
	}

	if imageFile == nil {
		loadedFile, stream, err := s.fileService.LoadFrom(ctx, inp.Image)
		if err != nil {
			l.Error("Error loading image file", "error", err)
			return nil, app.TranslateError("message service", err)
		}
		loadedFile.UpdateMeta(file.Metadata{
			Name:   inp.Name,
			Mime:   inp.Mime,
			Width:  inp.Width,
			Height: inp.Height,
		})

		imageFile, err = loadedFile.ToImageFile()
		if err != nil {
			l.Error("Error converting to image file", "error", err)
			return nil, app.NewAppError("message service", app.CodeInvalidImage, err)
		}

		uploadedFile, err := s.uploadFile(ctx, inst, inp.Image, stream, whatsapp.MediaImage, imageFile.Mime)
		if err != nil {
			return nil, app.NewAppError("message service", app.CodeFailWhatsappUpload, err)
		}

		imageFile.URL = uploadedFile.URL
		imageFile.DirectPath = uploadedFile.DirectPath
		imageFile.Sha256 = uploadedFile.Sha256
		imageFile.Size = uploadedFile.Size
		imageFile.Sha256Enc = uploadedFile.Sha256Enc
		imageFile.MediaKey = uploadedFile.MediaKey

		if useCache {
			l.Debug("Caching uploaded file", "cacheKey", cacheKey)
			c.Set(s.cache, cacheKey, imageFile, s.cacheFileUploadTTL)
		}
	}

	thumbnail := inp.Thumbnail
	if thumbnail == nil && imageFile.HasThumbnail() {
		l.Debug("Using video thumbnail as message thumbnail")
		thumbnail = &imageFile.Thumbnail.ID
	}

	var err error
	if thumbnail != nil {
		thumbnail, err = s.getThumbnail(ctx, *thumbnail, useCache)
		if err != nil {
			return nil, app.NewAppError("message service", app.CodeInvalidThumbnail, err)
		}
	}

	content := message.NewImageContent(imageFile, thumbnail, inp.Caption, inp.Mentions, inp.ViewOnce)
	message := message.NewMessage(inp.ID, inst.JID, inp.To, content, &inst.ID, inp.Expiration, true)

	l.Debug("Sending image message", "instance", inst.ID, "chat", inp.To)
	msg, err := s.whatsapp.SendImageMessage(ctx, inst, message)
	if err != nil {
		l.Error("Error sending image message", "error", err)
		return nil, app.TranslateError("message service", err)
	}

	l.Info("Image message sent successfully", "id", msg.ID, "instance", *msg.InstanceID, "sender", msg.Sender, "chat", msg.Chat)
	return msg, nil
}

func (s *MessageService) SendVideoMessage(ctx context.Context, inst *instance.Instance, inp input.SendVideoMessageInput) (*message.Message, error) {
	l := app.GetMessageServiceLogger()

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("message service", err)
	}

	l.Debug("Sending video message", "instance", inst.ID, "phone", inst.Phone, "chat", inp.To)

	var videoFile *file.VideoFile

	useCache := inp.Cache == nil || *inp.Cache
	source256 := sha256.Sum256([]byte(inp.Video))
	cacheKey := cache.CacheKeyFileUploadPrefix + hex.EncodeToString(source256[:])

	if useCache {
		cachedFile := s.getFileFromCache(ctx, cacheKey)
		if cachedFile != nil {
			videoFile, _ = cachedFile.ToVideoFile()
		}
	}

	if videoFile == nil {
		loadedFile, stream, err := s.fileService.LoadFrom(ctx, inp.Video)
		if err != nil {
			l.Error("Error loading video file", "error", err)
			return nil, app.TranslateError("message service", err)
		}
		loadedFile.UpdateMeta(file.Metadata{
			Name:   inp.Name,
			Mime:   inp.Mime,
			Width:  inp.Width,
			Height: inp.Height,
		})

		videoFile, err = loadedFile.ToVideoFile()
		if err != nil {
			l.Error("Error converting to image file", "error", err)
			return nil, app.NewAppError("message service", app.CodeInvalidImage, err)
		}

		uploadedFile, err := s.uploadFile(ctx, inst, inp.Video, stream, whatsapp.MediaVideo, videoFile.Mime)
		if err != nil {
			return nil, app.NewAppError("message service", app.CodeFailWhatsappUpload, err)
		}

		videoFile.URL = uploadedFile.URL
		videoFile.DirectPath = uploadedFile.DirectPath
		videoFile.Sha256 = uploadedFile.Sha256
		videoFile.Size = uploadedFile.Size
		videoFile.Sha256Enc = uploadedFile.Sha256Enc
		videoFile.MediaKey = uploadedFile.MediaKey

		if useCache {
			l.Debug("Caching uploaded file", "cacheKey", cacheKey)
			c.Set(s.cache, cacheKey, videoFile, s.cacheFileUploadTTL)
		}
	}

	thumbnail := inp.Thumbnail
	if thumbnail == nil && videoFile.HasThumbnail() {
		l.Debug("Using video thumbnail as message thumbnail")
		thumbnail = &videoFile.Thumbnail.ID
	}

	var err error
	if thumbnail != nil {
		thumbnail, err = s.getThumbnail(ctx, *thumbnail, useCache)
		if err != nil {
			return nil, app.NewAppError("message service", app.CodeInvalidThumbnail, err)
		}
	}

	content := message.NewVideoContent(*videoFile, thumbnail, inp.Caption, inp.Mentions, inp.ViewOnce)
	message := message.NewMessage(inp.ID, inst.JID, inp.To, content, &inst.ID, inp.Expiration, true)

	l.Debug("Sending video message", "instance", inst.ID, "chat", inp.To)
	msg, err := s.whatsapp.SendVideoMessage(ctx, inst, message)
	if err != nil {
		l.Error("Error sending video message", "error", err)
		return nil, app.TranslateError("message service", err)
	}

	l.Info("Video message sent successfully", "id", msg.ID, "instance", *msg.InstanceID, "sender", msg.Sender, "chat", msg.Chat)
	return msg, nil
}

func (s *MessageService) SendAudioMessage(ctx context.Context, inst *instance.Instance, inp input.SendAudioMessageInput) (*message.Message, error) {
	l := app.GetMessageServiceLogger()

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("message service", err)
	}

	l.Debug("Sending audio message", "instance", inst.ID, "phone", inst.Phone, "chat", inp.To)

	var audioFile *file.AudioFile

	useCache := inp.Cache == nil || *inp.Cache
	source256 := sha256.Sum256([]byte(inp.Audio))
	cacheKey := cache.CacheKeyFileUploadPrefix + hex.EncodeToString(source256[:])

	if useCache {
		cachedFile := s.getFileFromCache(ctx, cacheKey)
		if cachedFile != nil {
			audioFile, _ = cachedFile.ToAudioFile()
		}
	}

	if audioFile == nil {
		loadedFile, stream, err := s.fileService.LoadFrom(ctx, inp.Audio)
		if err != nil {
			l.Error("Error loading audio file", "error", err)
			return nil, app.TranslateError("message service", err)
		}
		loadedFile.UpdateMeta(file.Metadata{
			Name:     inp.Name,
			Mime:     inp.Mime,
			Duration: inp.Duration,
		})

		audioFile, err = loadedFile.ToAudioFile()
		if err != nil {
			l.Error("Error converting to image file", "error", err)
			return nil, app.NewAppError("message service", app.CodeInvalidImage, err)
		}

		uploadedFile, err := s.uploadFile(ctx, inst, inp.Audio, stream, whatsapp.MediaAudio, audioFile.Mime)
		if err != nil {
			return nil, app.NewAppError("message service", app.CodeFailWhatsappUpload, err)
		}

		audioFile.URL = uploadedFile.URL
		audioFile.DirectPath = uploadedFile.DirectPath
		audioFile.Sha256 = uploadedFile.Sha256
		audioFile.Size = uploadedFile.Size
		audioFile.Sha256Enc = uploadedFile.Sha256Enc
		audioFile.MediaKey = uploadedFile.MediaKey

		if useCache {
			l.Debug("Caching uploaded file", "cacheKey", cacheKey)
			c.Set(s.cache, cacheKey, audioFile, s.cacheFileUploadTTL)
		}
	}

	content := message.NewAudioContent(*audioFile)
	message := message.NewMessage(inp.ID, inst.JID, inp.To, content, &inst.ID, inp.Expiration, true)

	l.Debug("Sending audio message", "instance", inst.ID, "chat", inp.To)
	msg, err := s.whatsapp.SendAudioMessage(ctx, inst, message)
	if err != nil {
		l.Error("Error sending audio message", "error", err)
		return nil, app.TranslateError("message service", err)
	}

	l.Info("Audio message sent successfully", "id", msg.ID, "instance", *msg.InstanceID, "sender", msg.Sender, "chat", msg.Chat)
	return msg, nil
}

func (s *MessageService) SendVoiceMessage(ctx context.Context, inst *instance.Instance, inp input.SendVoiceMessageInput) (*message.Message, error) {
	l := app.GetMessageServiceLogger()

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("message service", err)
	}

	l.Debug("Sending voice message", "instance", inst.ID, "phone", inst.Phone, "chat", inp.To)

	var voiceFile *file.VoiceFile

	useCache := inp.Cache == nil || *inp.Cache
	source256 := sha256.Sum256([]byte(inp.Voice))
	cacheKey := cache.CacheKeyFileUploadPrefix + hex.EncodeToString(source256[:])

	if useCache {
		cachedFile := s.getFileFromCache(ctx, cacheKey)
		if cachedFile != nil {
			voiceFile, _ = cachedFile.ToVoiceFile()
		}
	}

	if voiceFile == nil {
		loadedFile, stream, err := s.fileService.LoadFrom(ctx, inp.Voice)
		if err != nil {
			l.Error("Error loading voice file", "error", err)
			return nil, app.TranslateError("message service", err)
		}
		loadedFile.UpdateMeta(file.Metadata{
			Name:     inp.Name,
			Mime:     inp.Mime,
			Duration: inp.Duration,
		})

		voiceFile, err = loadedFile.ToVoiceFile()
		if err != nil {
			l.Error("Error converting to voice file", "error", err)
			return nil, app.NewAppError("message service", app.CodeInvalidImage, err)
		}

		uploadedFile, err := s.uploadFile(ctx, inst, inp.Voice, stream, whatsapp.MediaAudio, voiceFile.Mime)
		if err != nil {
			return nil, app.NewAppError("message service", app.CodeFailWhatsappUpload, err)
		}

		voiceFile.URL = uploadedFile.URL
		voiceFile.DirectPath = uploadedFile.DirectPath
		voiceFile.Sha256 = uploadedFile.Sha256
		voiceFile.Size = uploadedFile.Size
		voiceFile.Sha256Enc = uploadedFile.Sha256Enc
		voiceFile.MediaKey = uploadedFile.MediaKey

		if useCache {
			l.Debug("Caching uploaded file", "cacheKey", cacheKey)
			c.Set(s.cache, cacheKey, voiceFile, s.cacheFileUploadTTL)
		}
	}

	content := message.NewVoiceContent(*voiceFile, inp.ViewOnce)
	message := message.NewMessage(inp.ID, inst.JID, inp.To, content, &inst.ID, inp.Expiration, true)

	l.Debug("Sending voice message", "instance", inst.ID, "chat", inp.To)
	msg, err := s.whatsapp.SendVoiceMessage(ctx, inst, message)
	if err != nil {
		l.Error("Error sending voice message", "error", err)
		return nil, app.TranslateError("message service", err)
	}

	l.Info("Voice message sent successfully", "id", msg.ID, "instance", *msg.InstanceID, "sender", msg.Sender, "chat", msg.Chat)
	return msg, nil
}

func (s *MessageService) SendDocumentMessage(ctx context.Context, inst *instance.Instance, inp input.SendDocumentMessageInput) (*message.Message, error) {
	l := app.GetMessageServiceLogger()

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("message service", err)
	}

	l.Debug("Sending document message", "instance", inst.ID, "chat", inp.To)

	var docFile *file.File

	useCache := inp.Cache == nil || *inp.Cache
	source256 := sha256.Sum256([]byte(inp.Document))
	cacheKey := cache.CacheKeyFileUploadPrefix + hex.EncodeToString(source256[:])

	if useCache {
		cachedFile := s.getFileFromCache(ctx, cacheKey)
		if cachedFile != nil {
			docFile = cachedFile
		}
	}

	if docFile == nil {
		loadedFile, stream, err := s.fileService.LoadFrom(ctx, inp.Document)
		if err != nil {
			l.Error("Error loading document file", "error", err)
			return nil, app.TranslateError("message service", err)
		}
		loadedFile.UpdateMeta(file.Metadata{
			Name:  inp.Name,
			Mime:  inp.Mime,
			Pages: inp.Pages,
		})

		docFile = loadedFile

		uploadedFile, err := s.uploadFile(ctx, inst, inp.Document, stream, whatsapp.MediaDocument, docFile.Mime)
		if err != nil {
			return nil, app.NewAppError("message service", app.CodeFailWhatsappUpload, err)
		}

		docFile.URL = uploadedFile.URL
		docFile.DirectPath = uploadedFile.DirectPath
		docFile.Sha256 = uploadedFile.Sha256
		docFile.Size = uploadedFile.Size
		docFile.Sha256Enc = uploadedFile.Sha256Enc
		docFile.MediaKey = uploadedFile.MediaKey

		if useCache {
			l.Debug("Caching uploaded file", "cacheKey", cacheKey)
			c.Set(s.cache, cacheKey, docFile, s.cacheFileUploadTTL)
		}
	}

	thumbnail := inp.Thumbnail
	if thumbnail == nil && docFile.HasThumbnail() {
		l.Debug("Using document thumbnail as message thumbnail")
		thumbnail = &docFile.Thumbnail.ID
	}

	var err error
	if thumbnail != nil {
		thumbnail, err = s.getThumbnail(ctx, *thumbnail, useCache)
		if err != nil {
			return nil, app.NewAppError("message service", app.CodeInvalidThumbnail, err)
		}
	}

	content := message.NewDocumentContent(*docFile, thumbnail, inp.Caption, inp.Mentions)
	message := message.NewMessage(inp.ID, inst.JID, inp.To, content, &inst.ID, inp.Expiration, true)

	l.Debug("Sending document message", "instance", inst.ID, "chat", inp.To)
	msg, err := s.whatsapp.SendDocumentMessage(ctx, inst, message)
	if err != nil {
		l.Error("Error sending document message", "error", err)
		return nil, app.TranslateError("message service", err)
	}

	l.Info("Document message sent successfully", "id", msg.ID, "instance", *msg.InstanceID, "sender", msg.Sender, "chat", msg.Chat)
	return msg, nil
}

func (s *MessageService) MarkMessagesAsRead(ctx context.Context, inst *instance.Instance, inp input.ReadMessagesInput) *app.AppError {
	l := app.GetMessageServiceLogger()

	if err := inp.Validate(); err != nil {
		return app.TranslateError("message service", err)
	}

	l.Debug("Marking messages as read", "instance", inst.ID, "chat", inp.Chat, "sender", inp.Sender, "ids", inp.IDs)

	err := s.whatsapp.ReadMessages(ctx, inst, inp.Chat, inp.IDs, inp.Sender)
	if err != nil {
		l.Error("Error marking messages as read", "error", err)
		return app.TranslateError("message service", err)
	}

	l.Info("Messages marked as read successfully", "instance", inst.ID, "chat", inp.Chat, "sender", inp.Sender, "ids", inp.IDs)
	return nil
}

func (s *MessageService) SendReaction(ctx context.Context, inst *instance.Instance, inp input.SendReactionInput) (*message.Message, *app.AppError) {
	l := app.GetMessageServiceLogger()

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("message service", err)
	}

	l.Debug("Sending reaction", "instance", inst.ID, "chat", inp.To, "message", inp.Message, "emoji", inp.Emoji)

	content := message.NewReactionContent(inp.Emoji, inp.Message)
	message := message.NewMessage(nil, inst.JID, inp.To, content, &inst.ID, nil, true)

	l.Debug("Sending reaction message", "instance", inst.ID, "chat", inp.To)

	msg, err := s.whatsapp.SendReaction(ctx, inst, message)
	if err != nil {
		l.Error("Error sending reaction", "error", err)
		return nil, app.TranslateError("message service", err)
	}

	l.Info("Reaction sent successfully", "id", msg.ID, "instance", *msg.InstanceID, "sender", msg.Sender, "chat", msg.Chat)
	return msg, nil
}

func (s *MessageService) getThumbnailFromCache(ctx context.Context, cacheKey string) *string {
	l := app.GetMessageServiceLogger()

	cachedThumbnail, err := c.Get[string](s.cache, cacheKey)
	if err == nil {
		l.Debug("Found cached thumbnail", "cacheKey", cacheKey)
		return &cachedThumbnail
	} else {
		l.Error("Error checking cache", "error", err)
	}

	return nil
}

func (s *MessageService) getFileFromCache(ctx context.Context, cacheKey string) *file.File {
	l := app.GetMessageServiceLogger()

	cachedFile, err := c.Get[file.File](s.cache, cacheKey)
	if err == nil {
		l.Debug("Found cached file", "cacheKey", cacheKey)
		return &cachedFile
	} else {
		l.Error("Error checking cache", "error", err)
	}

	return nil
}

func (s *MessageService) uploadFile(ctx context.Context, inst *instance.Instance, source string, stream io.ReadCloser, kind whatsapp.MediaKind, mime string) (*file.File, error) {
	l := app.GetMessageServiceLogger()

	source256 := sha256.Sum256([]byte(source))
	cacheKey := cache.CacheKeyFileUploadPrefix + hex.EncodeToString(source256[:])

	l.Debug("No cached file found, uploading", "cacheKey", cacheKey)
	l.Debug("Uploading file to WhatsApp servers")
	uploadedFile, err := s.whatsapp.UploadFile(ctx, inst, stream, kind, mime)
	if err != nil {
		l.Error("Error uploading file", "error", err)
		return nil, err
	}

	return uploadedFile, nil
}

func (s *MessageService) getThumbnail(ctx context.Context, source string, useCache bool) (*string, error) {
	l := app.GetMessageServiceLogger()

	source256 := sha256.Sum256([]byte(source))
	cacheKey := cache.CacheKeyThumbnailPrefix + hex.EncodeToString(source256[:])

	if useCache {
		thumbnail := s.getThumbnailFromCache(ctx, cacheKey)
		if thumbnail != nil {
			return thumbnail, nil
		}
	}

	l.Debug("Getting thumbnail")
	_, fileThumbnailData, err := s.fileService.GetFrom(ctx, source)
	if err != nil {
		l.Error("Error getting thumbnail", "error", err)
		return nil, err
	}

	thumbnailBase64 := base64.StdEncoding.EncodeToString(*fileThumbnailData)
	if useCache {
		c.Set(s.cache, cacheKey, thumbnailBase64, s.cacheFileUploadTTL)
	}

	return &thumbnailBase64, nil
}
