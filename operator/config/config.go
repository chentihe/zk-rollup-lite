package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server        Server        `mapstructure:"server"`
	Postgres      Postgres      `mapstructure:"postgres"`
	Redis         Redis         `mapstructure:"redis"`
	SmartContract SmartContract `mapstructure:"smartcontract"`
	EthClient     EthClient     `mapstructure:"ethclient"`
	Sender        Sender        `mapstructure:"sender"`
	Recipient     Recipient     `mapstructure:"recipient"`
	Circuit       Circuit       `mapstructure:"circuit"`
}

type Server struct {
	Port string `mapstructure:"port"`
}

type Postgres struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Port     string `mapstructure:"port"`
}

func (postgres *Postgres) DSN() string {
	return fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei", postgres.Username, postgres.Password, postgres.Name, postgres.Port)
}

type Redis struct {
	Host       string `mapstructure:"host"`
	Password   string `mapstructure:"password"`
	Port       string `mapstructure:"port"`
	Prometheus string `mapstructure:"prometheus"`
}

func (redis *Redis) Addr() string {
	return fmt.Sprintf("%s:%s", redis.Host, redis.Port)
}

type SmartContract struct {
	Address string `mapstructure:"address"`
	Abi     string `mapstructure:"abi"`
}

type EthClient struct {
	RPCUrl     string `mapstructure:"rpcurl"`
	WSUrl      string `mapstructure:"wsurl"`
	PrivateKey string `mapstructure:"privatekey"`
}

type Sender struct {
	PrivateKey string `mapstructure:"privatekey"`
}

type Recipient struct {
	PrivateKey string `mapstructure:"privatekey"`
}

type Circuit struct {
	Path string `mapstructure:"path"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")
	viper.SetConfigName("env.example")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
