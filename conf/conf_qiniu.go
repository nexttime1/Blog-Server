package conf

type QiNiu struct {
	Enable    bool   `json:"enable" yaml:"enable"`
	AccessKey string `json:"accessKey" yaml:"accessKey"`
	SecretKey string `json:"secretKey" yaml:"secretKey"`
	Bucket    string `json:"bucket" yaml:"bucket"`
	Uri       string `json:"uri" yaml:"uri"`
	Region    string `json:"region" yaml:"region"`
	Prefix    string `json:"prefix" yaml:"prefix"`
	Size      int    `json:"size" yaml:"size"` //大小限制  MB

}
