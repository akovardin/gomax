package types

import "github.com/akovardin/gomax/types/attachments"

type ReactionCounter struct {
	Count    FlexInt `json:"count"`
	Reaction string  `json:"reaction"`
}

type ReactionInfo struct {
	TotalCount   FlexInt           `json:"totalCount"`
	Counters     []ReactionCounter `json:"counters"`
	YourReaction *string           `json:"yourReaction,omitempty"`
}

type ReadState struct {
	Unread FlexInt `json:"unread"`
	Mark   FlexInt `json:"mark"`
}

type Message struct {
	ID            FlexInt                 `json:"id"`
	ChatID        *FlexInt                `json:"chatId,omitempty"`
	Sender        *FlexInt                `json:"sender,omitempty"`
	Text          string                  `json:"text"`
	Time          FlexInt                 `json:"time"`
	Type          string                  `json:"type"`
	CID           *FlexInt                `json:"cid,omitempty"`
	Attaches      []attachments.Attachment `json:"attaches,omitempty"`
	Stats         map[string]interface{}  `json:"stats,omitempty"`
	Status        *string                 `json:"status,omitempty"`
	ReactionInfo  *ReactionInfo           `json:"reactionInfo,omitempty"`
	Options       interface{}             `json:"options,omitempty"`
	PrevMessageID interface{}             `json:"prevMessageId,omitempty"`
	TTL           *bool                   `json:"ttl,omitempty"`
	Unread        *FlexInt                `json:"unread,omitempty"`
	Mark          *FlexInt                `json:"mark,omitempty"`
	Elements      []Element               `json:"elements,omitempty"`
}
