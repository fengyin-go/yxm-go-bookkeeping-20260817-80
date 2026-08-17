package service

import (
	"sort"
	"time"

	"bookkeeping/internal/model"
	"bookkeeping/pkg/idgen"
)

// CreateAccount 创建账户。
func (s *Service) CreateAccount(input model.Account) (*model.Account, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a := &model.Account{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Type:      input.Type,
		Balance:   input.Balance,
		Currency:  input.Currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.store.CreateAccount(a); err != nil {
		return nil, err
	}
	s.log.Infof("创建账户 %s", a.Name)
	return a, nil
}

// GetAccount 按 ID 查询账户。
func (s *Service) GetAccount(id string) (*model.Account, error) {
	return s.store.GetAccount(id)
}

// ListAccounts 列出全部账户。
func (s *Service) ListAccounts() ([]*model.Account, error) {
	list := s.store.ListAccounts()
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}

// UpdateAccount 更新账户名称与类型。
func (s *Service) UpdateAccount(id string, input model.Account) (*model.Account, error) {
	exist, err := s.store.GetAccount(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Type = input.Type
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateAccount(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteAccount 删除账户。
// 账户若已被流水引用则禁止删除，否则会留下指向已删除账户的悬空流水，
// 且这些流水当初造成的余额变动无法对账。
func (s *Service) DeleteAccount(id string) error {
	if _, err := s.store.GetAccount(id); err != nil {
		return err
	}
	for _, t := range s.store.ListTransactions() {
		if t.AccountID == id {
			return &model.AccountInUseError{AccountID: id}
		}
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	s.log.Infof("删除账户 %s", id)
	return nil
}
