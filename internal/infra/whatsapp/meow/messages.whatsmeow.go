package meow

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/group"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/message"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (g *WhatsmeowGateway) GenerateMessageID(ctx context.Context, inst *instance.Instance) (string, error) {
	val, ok := g.clients.Load(inst.ID)
	if !ok {
		return "", instance.ErrInstanceNotConnected
	}

	client := val.(*whatsmeow.Client)
	return client.GenerateMessageID(), nil
}

func (g *WhatsmeowGateway) SendTextMessage(ctx context.Context, inst *instance.Instance, msg *message.Message) (*message.Message, error) {
	l := app.GetMessageServiceLogger()

	if msg.Content.Kind() != message.MessageKindText {
		return nil, fmt.Errorf("invalid message content kind, expected %s but got %s", message.MessageKindText, msg.Content.Kind())
	}

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	content := msg.Content.(message.TextContent)

	context := &waE2E.ContextInfo{
		Expiration: msg.Expiration,
	}

	if content.HasMentions() {
		context.MentionedJID = *content.Mentions
	}

	extended := &waE2E.ExtendedTextMessage{
		Text:        &content.Text,
		ContextInfo: context,
	}

	whatsappMessage := &waE2E.Message{
		ExtendedTextMessage: extended,
	}

	to, err := types.ParseJID(msg.Chat)
	if err != nil {
		return nil, err
	}

	extra := whatsmeow.SendRequestExtra{}
	if msg.ExternalID != nil {
		extra.ID = *msg.ExternalID
	}

	resp, err := client.SendMessage(ctx, to, whatsappMessage, extra)
	if err != nil {
		if err.Error() == ErrServerStatus420.Error() {
			if to.Server == types.GroupServer {
				l.Error("Failed to send message to group, maybe the instance is not in the group", "phone", inst.Phone, "chat", msg.Chat)
				return nil, group.ErrMaybeNotMember
			}
		}
		return nil, err
	}

	msg.ExternalID = &resp.ID

	return msg, nil
}

func (g *WhatsmeowGateway) SendImageMessage(ctx context.Context, inst *instance.Instance, msg *message.Message) (*message.Message, error) {
	if msg.Content.Kind() != message.MessageKindImage {
		return nil, fmt.Errorf("invalid message content kind, expected %s but got %s", message.MessageKindImage, msg.Content.Kind())
	}

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	content := msg.Content.(message.ImageContent)

	context := &waE2E.ContextInfo{
		Expiration: msg.Expiration,
	}

	if content.HasMentions() {
		context.MentionedJID = *content.Mentions
	}

	mediaKey, _ := hex.DecodeString(content.Image.MediaKey)
	sa256Enc, _ := hex.DecodeString(content.Image.Sha256Enc)
	sa256, _ := hex.DecodeString(content.Image.Sha256)

	whatsappMessage := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			URL:           &content.Image.URL,
			DirectPath:    &content.Image.DirectPath,
			Mimetype:      &content.Image.Mime,
			MediaKey:      mediaKey,
			FileEncSHA256: sa256Enc,
			FileSHA256:    sa256,
			FileLength:    &content.Image.Size,
			Caption:       content.Caption,
			Height:        content.Image.Height,
			Width:         content.Image.Width,
			ContextInfo:   context,
			ViewOnce:      content.ViewOnce,
		},
	}

	if content.HasThumbnail() {
		data, err := base64.StdEncoding.DecodeString(*content.Thumbnail)
		if err == nil {
			whatsappMessage.ImageMessage.JPEGThumbnail = data
		}
	}

	to, err := types.ParseJID(msg.Chat)
	if err != nil {
		return nil, err
	}

	extra := whatsmeow.SendRequestExtra{}
	if msg.ExternalID != nil {
		extra.ID = *msg.ExternalID
	}

	resp, err := client.SendMessage(ctx, to, whatsappMessage, extra)
	if err != nil {
		return nil, err
	}

	msg.ExternalID = &resp.ID

	return msg, nil
}

