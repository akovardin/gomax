package types

type ElementAttributes struct {
	URL *string `json:"url,omitempty"`
}

type Element struct {
	Type       string             `json:"type"`
	From       *int               `json:"from,omitempty"`
	Length     *int               `json:"length,omitempty"`
	Attributes *ElementAttributes `json:"attributes,omitempty"`
}
