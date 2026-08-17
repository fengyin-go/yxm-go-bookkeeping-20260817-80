package store

import "bookkeeping/internal/model"

// CreateBudget 新增预算。
func (s *MemoryStore) CreateBudget(b *model.Budget) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.budgets[b.ID] = b
	return nil
}

// GetBudget 按 ID 查询预算。
func (s *MemoryStore) GetBudget(id string) (*model.Budget, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.budgets[id]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

// ListBudgets 返回全部预算。
func (s *MemoryStore) ListBudgets() []*model.Budget {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Budget, 0, len(s.budgets))
	for _, b := range s.budgets {
		list = append(list, b)
	}
	return list
}

// UpdateBudget 覆盖保存预算。
func (s *MemoryStore) UpdateBudget(b *model.Budget) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.budgets[b.ID]; !ok {
		return ErrNotFound
	}
	s.budgets[b.ID] = b
	return nil
}

// DeleteBudget 按 ID 删除预算。
func (s *MemoryStore) DeleteBudget(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.budgets[id]; !ok {
		return ErrNotFound
	}
	delete(s.budgets, id)
	return nil
}

// HasBudget 判断同一分类（或全局）是否已经存在预算。
func (s *MemoryStore) HasBudget(categoryID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, b := range s.budgets {
		if b.CategoryID == categoryID {
			return true
		}
	}
	return false
}
