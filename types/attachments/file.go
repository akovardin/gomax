package attachments

type FileAttachment struct {
	Type   string `json:"type"`
	FileID int    `json:"fileId"`
	Name   string `json:"name"`
	Size   int    `json:"size"`
	Token  string `json:"token"`
}

type FileRequest struct {
	Unsafe string `json:"unsafe"`
	URL    string `json:"url"`
}
