package conf

type Gaode struct {
	Enable bool   `yaml:"enable" json:"enable"`
	Key    string `json:"key" yaml:"key"`
}
