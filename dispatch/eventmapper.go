package dispatch

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/akovardin/gomax/protocol"
	"github.com/akovardin/gomax/types"
)

type EventMapper struct {
	App interface{}
}

func NewEventMapper(app interface{}) *EventMapper {
	return &EventMapper{App: app}
}

func (m *EventMapper) Map(eventType EventType, frame *protocol.InboundFrame) interface{} {
	if frame.Payload == nil {
		return frame
	}

	switch eventType {
	case EventTypeMessageNew, EventTypeMessageEdit:
		msg := m.parseMessage(frame.Payload)
		if msg != nil {
			return msg
		}
	case EventTypeChatUpdate:
		if chatRaw, ok := frame.Payload["chat"]; ok {
			data, _ := json.Marshal(chatRaw)
			var chat types.Chat
			dec := json.NewDecoder(bytes.NewReader(data))
			dec.UseNumber()
			if err := dec.Decode(&chat); err == nil {
				return &chat
			}
		}
		return frame
	}

	switch eventType {
	case EventTypeMessageDelete:
		evt := m.parseMessageDelete(frame.Payload)
		if evt != nil {
			return evt
		}
	case EventTypeMessageRead:
		data, _ := json.Marshal(frame.Payload)
		var evt types.MessageReadEvent
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber()
		if err := dec.Decode(&evt); err == nil {
			return &evt
		}
	case EventTypeTyping:
		data, _ := json.Marshal(frame.Payload)
		var evt types.TypingEvent
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber()
		if err := dec.Decode(&evt); err == nil {
			return &evt
		}
	case EventTypePresence:
		data, _ := json.Marshal(frame.Payload)
		var evt types.PresenceEvent
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber()
		if err := dec.Decode(&evt); err == nil {
			return &evt
		}
	case EventTypeReactionUpdate:
		data, _ := json.Marshal(frame.Payload)
		var evt types.ReactionUpdateEvent
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber()
		if err := dec.Decode(&evt); err == nil {
			return &evt
		}
	case EventTypeVideoReady:
		data, _ := json.Marshal(frame.Payload)
		var evt types.VideoUploadSignal
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber()
		if err := dec.Decode(&evt); err == nil {
			return &evt
		}
	case EventTypeFileReady:
		data, _ := json.Marshal(frame.Payload)
		var evt types.FileUploadSignal
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber()
		if err := dec.Decode(&evt); err == nil {
			return &evt
		}
	}

	return frame
}

func (m *EventMapper) parseMessageDelete(payload map[string]interface{}) *types.MessageDeleteEvent {
	if chatRaw, ok := payload["chat"].(map[string]interface{}); ok {
		chatID := toFlexInt(chatRaw["id"])
		messageIDs := toFlexIntSlice(payload, "messageIds")
		if chatID == 0 || messageIDs == nil {
			return nil
		}
		data, _ := json.Marshal(chatRaw)
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber()
		var chat types.Chat
		dec.Decode(&chat)
		return &types.MessageDeleteEvent{
			Chat:       &chat,
			ChatID:     chatID,
			MessageIDs: messageIDs,
			TTL:        toBool(payload["ttl"]),
		}
	}
	if msgRaw, ok := payload["message"].(map[string]interface{}); ok {
		msgID := toFlexInt(msgRaw["id"])
		chatID := toFlexInt(payload["chatId"])
		if chatID == 0 || msgID == 0 {
			return nil
		}
		return &types.MessageDeleteEvent{
			ChatID:     chatID,
			MessageIDs: []types.FlexInt{msgID},
			TTL:        toBool(payload["ttl"]),
		}
	}
	return nil
}

func toFlexInt(v interface{}) types.FlexInt {
	switch n := v.(type) {
	case float64:
		return types.FlexInt(int64(n))
	case string:
		i, _ := strconv.Atoi(n)
		return types.FlexInt(i)
	case int:
		return types.FlexInt(n)
	case int64:
		return types.FlexInt(n)
	case json.Number:
		i, _ := n.Int64()
		return types.FlexInt(i)
	}
	return 0
}

func toFlexIntSlice(payload map[string]interface{}, key string) []types.FlexInt {
	raw, ok := payload[key]
	if !ok {
		return nil
	}
	arr, ok := raw.([]interface{})
	if !ok {
		return []types.FlexInt{toFlexInt(raw)}
	}
	result := make([]types.FlexInt, 0, len(arr))
	for _, v := range arr {
		result = append(result, toFlexInt(v))
	}
	return result
}

func toBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func (m *EventMapper) parseMessage(payload map[string]interface{}) *types.Message {
	raw := payload

	if inner, ok := payload["message"].(map[string]interface{}); ok {
		raw = make(map[string]interface{})
		for k, v := range inner {
			raw[k] = v
		}
		for _, key := range []string{"chatId", "prevMessageId", "ttl", "unread", "mark"} {
			if v, ok := payload[key]; ok {
				raw[key] = v
			}
		}
	}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var msg types.Message
	if err := dec.Decode(&msg); err != nil {
		return nil
	}
	return &msg
}
