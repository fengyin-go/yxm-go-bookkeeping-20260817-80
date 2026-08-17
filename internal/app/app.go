// Package app 负责依赖装配。
package app

import (
	"net/http"

	"bookkeeping/internal/config"
	"bookkeeping/internal/handler"
	"bookkeeping/internal/service"
	"bookkeeping/internal/store"
	"bookkeeping/pkg/logger"
)

// App 应用容器。
type App struct {
	server *handler.Server
}

// New 装配全部依赖。
func New(cfg *config.Config, log *logger.Logger) (*App, error) {
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	server := handler.NewServer(svc, log, cfg)

	log.Infof("应用装配完成，配置：%s", cfg.String())
	return &App{server: server}, nil
}

// Routes 返回应用的 HTTP 处理器。
func (a *App) Routes() http.Handler {
	return a.server.Routes()
}
