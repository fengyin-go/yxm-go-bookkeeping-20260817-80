// Package service 实现记账系统的业务逻辑层。
package service

import (
	"bookkeeping/internal/config"
	"bookkeeping/internal/store"
	"bookkeeping/pkg/logger"
)

// Service 业务逻辑层入口。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

// New 创建业务服务实例。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
