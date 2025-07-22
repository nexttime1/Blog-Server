package conf

type Ai struct {
	Enable    bool   `json:"enable" yaml:"enable"`
	SecretKey string `json:"secretKey" yaml:"secretKey"`
	Nickname  string `json:"nickname" yaml:"nickname"`
	Avatar    string `json:"avatar" yaml:"avatar"`
}
