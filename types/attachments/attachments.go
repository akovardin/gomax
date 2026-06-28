package attachments

import "encoding/json"

type Attachment struct {
	Type string          `json:"type"`
	Raw  json.RawMessage `json:"-"`
}

func (a *Attachment) UnmarshalJSON(data []byte) error {
	a.Raw = make(json.RawMessage, len(data))
	copy(a.Raw, data)

	var typeOnly struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &typeOnly); err != nil {
		return err
	}
	a.Type = typeOnly.Type
	return nil
}
