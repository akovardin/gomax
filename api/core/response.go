package core

import (
	"encoding/json"

	"github.com/akovardin/gomax/protocol"
)

func PayloadItem(payload map[string]interface{}, key string, default_ interface{}) interface{} {
	if val, ok := payload[key]; ok {
		return val
	}
	return default_
}

func RequirePayloadModel[T any](payload map[string]interface{}) (*T, error) {
	return ConvertStruct[T](payload)
}

func ConvertStruct[T any](data interface{}) (*T, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result T
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func BuildAPIError(frame *protocol.InboundFrame) *ApiError {
	message := ""
	localizedMessage := ""
	title := ""
	if frame.Payload != nil {
		if m, ok := frame.Payload["message"].(string); ok {
			message = m
		}
		if lm, ok := frame.Payload["localizedMessage"].(string); ok {
			localizedMessage = lm
		}
		if t, ok := frame.Payload["title"].(string); ok {
			title = t
		}
	}
	return &ApiError{
		Opcode:           frame.Opcode,
		ErrorStr:         frame.Error,
		MessageStr:       message,
		LocalizedMessage: localizedMessage,
		Title:            title,
		Payload:          frame.Payload,
	}
}

func InvokeAPI(app AppInterface, opcode int, payload map[string]interface{}) (*protocol.InboundFrame, error) {
	timeout := app.Config().RequestTimeout
	frame, err := app.Invoke(opcode, payload, int(protocol.CommandRequest), timeout, false)
	if err != nil {
		return nil, err
	}
	if frame.Cmd == int(protocol.CommandError) || frame.Error != "" {
		return nil, BuildAPIError(frame)
	}
	return frame, nil
}

func BindMessage(app AppInterface, msg interface{}) {}
