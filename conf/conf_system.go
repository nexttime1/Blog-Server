package conf

import "fmt"

type System struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

func (s System) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}
