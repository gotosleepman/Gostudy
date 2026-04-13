package config

var Config *Conf

type Conf struct {
	Mysql     MysqlConfig
	Redis     RedisConfig
	TestNet   NetConfig
	MainNet   NetConfig
	Threshold ThresholdConfig
	Env       EnvConfig
}

type EnvConfig struct {
	Port       string `toml:"port"`
	Version    string `toml:"version"`
	Protocol   string `toml:"protocol"`
	DomainName string `toml:"domain_name"`
}

type ThresholdConfig struct {
	LendingPoolNativeThreshold string `toml:"lending_pool_native_threshold"`
}

type MysqlConfig struct {
	Address      string `toml:"address"`
	Port         string `toml:"port"`
	DbName       string `toml:"db_name"`
	UserName     string `toml:"user_name"`
	Password     string `toml:"password"`
	MaxOpenConns int    `toml:"max_open_conns"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxLifeTime  int    `toml:"max_life_time"`
}

type NetConfig struct {
	ChainId         string `toml:"chain_id"`
	NetUrl          string `toml:"net_url"`
	LendingPoolAddr string `toml:"lending_pool_addr"`
}

type RedisConfig struct {
	Address     string `toml:"address"`
	Port        string `toml:"port"`
	Db          int    `toml:"db"`
	Password    string `toml:"password"`
	MaxIdle     int    `toml:"max_idle"`
	MaxActive   int    `toml:"max_active"`
	IdleTimeout int    `toml:"idle_timeout"`
}
