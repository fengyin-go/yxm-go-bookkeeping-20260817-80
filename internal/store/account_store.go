package store

import "bookkeeping/internal/model"

// CreateAccount 新增账户，名称重复时返回 ErrConflict。
func (s *MemoryStore) CreateAccount(a *model.Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.accounts {
		if exist.Name == a.Name {
			return ErrConflict
		}
	}
	s.accounts[a.ID] = a
	return nil
}

// GetAccount 按 ID 查询账户。
func (s *MemoryStore) GetAccount(id string) (*model.Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.accounts[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

// GetAccountByName 按名称查询账户。
func (s *MemoryStore) GetAccountByName(name string) (*model.Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.accounts {
		if a.Name == name {
			return a, nil
		}
	}
	return nil, ErrNotFound
}

// ListAccounts 返回全部账户。
func (s *MemoryStore) ListAccounts() []*model.Account {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Account, 0, len(s.accounts))
	for _, a := range s.accounts {
		list = append(list, a)
	}
	return list
}

// UpdateAccount 覆盖保存账户。
func (s *MemoryStore) UpdateAccount(a *model.Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.accounts[a.ID]; !ok {
		return ErrNotFound
	}
	s.accounts[a.ID] = a
	return nil
}

// DeleteAccount 按 ID 删除账户。
func (s *MemoryStore) DeleteAccount(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.accounts[id]; !ok {
		return ErrNotFound
	}
	delete(s.accounts, id)
	return nil
}
