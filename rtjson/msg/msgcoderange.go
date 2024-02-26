package msg

// 业务类型
const (
	Built = "built" // 1-999错误码归内置错误

	UserRpc     = "user.rpc"     // 1000-1999 错误码归 user.rpc
	PayRpc      = "pay.rpc"      // 2000-2999 错误码归 pay.rpc
	WalletRpc   = "wallet.rpc"   // 3000-3999 错误码归 wallet.rpc
	RiskRpc     = "risk.rpc"     // 4000-4999 错误码归 risk.rpc
	CommonRpc   = "common.rpc"   // 5000-5999 错误码归 common.rpc
	ActivityRpc = "activity.rpc" // 6000-6999 错误码归 activity.rpc

	AdminHttp     = "admin.http"      // 10000-10999 错误码归 admin.http
	HomeHttp      = "home.http"       // 11000-11999 错误码归 home.http
	UserHttp      = "user.http"       // 12000-12999 错误码归 user.http
	OrderHttp     = "order.http"      // 13000-13999 错误码归 order.http
	PayHttp       = "pay.http"        // 14000-14999 错误码归 pay.http
	PayNotifyHttp = "pay.notify.http" // 15000-15999 错误码归 pay.http
	GameHttp      = "game.http"       // 16000-16999 错误码归 game.http
	ActivityHttp  = "activity.http"   // 17000-17999 错误码归 activity.http

	JiliGameHttp = "jiligame.http" // 20000-20999 错误码归 jiligame.http
	FcGameHttp   = "fcgame.http"   // 21000-21999 错误码归 jiligame.http
)

// 错误码区间分配
var codeRange = map[string][2]int{
	Built: {1, 999},

	UserRpc:     {1000, 1999},
	PayRpc:      {2000, 2999},
	WalletRpc:   {3000, 3999},
	RiskRpc:     {4000, 4999},
	CommonRpc:   {5000, 5999},
	ActivityRpc: {6000, 6999},

	AdminHttp:     {10000, 10999},
	HomeHttp:      {11000, 11999},
	UserHttp:      {12000, 12999},
	OrderHttp:     {13000, 13999},
	PayHttp:       {14000, 14999},
	PayNotifyHttp: {15000, 15999},
	GameHttp:      {16000, 16999},
	ActivityHttp:  {17000, 17999},

	JiliGameHttp: {20000, 20999},
	FcGameHttp:   {21000, 21999},
}
