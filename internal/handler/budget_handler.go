package handler

import (
	"net/http"
	"time"

	"bookkeeping/internal/model"
	"bookkeeping/pkg/httpx"
)

// registerBudgetRoutes 注册预算相关路由。
func (s *Server) registerBudgetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/budgets", s.createBudget)
	mux.HandleFunc("GET /api/budgets", s.listBudgets)
	mux.HandleFunc("GET /api/budgets/{id}", s.getBudget)
	mux.HandleFunc("PATCH /api/budgets/{id}", s.updateBudget)
	mux.HandleFunc("DELETE /api/budgets/{id}", s.deleteBudget)
	mux.HandleFunc("GET /api/budgets/check", s.checkBudgets)
}

type createBudgetRequest struct {
	CategoryID string `json:"category_id"` // 为空表示全部支出
	Amount     int64  `json:"amount"`      // 预算（分）
}

func (s *Server) createBudget(w http.ResponseWriter, r *http.Request) {
	var req createBudgetRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateBudget(model.Budget{CategoryID: req.CategoryID, Amount: req.Amount})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) listBudgets(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListBudgets()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) getBudget(w http.ResponseWriter, r *http.Request) {
	b, err := s.svc.GetBudget(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

type updateBudgetRequest struct {
	Amount int64 `json:"amount"`
}

func (s *Server) updateBudget(w http.ResponseWriter, r *http.Request) {
	var req updateBudgetRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.UpdateBudget(r.PathValue("id"), req.Amount)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) deleteBudget(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteBudget(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}

func (s *Server) checkBudgets(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.CheckBudgets(time.Now())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
