# Bug 是什么

同一个分类（或全局）可以创建多个预算，预算检查结果中会出现重复的预算项。

# 如何触发

对同一个分类连续创建两个预算，然后调用预算检查接口。

# 错误信息

`go test ./internal/service -run TestService_DuplicateBudgetIsRejected` 会失败，期望返回重复预算错误，但实际得到 `nil`。
