package config

import (
	"sync"
)

var (
	configInst *Config   // config instance
	configOnce sync.Once // config once
)

// Config
type Config struct {
	Agent           bool   `json:"agent"`
	dbRetryTimes    int    `json:"db_retry_times"`
	mySQLConnection string `json:"mysql_connection"`
}

// New factory method
func New() *Config {
	return &Config{dbRetryTimes: 60}
}

// SetDBRetryTimes set database connect failed retry times
func (object *Config) SetDBRetryTimes(times int) *Config {
	object.dbRetryTimes = times
	return object
}

// GetDBRetryTimes get database connect failed retry times
func (object *Config) GetDBRetryTimes() int {
	return object.dbRetryTimes
}

// GetMySQLConnection get mysql connection string
func (object *Config) GetMySQLConnection() string {
	return object.mySQLConnection
}

// SetMySQLConnection set mysql connection string
func (object *Config) SetMySQLConnection(conn string) *Config {
	object.mySQLConnection = conn
	return object
}

// Inst singleton
func Inst() *Config {
	configOnce.Do(func() {
		configInst = New()
	})
	return configInst
}

// init method
func init() {
	Inst().SetMySQLConnection("dev:dev@(db)/dcn?charset=utf8&parseTime=True&loc=Local")
}
