package configuration

type RedisConfiguration struct {
	Address  string `json:"address" yaml:"address"`
	Password string `json:"password" yaml:"password"`
}
