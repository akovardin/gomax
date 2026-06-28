package types

type LoginConfig struct {
	Hash *string `json:"hash,omitempty"`
}

type LoginResponse struct {
	Chats    []*Chat           `json:"chats"`
	Profile  *Profile          `json:"profile"`
	Messages map[int][]*Message `json:"messages"`
	Contacts []*User           `json:"contacts"`
	Token    *string           `json:"token,omitempty"`
	Time     *int              `json:"time,omitempty"`
	Config   *LoginConfig      `json:"config,omitempty"`
}
