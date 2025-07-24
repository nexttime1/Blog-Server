package enum

type RegisterType int8

const (
	SignQQ    RegisterType = 1
	SignGitee RegisterType = 2
	SignEmail RegisterType = 3
)

func (r RegisterType) ParseRegister() string {
	switch r {
	case SignQQ:
		return "QQ"
	case SignGitee:
		return "Gitee"
	case SignEmail:
		return "Email"
	}
	return "其他"
}
