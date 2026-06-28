package attachments

type StickerAttachment struct {
	Type        string `json:"type"`
	URL         string `json:"url"`
	StickerID   int    `json:"stickerId"`
	Tags        string `json:"tags"`
	Width       int    `json:"width"`
	SetID       int    `json:"setId"`
	Time        int    `json:"time"`
	StickerType string `json:"stickerType"`
	Audio       string `json:"audio"`
	Height      int    `json:"height"`
	AuthorType  string `json:"authorType"`
	LottieURL   string `json:"lottieUrl,omitempty"`
}
