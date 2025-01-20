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

type Server struct {
	config               *config.Config
	tlsConf              *tls.Config
	hookAfterConnStart   func(wiface.IConnection)
	hookBeforeConnStop   func(wiface.IConnection)
	hookAfterStreamStart func(wiface.IStream)
	hookBeforeStreamStop func(wiface.IStream)
	Ctx                  context.Context
	connMgr              wiface.IConnectionMgr
	routerMgr            wiface.IRouterMgr
}

type Options func(opts *Server)

func WithConfig(conf *config.Config) Options {
	return func(opts *Server) {
		opts.config = conf
	}
}

func WithTlsConf(conf *tls.Config) Options {
	return func(opts *Server) {
		opts.tlsConf = conf
	}
}

func WithHookAfterConnStart(hook func(wiface.IConnection)) Options {
	return func(opts *Server) {
		opts.hookAfterConnStart = hook
	}
}

func WithHookBeforeConnStop(hook func(wiface.IConnection)) Options {
	return func(opts *Server) {
		opts.hookBeforeConnStop = hook
	}
}

func WithHookAfterStreamStart(hook func(wiface.IStream)) Options {
	return func(opts *Server) {
		opts.hookAfterStreamStart = hook
	}
}

func WithHookBeforeStreamStop(hook func(wiface.IStream)) Options {
	return func(opts *Server) {
		opts.hookBeforeStreamStop = hook
	}
}

func NewDefaultServer() *Server {
	conf := config.NewDefault()
	opts := []Options{
		WithConfig(conf),
		WithTlsConf(generateTLSConfig(conf.IP, conf.Host)),
		WithHookAfterConnStart(func(conn wiface.IConnection) {
			SysPrintInfo("after connection start call...")
		}),
		WithHookBeforeConnStop(func(conn wiface.IConnection) {
			SysPrintInfo("before connection stop call...")
		}),
		WithHookAfterStreamStart(func(stream wiface.IStream) {
			SysPrintInfo("after stream start call...")
		}),
		WithHookBeforeStreamStop(func(stream wiface.IStream) {
			SysPrintInfo("before stream stop call...")
		}),
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
	SysPrintInfo(fmt.Sprintf("start server listener at ip: %s, port: %d", s.config.IP, s.config.Port))
	SysPrintInfo(fmt.Sprintf("server name: %s, address: %s, version: %s ", s.config.Name, s.getAddress(), s.config.Version))

	go func() {
		addr := s.getAddress()
		listener, err := quic.ListenAddr(addr, generateTLSConfig(s.config.IP, s.config.Host), &quic.Config{
			EnableDatagrams: true,
		})
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = listener.Close()
		}()

		//监听成功
		SysPrintInfo(fmt.Sprintf("listen tcp success, version: %v, addr: %v", s.config.Version, addr))
		connId := NewConnectionID()
		for {
			//获取新的连接
			acceptConn, err := listener.Accept(s.Ctx)
			if err != nil {
				SysPrintError(fmt.Sprintf("accept stream error: %v", err))
				continue
			}
			// 连接限制
			if s.connMgr.Len() >= s.config.MaxConn {
				SysPrintError("too many connections")
				continue
			}
			conn := NewConnection(s.Ctx, s, acceptConn, connId.Get())
			connId.Inc()

			SysPrintInfo(fmt.Sprintf("accept stream success, connId: %v", connId))
			// 启动链接
			go conn.Start()
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
	return fmt.Sprintf("%s:%d", s.config.IP, s.config.Port)
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

func (s *Server) CallHookAfterConnStart(conn wiface.IConnection) {
	if s.hookAfterConnStart != nil {
		s.hookAfterConnStart(conn)
	}
}

func (s *Server) CallHookBeforeConnStop(conn wiface.IConnection) {
	if s.hookBeforeConnStop != nil {
		s.hookBeforeConnStop(conn)
	}
}

func (s *Server) CallHookAfterStreamStart(stream wiface.IStream) {
	if s.hookAfterStreamStart != nil {
		s.hookAfterStreamStart(stream)
	}
}

func (s *Server) CallHookBeforeStreamStop(stream wiface.IStream) {
	if s.hookBeforeStreamStop != nil {
		s.hookAfterStreamStart(stream)
	}
}
