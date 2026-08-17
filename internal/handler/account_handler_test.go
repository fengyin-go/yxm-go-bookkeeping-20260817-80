package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bookkeeping/internal/config"
	"bookkeeping/internal/service"
	"bookkeeping/internal/store"
	"bookkeeping/pkg/logger"
)

func newTestServer() *Server {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return NewServer(service.New(store.NewMemoryStore(), log, cfg), log, cfg)
}

func decodeCode(t *testing.T, body []byte) int {
	t.Helper()
	var resp struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("解析响应失败: %v body=%s", err, body)
	}
	return resp.Code
}

// 已有流水的账户删除时应返回 409，且账户与流水保留。
func TestHandler_DeleteAccountInUseReturnsConflict(t *testing.T) {
	srv := newTestServer()
	mux := srv.Routes()

	// 创建账户。
	aID := createAccountOK(t, mux, `{"name":"工资卡","type":"bank","balance":100000}`)
	// 创建支出分类。
	cID := createCategoryOK(t, mux, `{"name":"餐饮","type":"expense"}`)
	// 记一笔支出流水。
	createTransactionOK(t, mux, `{"account_id":"`+aID+`","category_id":"`+cID+`","type":"expense","amount":100}`)

	// 删除被流水引用的账户。
	req := httptest.NewRequest(http.MethodDelete, "/api/accounts/"+aID, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("期望 409 Conflict，得到 %d body=%s", rec.Code, rec.Body.String())
	}
	if code := decodeCode(t, rec.Body.Bytes()); code != 409 {
		t.Fatalf("期望响应 code=409，得到 %d", code)
	}

	// 账户仍应存在，余额未被破坏。
	getRec := getAccount(t, mux, aID)
	if getRec.Code != http.StatusOK {
		t.Fatalf("拒绝删除后账户应仍存在，得到 %d", getRec.Code)
	}
}

// 无流水的账户可以正常删除。
func TestHandler_DeleteAccountWithoutTransactions(t *testing.T) {
	srv := newTestServer()
	mux := srv.Routes()

	aID := createAccountOK(t, mux, `{"name":"零钱","type":"cash","balance":0}`)

	req := httptest.NewRequest(http.MethodDelete, "/api/accounts/"+aID, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("期望 200，得到 %d body=%s", rec.Code, rec.Body.String())
	}
	if getAccount(t, mux, aID).Code != http.StatusNotFound {
		t.Fatalf("删除后账户应返回 404")
	}
}

// createAccountOK 创建账户并返回其 ID。
func createAccountOK(t *testing.T, mux http.Handler, body string) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/accounts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("创建账户失败: %d body=%s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析账户响应失败: %v", err)
	}
	return resp.Data.ID
}

// createCategoryOK 创建分类并返回其 ID。
func createCategoryOK(t *testing.T, mux http.Handler, body string) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("创建分类失败: %d body=%s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析分类响应失败: %v", err)
	}
	return resp.Data.ID
}

// createTransactionOK 记一笔流水。
func createTransactionOK(t *testing.T, mux http.Handler, body string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("创建流水失败: %d body=%s", rec.Code, rec.Body.String())
	}
}

// getAccount 查询账户详情。
func getAccount(t *testing.T, mux http.Handler, id string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/accounts/"+id, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}
