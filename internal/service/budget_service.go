package service

import (
	"time"

	"bookkeeping/internal/model"
	"bookkeeping/pkg/idgen"
)

// CreateBudget 创建预算。
func (s *Service) CreateBudget(input model.Budget) (*model.Budget, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if input.CategoryID != "" {
		if _, err := s.store.GetCategory(input.CategoryID); err != nil {
			return nil, err
		}
	}
	b := &model.Budget{
		ID:         idgen.Hex(),
		CategoryID: input.CategoryID,
		Period:     input.Period,
		Amount:     input.Amount,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := s.store.CreateBudget(b); err != nil {
		return nil, err
	}
	s.log.Infof("创建预算 amount=%d category=%s", b.Amount, b.CategoryID)
	return b, nil
}

// GetBudget 按 ID 查询预算。
func (s *Service) GetBudget(id string) (*model.Budget, error) {
	return s.store.GetBudget(id)
}

// ListBudgets 列出全部预算。
func (s *Service) ListBudgets() ([]*model.Budget, error) {
	return s.store.ListBudgets(), nil
}

// UpdateBudget 更新预算金额。
func (s *Service) UpdateBudget(id string, amount int64) (*model.Budget, error) {
	exist, err := s.store.GetBudget(id)
	if err != nil {
		return nil, err
	}
	exist.Amount = amount
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateBudget(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteBudget 删除预算。
func (s *Service) DeleteBudget(id string) error {
	if err := s.store.DeleteBudget(id); err != nil {
		return err
	}
	s.log.Infof("删除预算 %s", id)
	return nil
}

// BudgetStatus 预算执行状态。
type BudgetStatus struct {
	Budget   *model.Budget `json:"budget"`
	Spent    int64         `json:"spent"`     // 当月已支出（分）
	Remaining int64        `json:"remaining"` // 剩余额度（分）
	OverBudget bool        `json:"over_budget"`
}

// CheckBudgets 检查所有预算在当前月份的支出情况。
func (s *Service) CheckBudgets(now time.Time) ([]*BudgetStatus, error) {
	year, month := now.Year(), now.Month()
	result := make([]*BudgetStatus, 0)
	for _, b := range s.store.ListBudgets() {
		var spent int64
		for _, t := range s.store.ListTransactions() {
			if t.Type != model.TypeExpense {
				continue
			}
			if t.OccurredAt.Year() != year || t.OccurredAt.Month() != month {
				continue
			}
			if b.CategoryID != "" && t.CategoryID != b.CategoryID {
				continue
			}
			spent += t.Amount
		}
		result = append(result, &BudgetStatus{
			Budget:     b,
			Spent:      spent,
			Remaining:  b.Amount - spent,
			OverBudget: spent > b.Amount,
		})
	}
	return result, nil
}
