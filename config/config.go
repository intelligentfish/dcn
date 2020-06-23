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
	MySQLConnection string `json:"mysql_connection"`
}

// New factory method
func New() *Config {
	return &Config{}
}

// GetMySQLConnection get mysql connection string
func (object *Config) GetMySQLConnection() string {
	return object.MySQLConnection
}

// SetMySQLConnection set mysql connection string
func (object *Config) SetMySQLConnection(conn string) *Config {
	object.MySQLConnection = conn
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
	Inst().SetMySQLConnection("dev:dev@/dcn?charset=utf8&parseTime=True&loc=Local")
}
