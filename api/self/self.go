package self

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

func (s *Service) RequestProfilePhotoUploadURL() (string, error) {
	payload := map[string]interface{}{
		"avatar": false,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodePhotoUpload), payload)
	if err != nil {
		return "", err
	}
	if url, ok := frame.Payload["url"].(string); ok {
		return url, nil
	}
	return "", core.NewPyMaxError("no upload URL in response")
}

func (s *Service) ChangeProfile(firstName string, lastName string, description string, photo interface{}, photoToken string) (bool, error) {
	payload := map[string]interface{}{
		"firstName":   firstName,
		"lastName":    lastName,
		"description": description,
	}
	if photo != nil {
		payload["photo"] = photo
	}
	if photoToken != "" {
		payload["photoToken"] = photoToken
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeProfile), payload)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) CreateFolder(title string, chatInclude []int, filters []map[string]interface{}) (*types.FolderUpdate, error) {
	payload := map[string]interface{}{
		"title":   title,
		"include": chatInclude,
	}
	if filters != nil {
		payload["filters"] = filters
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeFoldersUpdate), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.FolderUpdate](frame.Payload)
}

func (s *Service) GetFolders(folderSync int) (*types.FolderList, error) {
	payload := map[string]interface{}{
		"folderSync": folderSync,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeFoldersGet), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.FolderList](frame.Payload)
}

func (s *Service) UpdateFolder(folderID int, title string, chatInclude []int, filters []map[string]interface{}) (*types.FolderUpdate, error) {
	payload := map[string]interface{}{
		"id":      folderID,
		"title":   title,
		"include": chatInclude,
	}
	if filters != nil {
		payload["filters"] = filters
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeFoldersUpdate), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.FolderUpdate](frame.Payload)
}

func (s *Service) DeleteFolder(folderID int) (*types.FolderUpdate, error) {
	payload := map[string]interface{}{
		"sourceId": folderID,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeFoldersDelete), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.FolderUpdate](frame.Payload)
}

func (s *Service) CloseAllSessions() (bool, error) {
	payload := map[string]interface{}{}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeSessionsClose), payload)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) Logout() (bool, error) {
	payload := map[string]interface{}{}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeLogout), payload)
	if err != nil {
		return false, err
	}
	return true, nil
}
