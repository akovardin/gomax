package attachments

type VideoAttachment struct {
	Type        string `json:"type"`
	Height      int    `json:"height"`
	Width       int    `json:"width"`
	VideoID     int    `json:"videoId"`
	Duration    int    `json:"duration"`
	PreviewData string `json:"previewData,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	Token       string `json:"token"`
	VideoType   string `json:"videoType,omitempty"`
}

type VideoRequest struct {
	External string `json:"external"`
	Cache    string `json:"cache"`
	URL      string `json:"url"`
}
