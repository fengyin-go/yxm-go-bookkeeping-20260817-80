package store

import (
	"time"

	"bookkeeping/internal/model"
)

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

// ApplyTransaction 在同一把锁内完成流水写入和账户余额更新，避免并发丢更新。
func (s *MemoryStore) ApplyTransaction(t *model.Transaction, accountID string, balanceDelta int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	account, ok := s.accounts[accountID]
	if !ok {
		return ErrNotFound
	}
	if _, ok := s.transactions[t.ID]; ok {
		return ErrConflict
	}

	s.transactions[t.ID] = t
	account.Balance += balanceDelta
	account.UpdatedAt = time.Now()
	return nil
}

// RemoveTransaction 在同一把锁内删除流水并回滚账户余额。
func (s *MemoryStore) RemoveTransaction(id string) (*model.Transaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, ok := s.transactions[id]
	if !ok {
		return nil, ErrNotFound
	}
	account, ok := s.accounts[tx.AccountID]
	if !ok {
		return nil, ErrNotFound
	}
	if tx.Type == model.TypeIncome {
		account.Balance -= tx.Amount
	} else {
		account.Balance += tx.Amount
	}
	account.UpdatedAt = time.Now()
	delete(s.transactions, id)
	return tx, nil
}
