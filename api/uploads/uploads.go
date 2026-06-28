package uploads

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/akovardin/gomax/api/core"
	"github.com/akovardin/gomax/protocol"
)

type AttachPhotoPayload struct {
	Type       string `json:"type"`
	PhotoToken string `json:"photoToken"`
}

type AttachVideoPayload struct {
	Type    string `json:"type"`
	VideoID int    `json:"videoId"`
	Token   string `json:"token"`
}

type AttachFilePayload struct {
	Type   string `json:"type"`
	FileID int    `json:"fileId"`
}

type Service struct {
	app                core.AppInterface
	mu                 sync.Mutex
	videoUploadWaiters map[int]chan struct{}
	fileUploadWaiters  map[int]chan struct{}
}

func NewService(app core.AppInterface) *Service {
	return &Service{
		app:                app,
		videoUploadWaiters: make(map[int]chan struct{}),
		fileUploadWaiters:  make(map[int]chan struct{}),
	}
}

func (s *Service) RegisterVideoWaiter(videoID int) chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan struct{}, 1)
	s.videoUploadWaiters[videoID] = ch
	return ch
}

func (s *Service) SignalVideoReady(videoID int) {
	s.mu.Lock()
	ch, ok := s.videoUploadWaiters[videoID]
	if ok {
		delete(s.videoUploadWaiters, videoID)
	}
	s.mu.Unlock()
	if ok {
		ch <- struct{}{}
	}
}

func (s *Service) RegisterFileWaiter(fileID int) chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan struct{}, 1)
	s.fileUploadWaiters[fileID] = ch
	return ch
}

func (s *Service) SignalFileReady(fileID int) {
	s.mu.Lock()
	ch, ok := s.fileUploadWaiters[fileID]
	if ok {
		delete(s.fileUploadWaiters, fileID)
	}
	s.mu.Unlock()
	if ok {
		ch <- struct{}{}
	}
}

func (s *Service) UploadPhoto(photo io.Reader, filename string, profile bool) (*AttachPhotoPayload, error) {
	payload := map[string]interface{}{
		"avatar": profile,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodePhotoUpload), payload)
	if err != nil {
		return nil, err
	}

	uploadURL, _ := frame.Payload["url"].(string)
	if uploadURL == "" {
		return nil, core.NewUploadError("no upload URL received")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, core.NewUploadError("failed to create form file: " + err.Error())
	}
	if _, err := io.Copy(part, photo); err != nil {
		return nil, core.NewUploadError("failed to copy photo data: " + err.Error())
	}
	if err := writer.Close(); err != nil {
		return nil, core.NewUploadError("failed to close multipart writer: " + err.Error())
	}

	req, err := http.NewRequest("POST", uploadURL, &buf)
	if err != nil {
		return nil, core.NewUploadError("failed to create upload request: " + err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, core.NewUploadError("upload request failed: " + err.Error())
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, core.NewUploadError("failed to parse upload response: " + err.Error())
	}

	photoToken, _ := result["photoToken"].(string)

	return &AttachPhotoPayload{
		Type:       "PHOTO",
		PhotoToken: photoToken,
	}, nil
}

func (s *Service) UploadVideo(video io.Reader, filename string, size int64) (*AttachVideoPayload, error) {
	payload := map[string]interface{}{
		"fileName": filename,
		"fileSize": size,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeVideoUpload), payload)
	if err != nil {
		return nil, err
	}

	uploadURL, _ := frame.Payload["url"].(string)
	videoIDFloat, _ := frame.Payload["videoId"].(float64)
	videoID := int(videoIDFloat)

	if uploadURL == "" {
		return nil, core.NewUploadError("no upload URL received")
	}

	data, err := io.ReadAll(video)
	if err != nil {
		return nil, core.NewUploadError("failed to read video data: " + err.Error())
	}

	req, err := http.NewRequest("POST", uploadURL, bytes.NewReader(data))
	if err != nil {
		return nil, core.NewUploadError("failed to create upload request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, core.NewUploadError("upload request failed: " + err.Error())
	}
	defer resp.Body.Close()

	waiter := s.RegisterVideoWaiter(videoID)

	var token string
	select {
	case <-waiter:
		result := map[string]interface{}{}
		json.NewDecoder(resp.Body).Decode(&result)
		token, _ = result["token"].(string)
	}

	return &AttachVideoPayload{
		Type:    "VIDEO",
		VideoID: videoID,
		Token:   token,
	}, nil
}

func (s *Service) UploadFile(file io.Reader, filename string, size int64) (*AttachFilePayload, error) {
	payload := map[string]interface{}{
		"fileName": filename,
		"fileSize": size,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeFileUpload), payload)
	if err != nil {
		return nil, err
	}

	uploadURL, _ := frame.Payload["url"].(string)
	fileIDFloat, _ := frame.Payload["fileId"].(float64)
	fileID := int(fileIDFloat)

	if uploadURL == "" {
		return nil, core.NewUploadError("no upload URL received")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, core.NewUploadError("failed to read file data: " + err.Error())
	}

	req, err := http.NewRequest("POST", uploadURL, bytes.NewReader(data))
	if err != nil {
		return nil, core.NewUploadError("failed to create upload request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, core.NewUploadError("upload request failed: " + err.Error())
	}
	defer resp.Body.Close()

	waiter := s.RegisterFileWaiter(fileID)

	select {
	case <-waiter:
	}

	return &AttachFilePayload{
		Type:   "FILE",
		FileID: fileID,
	}, nil
}
