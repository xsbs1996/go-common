package msg

// Built 业务类型
const Built = "built" // 1-999错误码归内置错误

// 错误码区间分配
var codeRange = map[string][2]int{
	Built: {1, 999},
}
