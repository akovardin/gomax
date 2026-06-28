package attachments

type CallAttachment struct {
	Type           string `json:"type"`
	Duration       int    `json:"duration"`
	ConversationID int    `json:"conversationId"`
	ContactIDs     []int  `json:"contactIds"`
}
