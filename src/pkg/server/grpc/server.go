package grpc

import (
	"context"
	"fmt"
	"math"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	listener net.Listener
	srv      *grpc.Server
	opts     options
}

func NewServer(opts ...Option) *Server {

	options := options{
		network: "tcp",
	}
	for _, o := range opts {
		o.apply(&options)
	}

	options.grpcOpts = append(options.grpcOpts, []grpc.ServerOption{
		grpc.MaxSendMsgSize(math.MaxInt32),
		grpc.MaxRecvMsgSize(math.MaxInt32),
	}...)

	if options.maxConnAge > 0 {
		options.grpcOpts = append(options.grpcOpts, grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: options.maxConnAge}))
	}
	s := grpc.NewServer(options.grpcOpts...)

	for _, f := range options.registerServerFuncs {
		f(s)
	}

	reflection.Register(s)

	return &Server{
		srv:  s,
		opts: options,
	}
}

func (s *Server) ListenAndServe() error {
	l, err := net.Listen(s.opts.network, s.opts.listen)
	if err != nil {
		return errors.Wrapf(err, "listen gRPC on address %s", s.opts.listen)
	}
	s.listener = l

	fmt.Println("msg", "listening for serving gRPC", "address", s.opts.listen)
	return errors.Wrap(s.srv.Serve(s.listener), "serve gRPC")
}

func (s *Server) Shutdown(err error) {
	fmt.Println("msg", "internal server is shutting down", "err", err)

	if s.opts.gracePeriod == 0 {
		s.srv.Stop()
		fmt.Println("msg", "internal server is shutdown", "err", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.opts.gracePeriod)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		fmt.Println("msg", "gracefully stopping internal server")
		s.srv.GracefulStop() // Also closes s.listener.
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		fmt.Println("msg", "grace period exceeded enforcing shutdown")
		s.srv.Stop()
		return
	case <-stopped:
		cancel()
	}
	fmt.Println("msg", "internal server is shutdown gracefully", "err", err)
}
