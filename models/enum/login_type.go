package enum

type LoginType int8

const (
	UserPwdLoginType LoginType = 1
	QQLoginType      LoginType = 2
	EmailLoginType   LoginType = 3
)
