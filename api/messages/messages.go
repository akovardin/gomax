package messages

import (
	"fmt"
	"math/rand"

	"github.com/akovardin/gomax/api/core"
	"github.com/akovardin/gomax/protocol"
	"github.com/akovardin/gomax/types"
	"github.com/akovardin/gomax/types/attachments"
)

type Service struct {
	app core.AppInterface
}

func NewService(app core.AppInterface) *Service {
	return &Service{app: app}
}

func (s *Service) SendMessage(chatID int, text string, replyTo *int, atts []interface{}, notify bool) (*types.Message, error) {
	cid := rand.Intn(1<<31 - 1)
	msg := map[string]interface{}{
		"text":     text,
		"cid":      cid,
		"elements": []interface{}{},
		"attaches": []interface{}{},
	}
	if replyTo != nil {
		msg["link"] = map[string]interface{}{
			"type":      "REPLY",
			"messageId": *replyTo,
		}
	}
	if atts != nil {
		msg["attaches"] = atts
	}

	payload := map[string]interface{}{
		"chatId":  chatID,
		"message": msg,
		"notify":  notify,
	}

	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeMsgSend), payload)
	if err != nil {
		return nil, err
	}
	return parseMessageResponse(frame.Payload)
}

func (s *Service) ForwardMessage(chatID int, messageID int, sourceChatID int, notify bool) (*types.Message, error) {
	cid := rand.Intn(1<<31 - 1)
	msg := map[string]interface{}{
		"cid":      cid,
		"attaches": []interface{}{},
		"link": map[string]interface{}{
			"type":      "FORWARD",
			"messageId": fmt.Sprintf("%d", messageID),
			"chatId":    sourceChatID,
		},
	}

	payload := map[string]interface{}{
		"chatId":  chatID,
		"message": msg,
		"notify":  notify,
	}

	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeMsgSend), payload)
	if err != nil {
		return nil, err
	}
	return parseMessageResponse(frame.Payload)
}

func (s *Service) GetMessages(chatID int, messageIDs []int) ([]*types.Message, error) {
	payload := map[string]interface{}{
		"chatId":     chatID,
		"messageIds": messageIDs,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeMsgGet), payload)
	if err != nil {
		return nil, err
	}
	if msgsRaw, ok := frame.Payload["messages"]; ok {
		msgs, err := core.ConvertStruct[[]*types.Message](msgsRaw)
		if err != nil {
			return nil, err
		}
		return *msgs, nil
	}
	if msgsRaw, ok := frame.Payload["message"]; ok {
		msgs, err := core.ConvertStruct[[]*types.Message](msgsRaw)
		if err != nil {
			return nil, err
		}
		return *msgs, nil
	}
	msgs, err := core.ConvertStruct[[]*types.Message](frame.Payload)
	if err != nil {
		return nil, err
	}
	return *msgs, nil
}

func (s *Service) GetMessage(chatID int, messageID int) (*types.Message, error) {
	payload := map[string]interface{}{
		"chatId":     chatID,
		"messageIds": []int{messageID},
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeMsgGet), payload)
	if err != nil {
		return nil, err
	}
	if msgsRaw, ok := frame.Payload["messages"]; ok {
		msgs, err := core.ConvertStruct[[]*types.Message](msgsRaw)
		if err != nil {
			return nil, err
		}
		if len(*msgs) > 0 {
			return (*msgs)[0], nil
		}
		return nil, nil
	}
	return core.RequirePayloadModel[types.Message](frame.Payload)
}

func (s *Service) EditMessage(chatID int, messageID int, text string, atts []interface{}) (*types.Message, error) {
	payload := map[string]interface{}{
		"chatId":    chatID,
		"messageId": messageID,
		"text":      text,
		"elements":  []interface{}{},
	}
	if atts != nil {
		payload["attachments"] = atts
	} else {
		payload["attachments"] = []interface{}{}
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeMsgEdit), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.Message](frame.Payload)
}

