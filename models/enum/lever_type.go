package enum

type LevelType int8

const (
	LogInfoLevel LevelType = 1
	LogWainLevel LevelType = 2
	LogErrLevel  LevelType = 3
)

// 实现fmt.Stringer 接口。  直接%s
func (level LevelType) String() string {
	switch level {
	case LogInfoLevel:
		return "Info"
	case LogWainLevel:
		return "Wain"
	case LogErrLevel:
		return "Error"
	}
	return ""
}
