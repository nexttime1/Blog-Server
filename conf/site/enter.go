package site

type SiteInfo struct {
	Title string `yaml:"title" json:"title"`
	Logo  string `yaml:"logo" json:"logo"`
	Beian string `yaml:"beian" json:"beian"`
	Mode  int8   `yaml:"mode" json:"mode"` //1 社区模式   2 博客模式

}
type Project struct {
	Title   string `yaml:"title" json:"title"`
	Icon    string `yaml:"icon" json:"icon"`
	WebPath string `yaml:"webPath" json:"webPath"`
}
type Seo struct {
	Keywords    string `yaml:"keywords" json:"keywords"`
	Description string `yaml:"description" json:"description"`
}
type About struct {
	SiteDate string `yaml:"siteDate" json:"siteDate"`
	QQ       string `yaml:"qq"  json:"qq"`
	WeChat   string `yaml:"weChat"  json:"weChat"`
	Gitee    string `yaml:"gitee" json:"gitee"`
	Bilibili string `yaml:"bilibili" json:"bilibili"`
	GitHub   string `yaml:"github" json:"gitHub"`
}

type Login struct {
	QQLogin          string `yaml:"qqLogin" json:"QQLogin"`
	UsernamePwdLogin string `yaml:"usernamePwdLogin"  json:"usernamePwdLogin"`
	EmailLogin       string `yaml:"emailLogin"  json:"emailLogin"`
	Captcha          bool   `yaml:"captcha" json:"captcha"` //验证码

}
type ComponentInfo struct {
	Title  string `yaml:"title" json:"title"`
	Enable bool   `yaml:"enable" json:"enable"`
}

type IndexRight struct {
	list []ComponentInfo `yaml:"list"  json:"list"`
}
type Article struct {
	NOExamine bool `yaml:"noExamine" json:"NOExamine"` //免审核

}