func (g *WhatsmeowGateway) SendVideoMessage(ctx context.Context, inst *instance.Instance, msg *message.Message) (*message.Message, error) {
	if msg.Content.Kind() != message.MessageKindVideo {
		return nil, fmt.Errorf("invalid message content kind, expected %s but got %s", message.MessageKindVideo, msg.Content.Kind())
	}

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	content := msg.Content.(*message.VideoContent)

	context := &waE2E.ContextInfo{
		Expiration: msg.Expiration,
	}

	if content.HasMentions() {
		context.MentionedJID = *content.Mentions
	}

	whatsappMessage := &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			URL:           &content.Video.URL,
			DirectPath:    &content.Video.DirectPath,
			MediaKey:      []byte(content.Video.MediaKey),
			Mimetype:      &content.Video.Mime,
			FileEncSHA256: []byte(content.Video.Sha256Enc),
			FileSHA256:    []byte(content.Video.Sha256),
			FileLength:    &content.Video.Size,
			Caption:       content.Caption,
			Height:        content.Video.Height,
			Width:         content.Video.Width,
			Seconds:       content.Video.Duration,
			ContextInfo:   context,
			ViewOnce:      content.ViewOnce,
		},
	}

	if content.HasThumbnail() {
		data, err := base64.StdEncoding.DecodeString(*content.Thumbnail)
		if err == nil {
			whatsappMessage.ImageMessage.JPEGThumbnail = data
		}
	}

	to, err := types.ParseJID(msg.Chat)
	if err != nil {
		return nil, err
	}

	extra := whatsmeow.SendRequestExtra{}
	if msg.ExternalID != nil {
		extra.ID = *msg.ExternalID
	}

	resp, err := client.SendMessage(ctx, to, whatsappMessage, extra)
	if err != nil {
		return nil, err
	}

	msg.ExternalID = &resp.ID

	return msg, nil
}

func (g *WhatsmeowGateway) SendAudioMessage(ctx context.Context, inst *instance.Instance, msg *message.Message) (*message.Message, error) {
	if msg.Content.Kind() != message.MessageKindAudio {
		return nil, fmt.Errorf("invalid message content kind, expected %s but got %s", message.MessageKindAudio, msg.Content.Kind())
	}

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	content := msg.Content.(*message.AudioContent)

	context := &waE2E.ContextInfo{
		Expiration: msg.Expiration,
	}

	whatsappMessage := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           &content.Audio.URL,
			DirectPath:    &content.Audio.DirectPath,
			MediaKey:      []byte(content.Audio.MediaKey),
			Mimetype:      &content.Audio.Mime,
			FileEncSHA256: []byte(content.Audio.Sha256Enc),
			FileSHA256:    []byte(content.Audio.Sha256),
			FileLength:    &content.Audio.Size,
			PTT:           proto.Bool(false),
			Seconds:       content.Audio.Duration,
			ContextInfo:   context,
		},
	}

	to, err := types.ParseJID(msg.Chat)
	if err != nil {
		return nil, err
	}

	extra := whatsmeow.SendRequestExtra{}
	if msg.ExternalID != nil {
		extra.ID = *msg.ExternalID
	}

	resp, err := client.SendMessage(ctx, to, whatsappMessage, extra)
	if err != nil {
		return nil, err
	}

	msg.ExternalID = &resp.ID

	return msg, nil
}

func (g *WhatsmeowGateway) SendVoiceMessage(ctx context.Context, inst *instance.Instance, msg *message.Message) (*message.Message, error) {
	if msg.Content.Kind() != message.MessageKindVoice {
		return nil, fmt.Errorf("invalid message content kind, expected %s but got %s", message.MessageKindVoice, msg.Content.Kind())
	}

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	content := msg.Content.(*message.VoiceContent)

	context := &waE2E.ContextInfo{
		Expiration: msg.Expiration,
	}

	whatsappMessage := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           &content.Voice.URL,
			DirectPath:    &content.Voice.DirectPath,
			MediaKey:      []byte(content.Voice.MediaKey),
			Mimetype:      &content.Voice.Mime,
			FileEncSHA256: []byte(content.Voice.Sha256Enc),
			FileSHA256:    []byte(content.Voice.Sha256),
			FileLength:    &content.Voice.Size,
			PTT:           proto.Bool(true),
			Seconds:       content.Voice.Duration,
			ContextInfo:   context,
			ViewOnce:      content.ViewOnce,
		},
	}

	to, err := types.ParseJID(msg.Chat)
	if err != nil {
		return nil, err
	}

	extra := whatsmeow.SendRequestExtra{}
	if msg.ExternalID != nil {
		extra.ID = *msg.ExternalID
	}

	resp, err := client.SendMessage(ctx, to, whatsappMessage, extra)
	if err != nil {
		return nil, err
	}

	msg.ExternalID = &resp.ID

	return msg, nil
}

