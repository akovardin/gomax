package users

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

func (s *Service) GetUsers(userIDs []int) ([]*types.User, error) {
	payload := map[string]interface{}{
		"userIds": userIDs,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeContactInfo), payload)
	if err != nil {
		return nil, err
	}
	if usersRaw, ok := frame.Payload["users"]; ok {
		users, err := core.ConvertStruct[[]*types.User](usersRaw)
		if err != nil {
			return nil, err
		}
		return *users, nil
	}
	users, err := core.ConvertStruct[[]*types.User](frame.Payload)
	if err != nil {
		return nil, err
	}
	return *users, nil
}

func (s *Service) GetUser(userID int) (*types.User, error) {
	payload := map[string]interface{}{
		"userIds": []int{userID},
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeContactInfo), payload)
	if err != nil {
		return nil, err
	}
	if usersRaw, ok := frame.Payload["users"]; ok {
		users, err := core.ConvertStruct[[]*types.User](usersRaw)
		if err != nil {
			return nil, err
		}
		if len(*users) > 0 {
			return (*users)[0], nil
		}
		return nil, nil
	}
	if userRaw, ok := frame.Payload["user"]; ok {
		return core.ConvertStruct[types.User](userRaw)
	}
	return core.RequirePayloadModel[types.User](frame.Payload)
}

func (s *Service) FetchUsers(userIDs []int) ([]*types.User, error) {
	return s.GetUsers(userIDs)
}

func (s *Service) SearchByPhone(phone string) (*types.User, error) {
	payload := map[string]interface{}{
		"phone": phone,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeContactInfoByPhone), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.User](frame.Payload)
}

func (s *Service) AddContact(contactID int) error {
	payload := map[string]interface{}{
		"contactId": contactID,
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeContactAdd), payload)
	return err
}

func (s *Service) RemoveContact(contactID int) error {
	payload := map[string]interface{}{
		"contactId": contactID,
		"action":    "remove",
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeContactUpdate), payload)
	return err
}

func (s *Service) ImportContacts(contacts []types.ContactInfo) ([]*types.User, error) {
	payload := map[string]interface{}{
		"contacts": contacts,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeContactList), payload)
	if err != nil {
		return nil, err
	}
	if usersRaw, ok := frame.Payload["users"]; ok {
		users, err := core.ConvertStruct[[]*types.User](usersRaw)
		if err != nil {
			return nil, err
		}
		return *users, nil
	}
	users, err := core.ConvertStruct[[]*types.User](frame.Payload)
	if err != nil {
		return nil, err
	}
	return *users, nil
}

func (s *Service) GetSessions() ([]*types.Session, error) {
	payload := map[string]interface{}{}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeSessionsInfo), payload)
	if err != nil {
		return nil, err
	}
	if sessionsRaw, ok := frame.Payload["sessions"]; ok {
		sessions, err := core.ConvertStruct[[]*types.Session](sessionsRaw)
		if err != nil {
			return nil, err
		}
		return *sessions, nil
	}
	sessions, err := core.ConvertStruct[[]*types.Session](frame.Payload)
	if err != nil {
		return nil, err
	}
	return *sessions, nil
}

func (s *Service) GetChatID(firstUserID int, secondUserID int) int {
	return firstUserID ^ secondUserID
}
