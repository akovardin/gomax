package types

type MessageDeleteEvent struct {
	MessageIDs []FlexInt `json:"messageIds"`
	ChatID     FlexInt   `json:"chatId"`
	Chat       *Chat     `json:"chat,omitempty"`
	Message    *Message  `json:"message,omitempty"`
	TTL        bool      `json:"ttl"`
}

type MessageReadEvent struct {
	SetAsUnread bool   `json:"setAsUnread"`
	ChatID      FlexInt `json:"chatId"`
	UserID      FlexInt `json:"userId"`
	Mark        FlexInt `json:"mark"`
}

type TypingEvent struct {
	ChatID FlexInt `json:"chatId"`
	UserID FlexInt `json:"userId"`
}

type PresenceEvent struct {
	Presence Presence `json:"presence"`
	UserID   FlexInt  `json:"userId"`
}

type ReactionUpdateEvent struct {
	MessageID  string            `json:"messageId"`
	ChatID     FlexInt           `json:"chatId"`
	Counters   []ReactionCounter `json:"counters"`
	TotalCount FlexInt           `json:"totalCount"`
}

type VideoUploadSignal struct {
	VideoID FlexInt `json:"videoId"`
}

type FileUploadSignal struct {
	FileID FlexInt `json:"fileId"`
}
