package service

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Option func(*App)

type ShutdownCallback func(ctx context.Context)

func WithShutdownCallback(cbs ...ShutdownCallback) Option {
	panic("implement me")
}

type App struct {
	servers []*Server

	shutdownTimeout time.Duration

	waitTime time.Duration

	cbTimeout time.Duration

	cbs []ShutdownCallback
}

func NewApp(servers []*Server, opts ...Option) *App {
	panic("implement me")
}

func (app *App) StartAndServe() {
	for _, s := range app.servers {
		srv := s
		go func() {
			if err := srv.Start(); err != nil {
				log.Printf("服务器%s已关闭", srv.name)
			} else {
				log.Printf("服务器%s异常退出", srv.name)
			}

		}()
	}
}

func (app *App) shutdown() {
	log.Println("开始关闭应用，停止接受新请求")

	log.Println("等待正在执行请求完结")

	log.Println("开始关闭服务器")

	log.Println("开始执行自定义回调")

	log.Println("开始释放资源")
	app.close()
}

func (app *App) close() {
	time.Sleep(time.Second)
	log.Println("应用关闭")
}

type Server struct {
	srv  *http.Server
	name string
	mux  *serverMux
}

type serverMux struct {
	reject bool
	*http.ServeMux
}

func (s *serverMux) ServerHttp(w http.ResponseWriter, r *http.Request) {
	if s.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("服务已关闭"))
		return
	}
	s.ServeMux.ServeHTTP(w, r)
}

func NewServer(name string, addr string) *Server {
	mux := &serverMux{ServeMux: http.NewServeMux()}
	return &Server{
		name: name,
		mux:  mux,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) rejectReq() {
	s.mux.reject = true
}

func (s *Server) stop() error {
	log.Printf("服务器%s关闭中", s.name)
	return s.srv.Shutdown(context.Background())
}
