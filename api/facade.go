package api

import (
	"github.com/akovardin/gomax/api/auth"
	"github.com/akovardin/gomax/api/bots"
	"github.com/akovardin/gomax/api/chats"
	"github.com/akovardin/gomax/api/core"
	"github.com/akovardin/gomax/api/messages"
	"github.com/akovardin/gomax/api/self"
	"github.com/akovardin/gomax/api/session"
	"github.com/akovardin/gomax/api/uploads"
	"github.com/akovardin/gomax/api/users"
)

type Facade struct {
	Messages *messages.Service
	Chats    *chats.Service
	Users    *users.Service
	Self     *self.Service
	Auth     *auth.Service
	Session  *session.Service
	Uploads  *uploads.Service
	Bots     *bots.Service
}

func NewFacade(app core.AppInterface) *Facade {
	return &Facade{
		Messages: messages.NewService(app),
		Chats:    chats.NewService(app),
		Users:    users.NewService(app),
		Self:     self.NewService(app),
		Auth:     auth.NewService(app),
		Session:  session.NewService(app),
		Uploads:  uploads.NewService(app),
		Bots:     bots.NewService(app),
	}
}
