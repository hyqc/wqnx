package wnet

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"os"
	"os/signal"
	"wqnx/config"
	"wqnx/wiface"
)

const ()

type Server struct {
	Config    *config.Config
	TlsConf   *tls.Config
	Ctx       context.Context
	connMgr   wiface.IConnectionMgr
	routerMgr wiface.IRouterMgr
}

type Options func(opts *Server)

func WithConfig(conf *config.Config) Options {
	return func(opts *Server) {
		opts.Config = conf
	}
}

func WithTlsConf(conf *tls.Config) Options {
	return func(opts *Server) {
		opts.TlsConf = conf
	}
}

func NewDefaultServer() *Server {
	conf := config.NewDefault()
	opts := []Options{
		WithConfig(conf),
		WithTlsConf(generateTLSConfig(conf.IP, conf.Host)),
	}
	opt := &Server{
		Ctx:       context.Background(),
		connMgr:   NewConnectionMgr(),
		routerMgr: NewRouterMgr(),
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

func (s *Server) Start() {
	SysPrintInfo(fmt.Sprintf("start server listener at ip: %s, port: %d", s.Config.IP, s.Config.Port))
	SysPrintInfo(fmt.Sprintf("server name: %s, address: %s, version: %s ", s.Config.Name, s.getAddress(), s.Config.Version))

	go func() {
		addr := s.getAddress()
		listener, err := quic.ListenAddr(addr, generateTLSConfig(s.Config.IP, s.Config.Host), &quic.Config{
			EnableDatagrams: true,
		})
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = listener.Close()
		}()

		//监听成功
		SysPrintInfo(fmt.Sprintf("listen tcp success, version: %v, addr: %v", s.Config.Version, addr))
		connId := NewConnectionID()
		for {
			//获取新的连接
			conn, err := listener.Accept(s.Ctx)
			if err != nil {
				SysPrintError(fmt.Sprintf("accept stream error: %v", err))
				continue
			}
			// TODO 连接限制
			if s.connMgr.Len() >= s.Config.MaxConn {
				SysPrintError("too many connections")
				continue
			}
			qconn := NewConnection(s.Ctx, s, conn, connId.Get())
			connId.Inc()

			SysPrintInfo(fmt.Sprintf("accept stream success, connId: %v", connId))
			// 启动链接
			go qconn.Start()
		}
	}()
}

func (s *Server) Stop() {
	s.GetConnMgr().Clear()
	SysPrintInfo("server stop")
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

func (s *Server) getAddress() string {
	return fmt.Sprintf("%s:%d", s.Config.IP, s.Config.Port)
}

func (s *Server) GetConnMgr() wiface.IConnectionMgr {
	return s.connMgr
}

func (s *Server) GetCtx() context.Context {
	return s.Ctx
}

func (s *Server) AddRouters(routers ...wiface.IRouter) wiface.IServer {
	if len(routers) == 0 {
		return s
	}
	for _, router := range routers {
		s.routerMgr.Add(router)
	}
	return s
}

func (s *Server) GetRouterMgr() wiface.IRouterMgr {
	return s.routerMgr
}
