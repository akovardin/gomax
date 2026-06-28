package core

import (
	"github.com/akovardin/gomax/protocol"
	"github.com/akovardin/gomax/types"
)

type AppInterface interface {
	Invoke(opcode int, payload map[string]interface{}, cmd int, timeout float64, compress bool) (*protocol.InboundFrame, error)
	Config() *types.ClientConfig
	Token() string
	UpdateToken(oldToken, newToken string)
	Me() *types.Profile
	Chats() []*types.Chat
	Users() map[types.FlexInt]*types.User
	Contacts() []*types.User
	Messages() map[types.FlexInt][]*types.Message
	CacheChat(chat *types.Chat)
	GetCachedChat(chatID types.FlexInt) *types.Chat
	CacheUser(user *types.User)
	GetCachedUser(userID types.FlexInt) *types.User
	BindMessage(msg *types.Message, chatID types.FlexInt)
}
