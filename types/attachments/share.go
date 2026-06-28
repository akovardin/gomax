package attachments

type ShareAttachment struct {
	Type        string `json:"type"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
