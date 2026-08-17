# Bug 是什么

多个 goroutine 同时给同一个账户记账时，账户余额更新不是原子的，出现丢失更新，并发安全检测会报告 data race。

# 如何触发

并发调用 `CreateTransaction` 给同一账户增加多笔收入，然后检查账户余额。

# 错误信息

`go test ./internal/service -run TestService_ConcurrentTransactionKeepsBalance -race` 会报告账户余额不等于期望值，并给出 `WARNING: DATA RACE`。