func (s *Service) FetchHistory(chatID int, forward int, backward int) ([]*types.Message, error) {
	payload := map[string]interface{}{
		"chatId":       chatID,
		"forward":      forward,
		"backward":     backward,
		"backwardTime": 0,
		"forwardTime":  0,
		"getChat":      false,
		"itemType":     "REGULAR",
		"getMessages":  true,
		"interactive":  false,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatHistory), payload)
	if err != nil {
		return nil, err
	}
	if msgsRaw, ok := frame.Payload["messages"]; ok {
		msgs, err := core.ConvertStruct[[]*types.Message](msgsRaw)
		if err != nil {
			return nil, err
		}
		return *msgs, nil
	}
	msgs, err := core.ConvertStruct[[]*types.Message](frame.Payload)
	if err != nil {
		return nil, err
	}
	return *msgs, nil
}

func (s *Service) DeleteMessage(chatID int, messageIDs []int, forMe bool) (bool, error) {
	payload := map[string]interface{}{
		"chatId":     chatID,
		"messageIds": messageIDs,
		"forMe":      forMe,
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeMsgDelete), payload)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) PinMessage(chatID int, messageID int, notifyPin bool) (bool, error) {
	payload := map[string]interface{}{
		"chatId": chatID,
		"pinnedMessage": map[string]interface{}{
			"messageId": messageID,
			"notify":    notifyPin,
		},
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatUpdate), payload)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) GetVideoByID(chatID int, messageID int, videoID int) (*attachments.VideoRequest, error) {
	payload := map[string]interface{}{
		"chatId":    chatID,
		"messageId": messageID,
		"videoId":   videoID,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeVideoPlay), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[attachments.VideoRequest](frame.Payload)
}

func (s *Service) GetFileByID(chatID int, messageID int, fileID int) (*attachments.FileRequest, error) {
	payload := map[string]interface{}{
		"chatId":    chatID,
		"messageId": messageID,
		"fileId":    fileID,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeFileDownload), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[attachments.FileRequest](frame.Payload)
}

func (s *Service) AddReaction(chatID int, messageID int, reaction string) (*types.ReactionInfo, error) {
	payload := map[string]interface{}{
		"chatId":    chatID,
		"messageId": fmt.Sprintf("%d", messageID),
		"reaction": map[string]interface{}{
			"reactionType": "EMOJI",
			"id":           reaction,
		},
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeMsgReaction), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.ReactionInfo](frame.Payload)
}

func (s *Service) GetReactions(chatID int, messageIDs []int) (map[string]*types.ReactionInfo, error) {
	strIDs := make([]string, len(messageIDs))
	for i, id := range messageIDs {
		strIDs[i] = fmt.Sprintf("%d", id)
	}
	payload := map[string]interface{}{
		"chatId":     chatID,
		"messageIds": strIDs,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeMsgGetReactions), payload)
	if err != nil {
		return nil, err
	}
	reactions, err := core.ConvertStruct[map[string]*types.ReactionInfo](frame.Payload)
	if err != nil {
		return nil, err
	}
	return *reactions, nil
}

func (s *Service) RemoveReaction(chatID int, messageID int) (*types.ReactionInfo, error) {
	payload := map[string]interface{}{
		"chatId":    chatID,
		"messageId": fmt.Sprintf("%d", messageID),
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeMsgCancelReaction), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.ReactionInfo](frame.Payload)
}

func (s *Service) ReadMessage(messageID int, chatID int) (*types.ReadState, error) {
	payload := map[string]interface{}{
		"type":      0,
		"chatId":    chatID,
		"messageId": messageID,
		"mark":      messageID,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeChatMark), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.ReadState](frame.Payload)
}

func parseMessageResponse(payload map[string]interface{}) (*types.Message, error) {
	if inner, ok := payload["message"].(map[string]interface{}); ok {
		return core.RequirePayloadModel[types.Message](inner)
	}
	return core.RequirePayloadModel[types.Message](payload)
}
