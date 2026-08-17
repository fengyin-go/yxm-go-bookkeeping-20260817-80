package store

import "bookkeeping/internal/model"

// CreateTransaction 新增流水。
func (s *MemoryStore) CreateTransaction(t *model.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions[t.ID] = t
	return nil
}

// GetTransaction 按 ID 查询流水。
func (s *MemoryStore) GetTransaction(id string) (*model.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.transactions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

// ListTransactions 返回全部流水。
func (s *MemoryStore) ListTransactions() []*model.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Transaction, 0, len(s.transactions))
	for _, t := range s.transactions {
		list = append(list, t)
	}
	return list
}

// DeleteTransaction 按 ID 删除流水。
func (s *MemoryStore) DeleteTransaction(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.transactions[id]; !ok {
		return ErrNotFound
	}
	delete(s.transactions, id)
	return nil
}

// HasTransactionsByCategory 判断分类是否已被流水使用。
func (s *MemoryStore) HasTransactionsByCategory(categoryID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.transactions {
		if t.CategoryID == categoryID {
			return true
		}
	}
	return false
}