func (g *WhatsmeowGateway) SendDocumentMessage(ctx context.Context, inst *instance.Instance, msg *message.Message) (*message.Message, error) {
	if msg.Content.Kind() != message.MessageKindDocument {
		return nil, fmt.Errorf("invalid message content kind, expected %s but got %s", message.MessageKindDocument, msg.Content.Kind())
	}

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	content := msg.Content.(*message.DocumentContent)

	context := &waE2E.ContextInfo{
		Expiration: msg.Expiration,
	}

	whatsappMessage := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			URL:           &content.Document.URL,
			DirectPath:    &content.Document.DirectPath,
			MediaKey:      []byte(content.Document.MediaKey),
			Mimetype:      &content.Document.Mime,
			FileEncSHA256: []byte(content.Document.Sha256Enc),
			FileSHA256:    []byte(content.Document.Sha256),
			FileLength:    &content.Document.Size,
			Title:         proto.String(content.Document.Name),
			Caption:       content.Caption,
			PageCount:     content.Document.Pages,
			ContextInfo:   context,
		},
	}

	if content.HasThumbnail() {
		data, err := base64.StdEncoding.DecodeString(*content.Thumbnail)
		if err == nil {
			whatsappMessage.ImageMessage.JPEGThumbnail = data
		}
	}

	to, err := types.ParseJID(msg.Chat)
	if err != nil {
		return nil, err
	}

	extra := whatsmeow.SendRequestExtra{}
	if msg.ExternalID != nil {
		extra.ID = *msg.ExternalID
	}

	resp, err := client.SendMessage(ctx, to, whatsappMessage, extra)
	if err != nil {
		return nil, err
	}

	msg.ExternalID = &resp.ID

	return msg, nil
}

func (g *WhatsmeowGateway) ReadMessages(ctx context.Context, inst *instance.Instance, chat string, ids []string, sender string) error {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return err
	}

	chatJID, err := types.ParseJID(chat)
	if err != nil {
		return err
	}

	senderJID, err := types.ParseJID(sender)
	if err != nil {
		return err
	}

	return client.MarkRead(ids, time.Now().UTC(), chatJID, senderJID)
}

func (g *WhatsmeowGateway) SendReaction(ctx context.Context, inst *instance.Instance, msg *message.Message) (*message.Message, error) {
	if msg.Content.Kind() != message.MessageKindReaction {
		return nil, fmt.Errorf("invalid message content kind, expected %s but got %s", message.MessageKindReaction, msg.Content.Kind())
	}

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	content := msg.Content.(message.ReactionContent)

	to, err := types.ParseJID(msg.Chat)
	if err != nil {
		return nil, err
	}

	isMe := true

	whatsappMessage := &waE2E.Message{
		ReactionMessage: &waE2E.ReactionMessage{
			Text: &content.Emoji,
			Key: &waCommon.MessageKey{
				ID:        &content.Message,
				RemoteJID: &msg.Chat,
				FromMe:    &isMe,
			},
		},
	}

	extra := whatsmeow.SendRequestExtra{}
	if msg.ExternalID != nil {
		extra.ID = *msg.ExternalID
	}

	resp, err := client.SendMessage(ctx, to, whatsappMessage, extra)
	if err != nil {
		return nil, err
	}

	msg.ExternalID = &resp.ID

	return msg, nil
}
