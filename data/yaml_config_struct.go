package data

type Config_struct struct {
	Client struct {
		Ipv4_address string `yaml:"ipv4_address"`
		Port         string `yaml:"port"`
	}
	Server struct {
		Ipv4_address string `yaml:"ipv4_address"`
		Port         string `yaml:"port"`
	}
}
