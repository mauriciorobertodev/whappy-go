package message

type ReactionContent struct {
	Emoji   string `json:"emoji"`
	Message string `json:"message"`
	Removed bool   `json:"removed"`
}

func NewReactionContent(emoji string, message string) ReactionContent {
	return ReactionContent{
		Emoji:   emoji,
		Message: message,
		Removed: emoji == "",
	}
}

func (t ReactionContent) Kind() MessageKind {
	return MessageKindReaction
}
