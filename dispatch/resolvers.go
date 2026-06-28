package dispatch

import (
	"github.com/akovardin/gomax/protocol"
)

func ResolveMessage(frame *protocol.InboundFrame) EventType {
	if frame.Opcode != int(protocol.OpcodeNotifMessage) && frame.Opcode != int(protocol.OpcodeMsgEdit) {
		return ""
	}

	if frame.Payload == nil {
		return EventTypeMessageNew
	}

	raw := frame.Payload
	if inner, ok := raw["message"].(map[string]interface{}); ok {
		raw = inner
	}

	if status, ok := raw["status"].(string); ok {
		switch status {
		case "EDITED":
			return EventTypeMessageEdit
		case "REMOVED":
			return EventTypeMessageDelete
		}
	}

	return EventTypeMessageNew
}

func ResolveChat(frame *protocol.InboundFrame) EventType {
	return EventTypeChatUpdate
}

func ResolveMessageDelete(frame *protocol.InboundFrame) EventType {
	return EventTypeMessageDelete
}

func ResolveMessageRead(frame *protocol.InboundFrame) EventType {
	return EventTypeMessageRead
}

func ResolveTyping(frame *protocol.InboundFrame) EventType {
	return EventTypeTyping
}

func ResolvePresence(frame *protocol.InboundFrame) EventType {
	return EventTypePresence
}

func ResolveReactionUpdate(frame *protocol.InboundFrame) EventType {
	return EventTypeReactionUpdate
}

func ResolveAttach(frame *protocol.InboundFrame) EventType {
	if frame.Payload == nil {
		return ""
	}

	if _, ok := frame.Payload["fileId"]; ok {
		return EventTypeFileReady
	}

	if _, ok := frame.Payload["videoId"]; ok {
		return EventTypeVideoReady
	}

	return ""
}

type Resolver func(frame *protocol.InboundFrame) EventType

var EventMap = map[int]Resolver{
	int(protocol.OpcodeNotifMessage):             ResolveMessage,
	int(protocol.OpcodeMsgEdit):                  ResolveMessage,
	int(protocol.OpcodeNotifChat):                ResolveChat,
	int(protocol.OpcodeNotifMsgDelete):           ResolveMessageDelete,
	int(protocol.OpcodeNotifAttach):              ResolveAttach,
	int(protocol.OpcodeNotifTyping):              ResolveTyping,
	int(protocol.OpcodeNotifMark):                ResolveMessageRead,
	int(protocol.OpcodeNotifPresence):            ResolvePresence,
	int(protocol.OpcodeNotifMsgReactionsChanged): ResolveReactionUpdate,
}
