package utils

const (
	NoErr               = 0
	CustomErr           = 1
	TaskRewardPayErr01  = 10001
	TaskRewardPayErr02  = 10002
	TaskRewardPayErr03  = 10003
	InsufficientBalance = 20004
)

var ErrText = map[int]string{
	TaskRewardPayErr01:  "已成功给%d个成员发放奖金，有%d个员工奖励发放失败，请重新发放",
	TaskRewardPayErr02:  "已成功给%d个成员发放奖金，有%d个员工未绑定微信账号，无法发放奖金，快提醒员工绑定吧～",
	TaskRewardPayErr03:  "没有可发放的员工",
	InsufficientBalance: "账户余额不足",
}
