package configs

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13306)/cuddle",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
