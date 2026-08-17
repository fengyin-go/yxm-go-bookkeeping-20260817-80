// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"bookkeeping/internal/config"
	"bookkeeping/internal/model"
	"bookkeeping/internal/service"
	"bookkeeping/internal/store"
	"bookkeeping/pkg/httpx"
	"bookkeeping/pkg/logger"
)

// Server 聚合全部 HTTP 处理器。
type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

// NewServer 创建处理器集合。
func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

// Routes 组装全部路由并挂载中间件。
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerAccountRoutes(mux)
	s.registerCategoryRoutes(mux)
	s.registerTransactionRoutes(mux)
	s.registerBudgetRoutes(mux)
	s.registerStatsRoutes(mux)
	return s.loggingMiddleware(s.recoveryMiddleware(mux))
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// writeServiceError 将业务错误映射为 HTTP 响应。
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case model.IsCategoryInUse(err):
		httpx.Conflict(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
