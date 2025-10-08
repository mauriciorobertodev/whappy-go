package input

import (
	"github.com/mauriciorobertodev/whappy-go/internal/domain/message"
)

type SendTextMessageInput struct {
	ID         *string   `json:"id"`
	To         string    `json:"to"`
	Text       string    `json:"text"`
	Mentions   *[]string `json:"mentions"`
	Expiration *uint32   `json:"expiration"`
}

func (inp *SendTextMessageInput) Validate() error {
	if inp.To == "" {
		return message.ErrInvalidJID
	}

	if inp.Text == "" {
		return message.ErrEmptyText
	}

	if len(inp.Text) > message.MaxMessageTextLength {
		return message.ErrTextTooLong
	}

	return nil
}

type SendImageMessageInput struct {
	ID         *string   `json:"id"`
	To         string    `json:"to"`
	Image      string    `json:"image"`
	Name       *string   `json:"name"`
	Mime       *string   `json:"mime"`
	Width      *uint32   `json:"width"`
	Height     *uint32   `json:"height"`
	Thumbnail  *string   `json:"thumbnail"`
	Caption    *string   `json:"caption"`
	Mentions   *[]string `json:"mentions"`
	Expiration *uint32   `json:"expiration"`
	ViewOnce   *bool     `json:"view_once"`
	Cache      *bool     `json:"cache"`
}

func (inp *SendImageMessageInput) Validate() error {
	if inp.To == "" {
		return message.ErrInvalidJID
	}

	if inp.Image == "" {
		return message.ErrImageRequired
	}

	if inp.Caption != nil && len(*inp.Caption) > message.MaxCaptionLength {
		return message.ErrCaptionTooLong
	}

	return nil
}

type SendVideoMessageInput struct {
	ID         *string   `json:"id"`
	To         string    `json:"to"`
	Video      string    `json:"video"`
	Name       *string   `json:"name"`
	Mime       *string   `json:"mime"`
	Width      *uint32   `json:"width"`
	Height     *uint32   `json:"height"`
	Duration   *uint32   `json:"duration"`
	Thumbnail  *string   `json:"thumbnail"`
	Caption    *string   `json:"caption"`
	Mentions   *[]string `json:"mentions"`
	ViewOnce   *bool     `json:"view_once"`
	Expiration *uint32   `json:"expiration"`
	Cache      *bool     `json:"cache"`
}

func (inp *SendVideoMessageInput) Validate() error {
	if inp.To == "" {
		return message.ErrInvalidJID
	}

	if inp.Video == "" {
		return message.ErrVideoRequired
	}

	if inp.Caption != nil && len(*inp.Caption) > message.MaxCaptionLength {
		return message.ErrCaptionTooLong
	}

	return nil
}

type SendAudioMessageInput struct {
	ID         *string `json:"id"`
	To         string  `json:"to"`
	Audio      string  `json:"audio"`
	Name       *string `json:"name"`
	Mime       *string `json:"mime"`
	Duration   *uint32 `json:"duration"`
	Expiration *uint32 `json:"expiration"`
	Cache      *bool   `json:"cache"`
}

func (inp *SendAudioMessageInput) Validate() error {
	if inp.To == "" {
		return message.ErrInvalidJID
	}

	if inp.Audio == "" {
		return message.ErrAudioRequired
	}

	return nil
}

type SendVoiceMessageInput struct {
	ID         *string `json:"id"`
	To         string  `json:"to"`
	Voice      string  `json:"voice"`
	Name       *string `json:"name"`
	Mime       *string `json:"mime"`
	Duration   *uint32 `json:"duration"`
	ViewOnce   *bool   `json:"view_once"`
	Expiration *uint32 `json:"expiration"`
	Cache      *bool   `json:"cache"`
}

func (inp *SendVoiceMessageInput) Validate() error {
	if inp.To == "" {
		return message.ErrInvalidJID
	}

	if inp.Voice == "" {
		return message.ErrVoiceRequired
	}

	return nil
}

type SendDocumentMessageInput struct {
	ID         *string   `json:"id"`
	To         string    `json:"to"`
	Document   string    `json:"document"`
	Name       *string   `json:"name"`
	Mime       *string   `json:"mime"`
	Pages      *uint32   `json:"pages"`
	Thumbnail  *string   `json:"thumbnail"`
	Caption    *string   `json:"caption"`
	Mentions   *[]string `json:"mentions"`
	Expiration *uint32   `json:"expiration"`
	Cache      *bool     `json:"cache"`
}

func (inp *SendDocumentMessageInput) Validate() error {
	if inp.To == "" {
		return message.ErrInvalidJID
	}

	if inp.Document == "" {
		return message.ErrDocumentRequired
	}

	if inp.Caption != nil && len(*inp.Caption) > message.MaxCaptionLength {
		return message.ErrCaptionTooLong
	}

	return nil
}

type SendReactionMessageInput struct {
	To      string `json:"to"`
	Message string `json:"message"`
	Emoji   string `json:"emoji"`
	Cache   *bool  `json:"cache"`
}

func (inp *SendReactionMessageInput) Validate() error {
	if inp.To == "" {
		return message.ErrInvalidJID
	}

	if inp.Message == "" {
		return message.ErrInvalidMessageID
	}

	return nil
}

type SendReactionInput struct {
	To      string `json:"to"`
	Message string `json:"message"`
	Emoji   string `json:"emoji"`
}

func (inp *SendReactionInput) Validate() error {
	if inp.To == "" {
		return message.ErrInvalidJID
	}

	if inp.Message == "" {
		return message.ErrInvalidMessageID
	}

	return nil
}
