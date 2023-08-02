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
	Accounts      []*Account    `mapstructure:"accounts"`
	Circuit       Circuit       `mapstructure:"circuit"`
}

type Server struct {
	Port string `mapstructure:"port"`
}

type Postgres struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}

func (postgres *Postgres) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei", postgres.Host, postgres.Username, postgres.Password, postgres.Name, postgres.Port)
}

func (postgres *Postgres) Url() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", postgres.Username, postgres.Password, postgres.Host, postgres.Port, postgres.Name)
}

type Redis struct {
	Host     string   `mapstructure:"host"`
	Password string   `mapstructure:"password"`
	Port     string   `mapstructure:"port"`
	Keys     Keys     `mapstructure:"keys"`
	Commands Commands `mapstructure:"commands"`
	Channels Channels `mapstructure:"channels"`
}

type Keys struct {
	LastInsertedKey string `mapstructure:"lastinsertedkey"`
	RollupedTxsKey  string `mapsturcture:"rollupedtxskey"`
}

type Commands struct {
	RollupCommand string `mapstructure:"rollupcommand"`
}

type Channels struct {
	RollupCh string `mapstructure:"rollupchannel"`
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

type Account struct {
	EcdsaPrivKey string `mapstructure:"ecdsaprivkey"`
	EddsaPrivKey string `mapstructure:"eddsaprivkey"`
	Index        int64  `mapstructure:"index"`
}

type Circuit struct {
	Path string `mapstructure:"path"`
}

func LoadConfig(node string, paths ...string) (config *Config, err error) {
	for _, path := range paths {
		viper.AddConfigPath(path)
	}

	yamlFileName := fmt.Sprintf("env.example.%s", node)
	viper.SetConfigType("yaml")
	viper.SetConfigName(yamlFileName)

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
