package service

import (
	"sort"
	"time"

	"bookkeeping/internal/model"
	"bookkeeping/pkg/idgen"
)

// CreateCategory 创建收支分类。
func (s *Service) CreateCategory(input model.Category) (*model.Category, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c := &model.Category{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Type:      input.Type,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateCategory(c); err != nil {
		return nil, err
	}
	s.log.Infof("创建分类 %s(%s)", c.Name, c.Type)
	return c, nil
}

// GetCategory 按 ID 查询分类。
func (s *Service) GetCategory(id string) (*model.Category, error) {
	return s.store.GetCategory(id)
}

// ListCategories 列出全部分类，可按类型过滤。
func (s *Service) ListCategories(categoryType string) ([]*model.Category, error) {
	all := s.store.ListCategories()
	list := make([]*model.Category, 0, len(all))
	for _, c := range all {
		if categoryType == "" || c.Type == categoryType {
			list = append(list, c)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}

// UpdateCategory 更新分类。
func (s *Service) UpdateCategory(id string, input model.Category) (*model.Category, error) {
	exist, err := s.store.GetCategory(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Type = input.Type
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCategory(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteCategory 删除分类。
func (s *Service) DeleteCategory(id string) error {
	if err := s.store.DeleteCategory(id); err != nil {
		return err
	}
	s.log.Infof("删除分类 %s", id)
	return nil
}
