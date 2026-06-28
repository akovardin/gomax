package attachments

type PhotoAttachment struct {
	Type        string `json:"type"`
	BaseURL     string `json:"baseUrl"`
	Height      int    `json:"height"`
	Width       int    `json:"width"`
	PhotoID     int    `json:"photoId"`
	PhotoToken  string `json:"photoToken"`
	PreviewData string `json:"previewData,omitempty"`
}
