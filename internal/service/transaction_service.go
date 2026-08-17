package service

import (
	"sort"
	"time"

	"bookkeeping/internal/model"
	"bookkeeping/pkg/idgen"
)

// CreateTransaction 记一笔账，同步更新账户余额。
func (s *Service) CreateTransaction(input model.Transaction) (*model.Transaction, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	account, err := s.store.GetAccount(input.AccountID)
	if err != nil {
		return nil, err
	}
	category, err := s.store.GetCategory(input.CategoryID)
	if err != nil {
		return nil, err
	}
	// 分类类型须与收支类型一致。
	if category.Type != input.Type {
		return nil, model.NewValidationError("category_id", "分类类型与收支类型不一致")
	}
	// 支出时校验余额充足。
	if input.Type == model.TypeExpense && account.Balance < input.Amount {
		return nil, model.NewValidationError("amount", "账户余额不足")
	}

	t := &model.Transaction{
		ID:         idgen.Hex(),
		AccountID:  input.AccountID,
		CategoryID: input.CategoryID,
		Type:       input.Type,
		Amount:     input.Amount,
		Note:       input.Note,
		OccurredAt: input.OccurredAt,
		CreatedAt:  time.Now(),
	}
	if err := s.store.ApplyTransaction(t, account.ID, t.BalanceDelta()); err != nil {
		return nil, err
	}

	s.log.Infof("记账 %s %d 分 (%s)", input.Type, input.Amount, input.CategoryID)
	return t, nil
}

// GetTransaction 按 ID 查询流水。
func (s *Service) GetTransaction(id string) (*model.Transaction, error) {
	return s.store.GetTransaction(id)
}

// ListTransactions 列出流水，支持筛选与分页，按发生时间倒序。
func (s *Service) ListTransactions(filter model.TransactionFilter, page, size int) ([]*model.Transaction, int, error) {
	all := s.store.ListTransactions()
	matched := make([]*model.Transaction, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].OccurredAt.After(matched[j].OccurredAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Transaction{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// DeleteTransaction 删除流水并回滚账户余额。
func (s *Service) DeleteTransaction(id string) error {
	t, err := s.store.RemoveTransaction(id)
	if err != nil {
		return err
	}
	_ = t
	s.log.Infof("删除流水 %s", id)
	return nil
}
