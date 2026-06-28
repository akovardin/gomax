package keyboards

type InlineKeyboardAttachment struct {
	Type     string                 `json:"type"`
	Keyboard map[string]interface{} `json:"keyboard"`
}
