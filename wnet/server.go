package wnet

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"os"
	"os/signal"
)

const (
	version = "1.0.0"
)

type Server struct {
	Name    string
	Version string
	IP      string
	Host    string
	Port    int
	TlsConf *tls.Config
	Ctx     context.Context
}

type Options func(opts *Server)

func WithIP(ip string) Options {
	return func(opts *Server) {
		opts.IP = ip
	}
}

func WithHost(host string) Options {
	return func(opts *Server) {
		opts.Host = host
	}
}

func WithPort(port int) Options {
	return func(opts *Server) {
		opts.Port = port
	}
}

func WithName(name string) Options {
	return func(opts *Server) {
		opts.Name = name
	}
}

func WithVersion(version string) Options {
	return func(opts *Server) {
		opts.Version = version
	}
}

func WithTlsConf(conf *tls.Config) Options {
	return func(opts *Server) {
		opts.TlsConf = conf
	}
}

func NewDefaultServer() *Server {
	ip := "127.0.0.1"
	host := "localhost"
	opts := []Options{
		WithName("wqx"),
		WithVersion(version),
		WithPort(6666),
		WithIP(ip),
		WithHost(host),
		WithTlsConf(generateTLSConfig(ip, host)),
	}
	opt := &Server{
		Ctx: context.Background(),
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

func NewServer(opts ...Options) *Server {
	opt := NewDefaultServer()
	for _, o := range opts {
		o(opt)
	}
	return opt
}

func (s *Server) getAddress() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

func (s *Server) Start() {
	SysPrintInfo(fmt.Sprintf("start server listener at ip: %s, port: %d", s.IP, s.Port))
	SysPrintInfo(fmt.Sprintf("server name: %s, address: %s, version: %s ", s.Name, s.getAddress(), s.Version))

	go func() {
		addr := s.getAddress()
		listener, err := quic.ListenAddr(addr, generateTLSConfig(s.IP, s.Host), &quic.Config{
			EnableDatagrams: true,
		})
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = listener.Close()
		}()

		//监听成功
		SysPrintInfo(fmt.Sprintf("listen tcp success, version: %v, addr: %v", s.Version, addr))
		connId := uint32(1)
		for {
			//获取新的连接
			conn, err := listener.Accept(s.Ctx)
			if err != nil {
				SysPrintError(fmt.Sprintf("accept stream error: %v", err))
				continue
			}
			qconn := NewConnection(s.Ctx, conn, connId)
			SysPrintInfo(fmt.Sprintf("accept stream success, connId: %v", connId))
			connId++
			// 启动链接
			go qconn.Start()
		}
	}()
}

func (s *Server) Stop() {

}

func (s *Server) Run() {
	//启动监听
	s.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		s.Stop()
		SysPrintInfo("exit")
	}
}
