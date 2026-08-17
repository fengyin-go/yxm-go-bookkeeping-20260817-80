package store

import "bookkeeping/internal/model"

// CreateCategory 新增分类，同名同类型重复时返回 ErrConflict。
func (s *MemoryStore) CreateCategory(c *model.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.categories {
		if exist.Name == c.Name && exist.Type == c.Type {
			return ErrConflict
		}
	}
	s.categories[c.ID] = c
	return nil
}

// GetCategory 按 ID 查询分类。
func (s *MemoryStore) GetCategory(id string) (*model.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.categories[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

// ListCategories 返回全部分类。
func (s *MemoryStore) ListCategories() []*model.Category {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Category, 0, len(s.categories))
	for _, c := range s.categories {
		list = append(list, c)
	}
	return list
}

// UpdateCategory 覆盖保存分类。
func (s *MemoryStore) UpdateCategory(c *model.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.categories[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.categories {
		if exist.ID != c.ID && exist.Name == c.Name && exist.Type == c.Type {
			return ErrConflict
		}
	}
	s.categories[c.ID] = c
	return nil
}

// DeleteCategory 按 ID 删除分类。
func (s *MemoryStore) DeleteCategory(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.categories[id]; !ok {
		return ErrNotFound
	}
	delete(s.categories, id)
	return nil
}
