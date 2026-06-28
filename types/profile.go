package types

type Profile struct {
	Contact        *User `json:"contact"`
	ProfileOptions []int `json:"profileOptions,omitempty"`
}
