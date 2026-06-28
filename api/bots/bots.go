package bots

import (
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

func (s *Service) GetInitData(botID int, chatID int, startParam string) (*types.InitData, error) {
	payload := map[string]interface{}{
		"botId":      botID,
		"chatId":     chatID,
		"startParam": startParam,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeWebAppInitData), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.InitData](frame.Payload)
}
