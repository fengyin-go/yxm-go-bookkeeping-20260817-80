package handler

import (
	"net/http"
	"time"

	"bookkeeping/internal/model"
	"bookkeeping/pkg/httpx"
)

// registerTransactionRoutes 注册流水相关路由。
func (s *Server) registerTransactionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/transactions", s.createTransaction)
	mux.HandleFunc("GET /api/transactions", s.listTransactions)
	mux.HandleFunc("GET /api/transactions/{id}", s.getTransaction)
	mux.HandleFunc("DELETE /api/transactions/{id}", s.deleteTransaction)
}

type createTransactionRequest struct {
	AccountID  string `json:"account_id"`
	CategoryID string `json:"category_id"`
	Type       string `json:"type"`
	Amount     int64  `json:"amount"` // 金额（分）
	Note       string `json:"note"`
	OccurredAt string `json:"occurred_at"` // RFC3339，可选
}

func (s *Server) createTransaction(w http.ResponseWriter, r *http.Request) {
	var req createTransactionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	occurredAt := time.Now()
	if req.OccurredAt != "" {
		t, err := time.Parse(time.RFC3339, req.OccurredAt)
		if err != nil {
			httpx.BadRequest(w, "occurred_at 需为 RFC3339 时间")
			return
		}
		occurredAt = t
	}
	t, err := s.svc.CreateTransaction(model.Transaction{
		AccountID:  req.AccountID,
		CategoryID: req.CategoryID,
		Type:       req.Type,
		Amount:     req.Amount,
		Note:       req.Note,
		OccurredAt: occurredAt,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

// listTransactions 流水列表：GET /api/transactions?account_id=&category_id=&type=&from=&to=&page=&size=
func (s *Server) listTransactions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TransactionFilter{
		AccountID:  r.URL.Query().Get("account_id"),
		CategoryID: r.URL.Query().Get("category_id"),
		Type:       r.URL.Query().Get("type"),
	}
	if from := r.URL.Query().Get("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			filter.From = &t
		}
	}
	if to := r.URL.Query().Get("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			filter.To = &t
		}
	}
	items, total, err := s.svc.ListTransactions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTransaction(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.GetTransaction(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteTransaction(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
