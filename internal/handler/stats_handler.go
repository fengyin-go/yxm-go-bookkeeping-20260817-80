package handler

import (
	"net/http"
	"strconv"
	"time"

	"bookkeeping/pkg/httpx"
)

// registerStatsRoutes 注册统计相关路由。
func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.overview)
	mux.HandleFunc("GET /api/stats/monthly", s.monthlyReport)
}

func (s *Server) overview(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.Overview()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

// monthlyReport 月度报表：GET /api/stats/monthly?year=2026&month=8
func (s *Server) monthlyReport(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	if v := r.URL.Query().Get("year"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			year = n
		}
	}
	if v := r.URL.Query().Get("month"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			month = n
		}
	}
	result, err := s.svc.MonthlyReport(year, month)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
