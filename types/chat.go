package types

type Chat struct {
	ID                     FlexInt                 `json:"id"`
	Type                   string                 `json:"type"`
	Status                 string                 `json:"status"`
	Owner                  FlexInt                `json:"owner"`
	Participants           map[string]FlexInt     `json:"participants"`
	Title                  *string                `json:"title,omitempty"`
	BaseRawIconURL         *string                `json:"baseRawIconUrl,omitempty"`
	BaseIconURL            *string                `json:"baseIconUrl,omitempty"`
	LastMessage            *Message               `json:"lastMessage,omitempty"`
	LastEventTime          FlexInt                `json:"lastEventTime"`
	LastDelayedUpdateTime  FlexInt                `json:"lastDelayedUpdateTime"`
	Created                FlexInt                `json:"created"`
	NewMessages            FlexInt                `json:"newMessages"`
	Link                   *string                `json:"link,omitempty"`
	Access                 *string                `json:"access,omitempty"`
	Restrictions           *FlexInt               `json:"restrictions,omitempty"`
	PinnedMessage          *Message               `json:"pinnedMessage,omitempty"`
	ParticipantsCount      FlexInt                `json:"participantsCount"`
	Description            *string                `json:"description,omitempty"`
	Options                interface{}            `json:"options,omitempty"`
	JoinTime               FlexInt                `json:"joinTime"`
	InvitedBy              *FlexInt               `json:"invitedBy,omitempty"`
	Modified               FlexInt                `json:"modified"`
	MessagesCount          FlexInt                `json:"messagesCount"`
	HasBots                *bool                  `json:"hasBots,omitempty"`
	PrevMessageID          *FlexInt               `json:"prevMessageId,omitempty"`
	AdminParticipants      map[string]interface{} `json:"adminParticipants"`
	Admins                 []FlexInt              `json:"admins"`
	CID                    *FlexInt               `json:"cid,omitempty"`
}

func (c *Chat) IsDialog() bool  { return c.Type == "DIALOG" }
func (c *Chat) IsGroup() bool   { return c.Type == "GROUP" }
func (c *Chat) IsChannel() bool { return c.Type == "CHANNEL" }
