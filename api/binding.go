package api

import (
	"github.com/akovardin/gomax/api/core"
	"github.com/akovardin/gomax/types"
)

func BindMessage(app core.AppInterface, msg *types.Message, chatID types.FlexInt) {
	app.BindMessage(msg, chatID)
}
