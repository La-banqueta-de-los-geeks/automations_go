package configs

import (
	"github.com/eduardogpg/gonv"
)

type RedisConfig struct {
	Addr string
	Password string
	DB string
}

var redisconfiguration *RedisConfig

func init(){
	redisconfiguration = &RedisConfig{}
	redisconfiguration.Addr	= gonv.GetStringEnv("REDIS_ADDR", ":6379")
	redisconfiguration.Password	= gonv.GetStringEnv("REDIS_PASSWORD", "vurokrazia")
	redisconfiguration.DB	= gonv.GetStringEnv("REDIS_DB", "0")
}
func (this *RedisConfig) Url() *RedisConfig {
	return redisconfiguration
}
func GetRedisconfiguration() *RedisConfig  {
	return redisconfiguration
}