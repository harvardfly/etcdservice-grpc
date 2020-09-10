package config

import (
	"github.com/go-ini/ini"
)

// AppConfig 应用的配置结构体
type AppConfig struct {
	*ServerConfig `json:"server" ini:"server"`
	*EtcdConfig   `json:"etcd" ini:"etcd"`
	*ClientConfig `json:"client" ini:"client"`
}

// ServerConfig web server配置
type ServerConfig struct {
	ServiceName string `json:"service_name" ini:"serviceName"`
	Port        int    `json:"port" ini:"port"`
}

// EtcdConfig etcd 配置
type EtcdConfig struct {
	Schema   string `json:"schema" ini:"schema"`
	EtcdAddr string `json:"etcd_addr" ini:"etcdAddr"`
	TTL      int64  `json:"ttl" ini:"ttl"`
}

// ClientConfig web client配置
type ClientConfig struct {
	Caller string `json:"caller" ini:"caller"`
	Callee string `json:"callee" ini:"callee"`
}

// Conf 定义了全局的配置文件实例
var Conf = new(AppConfig)

// InitFromIni 配置文件映射到全局变量 Conf
func InitFromIni(filename string) error {
	err := ini.MapTo(Conf, filename)
	if err != nil {
		panic(err)
	}
	return err
}
