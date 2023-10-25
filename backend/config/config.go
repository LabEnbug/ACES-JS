package config

type DatabaseConfig struct {
	Host string
	Port int
	User string
	Pass string
	Name string
}

var DefaultDatabase DatabaseConfig

type RedisConfig struct {
	Host    string
	Port    int
	Pass    string
	Channel int
}

var DefaultRedis RedisConfig

func InitConfig() {
	DefaultDatabase = DatabaseConfig{
		Host: "127.0.0.1",
		Port: 3306,
		User: "root",
		Pass: "acesaces",
		Name: "qiniu_sv",
	}
	DefaultRedis = RedisConfig{
		Host:    "127.0.0.1",
		Port:    6379,
		Pass:    "acesaces",
		Channel: 0,
	}

}
