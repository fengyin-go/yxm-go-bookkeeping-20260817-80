package service

import (
	"sort"
	"time"

	"bookkeeping/internal/model"
)

// Overview 账户总览统计。
type Overview struct {
	AccountCount     int   `json:"account_count"`
	TotalBalance     int64 `json:"total_balance"`      // 全部账户余额合计（分）
	TransactionCount int   `json:"transaction_count"`
}

// Overview 汇总账户与流水总览。
func (s *Service) Overview() (*Overview, error) {
	ov := &Overview{AccountCount: len(s.store.ListAccounts())}
	for _, a := range s.store.ListAccounts() {
		ov.TotalBalance += a.Balance
	}
	ov.TransactionCount = len(s.store.ListTransactions())
	return ov, nil
}

// MonthlyReport 月度收支报表。
type MonthlyReport struct {
	Year    int               `json:"year"`
	Month   int               `json:"month"`
	Income  int64             `json:"income"`  // 收入合计（分）
	Expense int64             `json:"expense"` // 支出合计（分）
	Net     int64             `json:"net"`     // 净收入（分）
	ByCategory []*CategorySum `json:"by_category"`
}

// CategorySum 分类汇总。
type CategorySum struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Type         string `json:"type"`
	Total        int64  `json:"total"` // 金额（分）
}

// MonthlyReport 生成指定月份的收支报表。
func (s *Service) MonthlyReport(year, month int) (*MonthlyReport, error) {
	if month < 1 || month > 12 {
		return nil, model.NewValidationError("month", "月份需在 1-12 之间")
	}
	report := &MonthlyReport{Year: year, Month: month}

	categoryName := make(map[string]string)
	for _, c := range s.store.ListCategories() {
		categoryName[c.ID] = c.Name
	}

	sumByCategory := make(map[string]*CategorySum)
	for _, t := range s.store.ListTransactions() {
		if t.OccurredAt.Year() != year || int(t.OccurredAt.Month()) != month {
			continue
		}
		switch t.Type {
		case model.TypeIncome:
			report.Income += t.Amount
		case model.TypeExpense:
			report.Expense += t.Amount
		}
		sum, ok := sumByCategory[t.CategoryID]
		if !ok {
			sum = &CategorySum{
				CategoryID:   t.CategoryID,
				CategoryName: categoryName[t.CategoryID],
				Type:         t.Type,
			}
			sumByCategory[t.CategoryID] = sum
		}
		sum.Total += t.Amount
	}
	report.Net = report.Income - report.Expense

	report.ByCategory = make([]*CategorySum, 0, len(sumByCategory))
	for _, sum := range sumByCategory {
		report.ByCategory = append(report.ByCategory, sum)
	}
	sort.Slice(report.ByCategory, func(i, j int) bool {
		return report.ByCategory[i].Total > report.ByCategory[j].Total
	})
	return report, nil
}

// currentMonth 返回当前年份与月份。
func currentMonth() (int, int) {
	now := time.Now()
	return now.Year(), int(now.Month())
}
