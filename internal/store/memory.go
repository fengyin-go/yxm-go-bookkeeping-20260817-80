package store

import (
	"sync"

	"bookkeeping/internal/model"
)

// MemoryStore 基于内存的 Store 实现。
type MemoryStore struct {
	mu           sync.RWMutex
	accounts     map[string]*model.Account
	categories   map[string]*model.Category
	transactions map[string]*model.Transaction
	budgets      map[string]*model.Budget
}

// NewMemoryStore 创建空的内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		accounts:     make(map[string]*model.Account),
		categories:   make(map[string]*model.Category),
		transactions: make(map[string]*model.Transaction),
		budgets:      make(map[string]*model.Budget),
	}
}

var _ Store = (*MemoryStore)(nil)
