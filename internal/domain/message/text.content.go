package message

type TextContent struct {
	Text     string    `json:"text"`
	Mentions *[]string `json:"mentions"`
}

func NewTextContent(text string, mentions *[]string) TextContent {
	return TextContent{
		Text:     text,
		Mentions: mentions,
	}
}

func (t TextContent) Kind() MessageKind {
	return MessageKindText
}

func (t *TextContent) TextForMentions() string {
	return t.Text
}

func (t *TextContent) HasMentions() bool {
	return t.Mentions != nil && len(*t.Mentions) > 0
}
