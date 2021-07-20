package pojo

type UserInfo struct {
	Uid      string `json:"uid"`      // 用户唯一标识
	UserName string `json:"userName"` // 用户名
	NickName string `json:"nickName"` // 用户昵称
	Avatar   string `json:"avatar"`   // 用户头像
	Sex      int    `json:"sex"`      //性别
	Mobile   string `json:"mobile"`   //手机号码
	MailBox  string `json:"mailBox"`  //邮箱
	RoleCode string `json:"roleCode"` // 账号角色 basic:普通员工 admin:管理员
}
