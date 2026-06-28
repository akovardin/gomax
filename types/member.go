package types

type Presence struct {
	Seen   *FlexInt `json:"seen,omitempty"`
	Status *FlexInt `json:"status,omitempty"`
}

type Member struct {
	Contact  *User     `json:"contact"`
	Presence *Presence `json:"presence,omitempty"`
}
