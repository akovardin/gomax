package attachments

type AudioAttachment struct {
	Type                string `json:"type"`
	Duration            int    `json:"duration"`
	AudioID             int    `json:"audioId"`
	Wave                string `json:"wave"`
	TranscriptionStatus string `json:"transcriptionStatus,omitempty"`
	URL                 string `json:"url"`
	Token               string `json:"token"`
}
