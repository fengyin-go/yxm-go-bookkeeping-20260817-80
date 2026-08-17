package service

import (
	"sync"
	"testing"
	"time"

	"bookkeeping/internal/model"
)

func TestService_ConcurrentTransactionKeepsBalance(t *testing.T) {
	s := newTestService()
	aid, incID, _ := s.seed(t)

	const workers = 8
	const rounds = 12

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < rounds; i++ {
				_, err := s.CreateTransaction(model.Transaction{
					AccountID:  aid,
					CategoryID: incID,
					Type:       model.TypeIncome,
					Amount:     1,
					OccurredAt: time.Now(),
				})
				if err != nil {
					t.Errorf("并发记账失败: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	a, err := s.GetAccount(aid)
	if err != nil {
		t.Fatalf("查询账户失败: %v", err)
	}
	want := int64(100000 + workers*rounds)
	if a.Balance != want {
		t.Fatalf("并发记账后余额错误，期望 %d，得到 %d", want, a.Balance)
	}
}
