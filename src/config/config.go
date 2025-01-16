package config

import (
	"wqnx/src/wiface"
)

const (
	version                  = "1.0.0"
	defaultName              = "wqnx"
	defaultIP                = "127.0.0.1"
	defaultHost              = "localhost"
	defaultPort              = 6666
	defaultMaxPacketSize     = 4096
	defaultMaxConn           = 10000
	defaultMaxWorkerPoolSize = 10
	defaultMaxWorkerTaskLen  = 10000
	defaultMaxMsgChanBuffLen = 10000
)

type Config struct {
	QUICServer        wiface.IServer
	Name              string `json:"name" yaml:"name"` //当前服务的名称
	IP                string `json:"ip" yaml:"ip"`
	Host              string `json:"host" yaml:"host"`                                   //主机
	Port              int    `json:"port" yaml:"port"`                                   //端口号
	Version           string `json:"version" yaml:"version"`                             //当前框架的版本
	MaxPacketSize     uint32 `json:"max_packet_size" yaml:"max_packet_size"`             //数据包最大大小
	MaxConn           int    `json:"max_conn" yaml:"max_conn"`                           //最大连接数
	WorkerPoolSize    uint32 `json:"worker_pool_size" yaml:"worker_pool_size"`           //工作池的数量
	MaxWorkerTaskLen  uint32 `json:"max_worker_task_len" yaml:"max_worker_task_len"`     //worker对应的队列的最大长度
	MaxMsgChanBuffLen uint32 `json:"max_msg_chan_buff_len" yaml:"max_msg_chan_buff_len"` //最大有缓冲读写goroutine通道的缓冲区大小
}

type Options func(config *Config)

func WithName(name string) Options {
	return func(o *Config) {
		o.Name = name
	}
}

func WithIP(ip string) Options {
	return func(o *Config) {
		o.IP = ip
	}
}

func WithHost(host string) Options {
	return func(o *Config) {
		o.Host = host
	}
}

func WithPort(port int) Options {
	return func(o *Config) {
		o.Port = port
	}
}

func WithVersion(v string) Options {
	return func(o *Config) {
		o.Version = v
	}
}

func WithMaxPacketSize(size uint32) Options {
	return func(o *Config) {
		o.MaxPacketSize = size
	}
}

func WithMaxConn(max int) Options {
	return func(o *Config) {
		o.MaxConn = max
	}
}

func WithWorkerPoolSize(size uint32) Options {
	return func(o *Config) {
		o.WorkerPoolSize = size
	}
}

func WithMaxWorkerTaskLen(size uint32) Options {
	return func(o *Config) {
		o.MaxWorkerTaskLen = size
	}
}

func WithMaxMsgChanBuffLen(size uint32) Options {
	return func(o *Config) {
		o.MaxMsgChanBuffLen = size
	}
}

func NewOption(opts ...Options) *Config {
	o := &Config{}
	opts = append(defaultOptions(), opts...)
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func defaultOptions() []Options {
	return []Options{
		WithName(defaultName),
		WithIP(defaultIP),
		WithHost(defaultHost),
		WithPort(defaultPort),
		WithVersion(version),
		WithMaxPacketSize(defaultMaxPacketSize),
		WithMaxConn(defaultMaxConn),
		WithWorkerPoolSize(defaultMaxWorkerPoolSize),
		WithMaxWorkerTaskLen(defaultMaxWorkerTaskLen),
		WithMaxMsgChanBuffLen(defaultMaxMsgChanBuffLen),
	}
}

func NewDefault() *Config {
	return NewOption()
}

var Conf *Config

func init() {
	Conf = NewDefault()
}
