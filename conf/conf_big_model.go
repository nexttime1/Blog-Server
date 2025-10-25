package conf

type ModelOption struct {
	Label    string `yaml:"label" json:"label"`
	Value    string `yaml:"value" json:"value"`
	Disabled bool   `yaml:"disabled" json:"disabled"`
}

type Setting struct {
	Name      string `yaml:"name" json:"name"`
	Enable    bool   `yaml:"enable" json:"enable"`
	ApiKey    string `yaml:"api_key" json:"api_key"`
	ApiSecret string `yaml:"api_secret" json:"api_secret"`
	Title     string `yaml:"title" json:"title"`
	Prompt    string `yaml:"prompt" json:"prompt"`
	Slogan    string `yaml:"slogan" json:"slogan"`
}

type SessionSetting struct {
	ChatScope    int `yaml:"chat_scope" json:"chat_scope"`       // 对话的积分消耗
	SessionScope int `yaml:"session_scope" json:"session_scope"` // 会话的积分消耗
	DayScope     int `yaml:"day_scope" json:"day_scope"`         // 每日可以领取的积分
}

type BigModel struct {
	Setting        Setting        `yaml:"setting"`
	ModelList      []ModelOption  `yaml:"model_list"`
	SessionSetting SessionSetting `yaml:"session_setting"`
}
