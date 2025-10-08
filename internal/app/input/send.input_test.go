package input_test

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/message"
	"github.com/mauriciorobertodev/whappy-go/internal/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Send Inputs", func() {
	Describe("SendTextMessageInput Input", func() {
		It("should validate successfully", func() {
			inp := &input.SendTextMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Text:       "Hello, World!",
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation due to empty To", func() {
			inp := &input.SendTextMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "",
				Text:       "Hello, World!",
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
			}
			Expect(inp.Validate()).To(Equal(message.ErrInvalidJID))
		})

		It("should fail validation due to empty Text", func() {
			inp := &input.SendTextMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Text:       "",
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
			}
			Expect(inp.Validate()).To(Equal(message.ErrEmptyText))
		})

		It("should fail validation due to Text too long", func() {
			longText := make([]byte, message.MaxMessageTextLength+1)
			for i := range longText {
				longText[i] = 'a'
			}
			inp := &input.SendTextMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Text:       string(longText),
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
			}
			Expect(inp.Validate()).To(Equal(message.ErrTextTooLong))
		})
	})

	Describe("SendImageMessageInput Input", func() {
		It("should validate successfully", func() {
			inp := &input.SendImageMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Image:      "base64encodedstring",
				Name:       utils.StringPtr("image.jpg"),
				Mime:       utils.StringPtr("image/jpeg"),
				Width:      utils.Uint32Ptr(800),
				Height:     utils.Uint32Ptr(600),
				Thumbnail:  utils.StringPtr("base64thumbnail"),
				Caption:    utils.StringPtr("Check this out!"),
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
				ViewOnce:   utils.BoolPtr(true),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation due to empty To", func() {
			inp := &input.SendImageMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "",
				Image:      "base64encodedstring",
				Name:       utils.StringPtr("image.jpg"),
				Mime:       utils.StringPtr("image/jpeg"),
				Width:      utils.Uint32Ptr(800),
				Height:     utils.Uint32Ptr(600),
				Thumbnail:  utils.StringPtr("base64thumbnail"),
				Caption:    utils.StringPtr("Check this out!"),
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
				ViewOnce:   utils.BoolPtr(true),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(Equal(message.ErrInvalidJID))
		})

		It("should fail validation due to empty Image", func() {
			inp := &input.SendImageMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Image:      "",
				Name:       utils.StringPtr("image.jpg"),
				Mime:       utils.StringPtr("image/jpeg"),
				Width:      utils.Uint32Ptr(800),
				Height:     utils.Uint32Ptr(600),
				Thumbnail:  utils.StringPtr("base64thumbnail"),
				Caption:    utils.StringPtr("Check this out!"),
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
				ViewOnce:   utils.BoolPtr(true),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(Equal(message.ErrImageRequired))
		})

		It("should fail validation due to Caption too long", func() {
			longCaption := make([]byte, message.MaxCaptionLength+1)
			for i := range longCaption {
				longCaption[i] = 'a'
			}
			inp := &input.SendImageMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Image:      "base64encodedstring",
				Name:       utils.StringPtr("image.jpg"),
				Mime:       utils.StringPtr("image/jpeg"),
				Width:      utils.Uint32Ptr(800),
				Height:     utils.Uint32Ptr(600),
				Thumbnail:  utils.StringPtr("base64thumbnail"),
				Caption:    utils.StringPtr(string(longCaption)),
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
				ViewOnce:   utils.BoolPtr(true),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(Equal(message.ErrCaptionTooLong))
		})
	})

	Describe("SendVideoMessageInput Input", func() {
		It("should validate successfully", func() {
			inp := &input.SendVideoMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Video:      "base64encodedstring",
				Name:       utils.StringPtr("video.mp4"),
				Mime:       utils.StringPtr("video/mp4"),
				Width:      utils.Uint32Ptr(1280),
				Height:     utils.Uint32Ptr(720),
				Thumbnail:  utils.StringPtr("base64thumbnail"),
				Caption:    utils.StringPtr("Watch this!"),
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
				ViewOnce:   utils.BoolPtr(true),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation due to empty To", func() {
			inp := &input.SendVideoMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "",
				Video:      "base64encodedstring",
				Name:       utils.StringPtr("video.mp4"),
				Mime:       utils.StringPtr("video/mp4"),
				Width:      utils.Uint32Ptr(1280),
				Height:     utils.Uint32Ptr(720),
				Thumbnail:  utils.StringPtr("base64thumbnail"),
				Caption:    utils.StringPtr("Watch this!"),
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
				ViewOnce:   utils.BoolPtr(true),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(Equal(message.ErrInvalidJID))
		})

		It("should fail validation due to empty Video", func() {
			inp := &input.SendVideoMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Video:      "",
				Name:       utils.StringPtr("video.mp4"),
				Mime:       utils.StringPtr("video/mp4"),
				Width:      utils.Uint32Ptr(1280),
				Height:     utils.Uint32Ptr(720),
				Thumbnail:  utils.StringPtr("base64thumbnail"),
				Caption:    utils.StringPtr("Watch this!"),
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
				ViewOnce:   utils.BoolPtr(true),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(Equal(message.ErrVideoRequired))
		})

		It("should fail validation due to Caption too long", func() {
			longCaption := make([]byte, message.MaxCaptionLength+1)
			for i := range longCaption {
				longCaption[i] = 'a'
			}
			inp := &input.SendVideoMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Video:      "base64encodedstring",
				Name:       utils.StringPtr("video.mp4"),
				Mime:       utils.StringPtr("video/mp4"),
				Width:      utils.Uint32Ptr(1280),
				Height:     utils.Uint32Ptr(720),
				Thumbnail:  utils.StringPtr("base64thumbnail"),
				Caption:    utils.StringPtr(string(longCaption)),
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
				ViewOnce:   utils.BoolPtr(true),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(Equal(message.ErrCaptionTooLong))
		})
	})

	Describe("SendAudioMessageInput Input", func() {
		It("should validate successfully", func() {
			inp := &input.SendAudioMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Audio:      "base64encodedstring",
				Name:       utils.StringPtr("audio.mp3"),
				Mime:       utils.StringPtr("audio/mpeg"),
				Duration:   utils.Uint32Ptr(120),
				Expiration: utils.Uint32Ptr(60),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation due to empty To", func() {
			inp := &input.SendAudioMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "",
				Audio:      "base64encodedstring",
				Name:       utils.StringPtr("audio.mp3"),
				Mime:       utils.StringPtr("audio/mpeg"),
				Duration:   utils.Uint32Ptr(120),
				Expiration: utils.Uint32Ptr(60),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(Equal(message.ErrInvalidJID))
		})

		It("should fail validation due to empty Audio", func() {
			inp := &input.SendAudioMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Audio:      "",
				Name:       utils.StringPtr("audio.mp3"),
				Mime:       utils.StringPtr("audio/mpeg"),
				Duration:   utils.Uint32Ptr(120),
				Expiration: utils.Uint32Ptr(60),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(Equal(message.ErrAudioRequired))
		})
	})

	Describe("SendVoiceMessageInput Input", func() {
		It("should validate successfully", func() {
			inp := &input.SendVoiceMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Voice:      "base64encodedstring",
				Name:       utils.StringPtr("voice.ogg"),
				Mime:       utils.StringPtr("audio/ogg"),
				Duration:   utils.Uint32Ptr(30),
				Expiration: utils.Uint32Ptr(60),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation due to empty To", func() {
			inp := &input.SendVoiceMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "",
				Voice:      "base64encodedstring",
				Name:       utils.StringPtr("voice.ogg"),
				Mime:       utils.StringPtr("audio/ogg"),
				Duration:   utils.Uint32Ptr(30),
				Expiration: utils.Uint32Ptr(60),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(Equal(message.ErrInvalidJID))
		})

		It("should fail validation due to empty Voice", func() {
			inp := &input.SendVoiceMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Voice:      "",
				Name:       utils.StringPtr("voice.ogg"),
				Mime:       utils.StringPtr("audio/ogg"),
				Duration:   utils.Uint32Ptr(30),
				Expiration: utils.Uint32Ptr(60),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(Equal(message.ErrVoiceRequired))
		})
	})

	Describe("SendDocumentMessageInput Input", func() {
		It("should validate successfully", func() {
			inp := &input.SendDocumentMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "551412345678",
				Document:   "base64encodedstring",
				Name:       utils.StringPtr("file.pdf"),
				Mime:       utils.StringPtr("application/pdf"),
				Caption:    utils.StringPtr("Please see the attached document."),
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation due to empty To", func() {
			inp := &input.SendDocumentMessageInput{
				ID:         utils.StringPtr("123e4567-e89b-12d3-a456-426614174000"),
				To:         "",
				Document:   "base64encodedstring",
				Name:       utils.StringPtr("file.pdf"),
				Mime:       utils.StringPtr("application/pdf"),
				Caption:    utils.StringPtr("Please see the attached document."),
				Mentions:   &[]string{"5514987654321"},
				Expiration: utils.Uint32Ptr(60),
				Cache:      utils.BoolPtr(true),
			}
			Expect(inp.Validate()).To(Equal(message.ErrInvalidJID))
		})
	})
})
