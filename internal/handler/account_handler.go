package handler

import (
	"net/http"

	"bookkeeping/internal/model"
	"bookkeeping/pkg/httpx"
)

// registerAccountRoutes 注册账户相关路由。
func (s *Server) registerAccountRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/accounts", s.createAccount)
	mux.HandleFunc("GET /api/accounts", s.listAccounts)
	mux.HandleFunc("GET /api/accounts/{id}", s.getAccount)
	mux.HandleFunc("PUT /api/accounts/{id}", s.updateAccount)
	mux.HandleFunc("DELETE /api/accounts/{id}", s.deleteAccount)
}

type createAccountRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Balance  int64  `json:"balance"`  // 初始余额（分）
	Currency string `json:"currency"`
}

func (s *Server) createAccount(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAccount(model.Account{
		Name:     req.Name,
		Type:     req.Type,
		Balance:  req.Balance,
		Currency: req.Currency,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAccounts(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListAccounts()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) getAccount(w http.ResponseWriter, r *http.Request) {
	a, err := s.svc.GetAccount(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type updateAccountRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (s *Server) updateAccount(w http.ResponseWriter, r *http.Request) {
	var req updateAccountRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateAccount(r.PathValue("id"), model.Account{Name: req.Name, Type: req.Type})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAccount(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteAccount(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
