package attachments

import "encoding/json"

type UnknownAttachment struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
