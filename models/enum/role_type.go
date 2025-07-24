package enum

type RoleType int8

const (
	AdminRole   RoleType = 1
	UserRole    RoleType = 2
	VisitorRole RoleType = 3
	BlackRole   RoleType = 4
)

func (r RoleType) ParseRole() string {
	switch r {
	case AdminRole:
		return "管理员"
	case UserRole:
		return "普通用户"
	case VisitorRole:
		return "游客"
	case BlackRole:
		return "黑名单"
	}
	return "游客（未知）"
}
