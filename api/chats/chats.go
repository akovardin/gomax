package chats

import (
	"strings"
	"time"

	"github.com/akovardin/gomax/api/core"
	"github.com/akovardin/gomax/protocol"
	"github.com/akovardin/gomax/types"
)

type Service struct {
	app core.AppInterface
}

func NewService(app core.AppInterface) *Service {
	return &Service{app: app}
}

func extractLinkHash(link string) string {
	parts := strings.Split(link, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			return parts[i]
		}
	}
	return link
}

func (s *Service) CreateGroup(name string, participantIDs []int, notify bool) (*types.Chat, *types.Message, error) {
	payload := map[string]interface{}{
		"name":    name,
		"userIds": participantIDs,
		"notify":  notify,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatCreate), payload)
	if err != nil {
		return nil, nil, err
	}
	chat, err := core.RequirePayloadModel[types.Chat](frame.Payload)
	if err != nil {
		return nil, nil, err
	}
	var msg *types.Message
	if msgRaw, ok := frame.Payload["message"]; ok {
		msg, _ = core.ConvertStruct[types.Message](msgRaw)
	}
	return chat, msg, nil
}

func (s *Service) InviteUsersToGroup(chatID int, userIDs []int, showHistory bool) (*types.Chat, error) {
	payload := map[string]interface{}{
		"chatId":      chatID,
		"userIds":     userIDs,
		"showHistory": showHistory,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatMembersUpdate), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.Chat](frame.Payload)
}

func (s *Service) RemoveUsersFromGroup(chatID int, userIDs []int, cleanMsgPeriod int) (bool, error) {
	payload := map[string]interface{}{
		"chatId":         chatID,
		"userIds":        userIDs,
		"cleanMsgPeriod": cleanMsgPeriod,
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatMembersUpdate), payload)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) ChangeGroupSettings(chatID int, title *string, description *string) error {
	payload := map[string]interface{}{
		"chatId": chatID,
	}
	if title != nil {
		payload["title"] = *title
	}
	if description != nil {
		payload["description"] = *description
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatUpdate), payload)
	return err
}

func (s *Service) JoinGroup(link string) (*types.Chat, error) {
	hash := extractLinkHash(link)
	payload := map[string]interface{}{
		"hash": hash,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatJoin), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.Chat](frame.Payload)
}

func (s *Service) JoinChannel(link string) (*types.Chat, error) {
	hash := extractLinkHash(link)
	payload := map[string]interface{}{
		"hash": hash,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatJoin), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.Chat](frame.Payload)
}

func (s *Service) ResolveGroupByLink(link string) (*types.Chat, error) {
	hash := extractLinkHash(link)
	payload := map[string]interface{}{
		"hash": hash,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatCheckLink), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.Chat](frame.Payload)
}

func (s *Service) ReworkInviteLink(chatID int) (*types.Chat, error) {
	payload := map[string]interface{}{
		"chatId": chatID,
		"action": "rework_invite_link",
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatUpdate), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.Chat](frame.Payload)
}

func (s *Service) GetChats(chatIDs []int) ([]*types.Chat, error) {
	payload := map[string]interface{}{
		"chatIds": chatIDs,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatInfo), payload)
	if err != nil {
		return nil, err
	}
	if chatsRaw, ok := frame.Payload["chats"]; ok {
		chats, err := core.ConvertStruct[[]*types.Chat](chatsRaw)
		if err != nil {
			return nil, err
		}
		return *chats, nil
	}
	if chatRaw, ok := frame.Payload["chat"]; ok {
		c, err := core.ConvertStruct[types.Chat](chatRaw)
		if err != nil {
			return nil, err
		}
		return []*types.Chat{c}, nil
	}
	chat, err := core.RequirePayloadModel[types.Chat](frame.Payload)
	if err != nil {
		return nil, err
	}
	return []*types.Chat{chat}, nil
}

func (s *Service) GetChat(chatID int) (*types.Chat, error) {
	payload := map[string]interface{}{
		"chatIds": []int{chatID},
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatInfo), payload)
	if err != nil {
		return nil, err
	}
	if chatRaw, ok := frame.Payload["chat"]; ok {
		return core.ConvertStruct[types.Chat](chatRaw)
	}
	if chatsRaw, ok := frame.Payload["chats"]; ok {
		chats, err := core.ConvertStruct[[]*types.Chat](chatsRaw)
		if err != nil {
			return nil, err
		}
		if len(*chats) > 0 {
			return (*chats)[0], nil
		}
		return nil, nil
	}
	return core.RequirePayloadModel[types.Chat](frame.Payload)
}

func (s *Service) LeaveGroup(chatID int) error {
	payload := map[string]interface{}{
		"chatId": chatID,
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatLeave), payload)
	return err
}

func (s *Service) LeaveChannel(chatID int) error {
	payload := map[string]interface{}{
		"chatId": chatID,
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatLeave), payload)
	return err
}

func (s *Service) DeleteChat(chatID int) error {
	payload := map[string]interface{}{
		"chatId": chatID,
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatDelete), payload)
	return err
}

func (s *Service) FetchChats(marker interface{}) ([]*types.Chat, error) {
	if marker == nil {
		marker = int(time.Now().UnixMilli())
	}
	payload := map[string]interface{}{
		"marker": marker,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatsList), payload)
	if err != nil {
		return nil, err
	}
	if chatsRaw, ok := frame.Payload["chats"]; ok {
		chats, err := core.ConvertStruct[[]*types.Chat](chatsRaw)
		if err != nil {
			return nil, err
		}
		return *chats, nil
	}
	return []*types.Chat{}, nil
}

func (s *Service) GetJoinRequests(chatID int, count int) ([]*types.Member, error) {
	payload := map[string]interface{}{
		"chatId": chatID,
		"count":  count,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatMembers), payload)
	if err != nil {
		return nil, err
	}
	if membersRaw, ok := frame.Payload["members"]; ok {
		members, err := core.ConvertStruct[[]*types.Member](membersRaw)
		if err != nil {
			return nil, err
		}
		return *members, nil
	}
	members, err := core.ConvertStruct[[]*types.Member](frame.Payload)
	if err != nil {
		return nil, err
	}
	return *members, nil
}

func (s *Service) ConfirmJoinRequests(chatID int, userIDs []int) (bool, error) {
	payload := map[string]interface{}{
		"chatId":  chatID,
		"userIds": userIDs,
		"action":  "confirm",
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatMembersUpdate), payload)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) DeclineJoinRequests(chatID int, userIDs []int) (bool, error) {
	payload := map[string]interface{}{
		"chatId":  chatID,
		"userIds": userIDs,
		"action":  "decline",
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatMembersUpdate), payload)
	if err != nil {
		return false, err
	}
	return true, nil
}
