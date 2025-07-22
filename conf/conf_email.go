package conf

type Email struct {
	Domain       string `yaml:"domain" json:"domain"`
	Port         int    `yaml:"port" json:"port"`
	SendEmail    bool   `yaml:"sendEmail" json:"sendEmail"`
	AuthCode     string `yaml:"authCode" json:"authCode"` //授权码
	SendNickname string `yaml:"sendNickname" json:"sendNickname"`
	SSL          bool   `yaml:"ssl" json:"ssl"`
	TLS          bool   `yaml:"tls" json:"tls"`
}
