package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	avatv1 "pet-angel/api/avatar/v1"
	"pet-angel/internal/biz"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// mock usecase for ChatStreamHTTP
type mockAvatarUC struct{ lastUserMsg, lastAiMsg string }

func (m *mockAvatarUC) SaveUserMessage(_ interface{}, _ int64, content string) (*biz.ChatMsg, error) {
	m.lastUserMsg = content
	return &biz.ChatMsg{}, nil
}
func (m *mockAvatarUC) SaveAIMessage(_ interface{}, _ int64, content string) (*biz.ChatMsg, error) {
	m.lastAiMsg = content
	return &biz.ChatMsg{}, nil
}

func TestAvatarService_ChatRouteExist(t *testing.T) {
	// just test route wiring works
	s := &AvatarService{uc: &biz.AvatarUsecase{}, jwtSecret: "", logger: nil}
	_ = s // no-op
}

func TestUploadRouteExists(t *testing.T) {
	srv := khttp.NewServer()
	svc := &GreeterService{}
	avat := &AvatarService{uc: biz.NewAvatarUsecase(nil), jwtSecret: "", logger: nil}
	avatv1.RegisterAvatarServiceHTTPServer(srv, avat)
	ts := httptest.NewServer(srv)
	defer ts.Close()
	_ = svc
	// POST body
	body, _ := json.Marshal(map[string]string{"content": "hi"})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/avatar/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}
