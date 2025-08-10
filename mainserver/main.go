package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"
	_ "user_growth/conf"
	"user_growth/pb"
	"user_growth/ugserver"
)

var AllowOrigin = map[string]bool{
	"http://a.site.com": true, // 允许的跨域源
	"http://b.site.com": true,
	"http://web.com":    true,
	"http://12.0.0.1":   true,
	"http://localhost":  true, // 允许的跨域源
}

// 定义跨域中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if AllowOrigin[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// 处理 OPTIONS 预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func init() {
	// default UTC time location
	time.Local = time.UTC
}

func main() {

	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//creds, err := credentials.NewServerTLSFromFile("cert/server.pem", "cert/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile error=%v", err)
	}
	opts := []grpc.ServerOption{
		grpc.WriteBufferSize(1024 * 1024 * 1), // 默认32KB
		grpc.ReadBufferSize(1024 * 1024 * 1),  // 默认32KB
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     10 * time.Minute, // 没有消息的最长时间
			MaxConnectionAge:      1 * time.Hour,    // 连接最长时间
			MaxConnectionAgeGrace: 10 * time.Minute, // 最长时间后延迟关闭
			Time:                  2 * time.Minute,  // ping间隔
			Timeout:               3 * time.Second,  // ping超时
		}),
		grpc.MaxConcurrentStreams(1000),
		grpc.ConnectionTimeout(time.Second * 1), // 连接超时
		//grpc.Creds(creds),
	}
	s := grpc.NewServer(opts...)
	// 注册服务
	pb.RegisterUserCoinServer(s, &ugserver.UgCoinServer{})
	pb.RegisterUserGradeServer(s, &ugserver.UgGradeServer{})
	reflection.Register(s)

	//=========================================================
	// 增加grpc-gateway的支持
	serveMuxOpt := []runtime.ServeMuxOption{
		runtime.WithOutgoingHeaderMatcher(func(s string) (string, bool) {
			return s, true
		}),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			origin := request.Header.Get("Origin")
			if AllowOrigin[origin] {
				md := metadata.New(map[string]string{
					"Access-Control-Allow-Origin":      origin,
					"Access-Control-Allow-Methods":     "GET,POST,PUT,DELETE,OPTION",
					"Access-Control-Allow-Headers":     "*",
					"Access-Control-Allow-Credentials": "true",
				})
				grpc.SetHeader(ctx, md)
			}
			return nil
		}),
	}
	mux := runtime.NewServeMux(serveMuxOpt...) // 创建grpc-gateway的ServeMux，负责将HTTP请求转发到gRPC服务
	ctx := context.Background()
	if err := pb.RegisterUserCoinHandlerServer(ctx, mux, &ugserver.UgCoinServer{}); err != nil {
		log.Printf("Faile to RegisterUserCoinHandlerServer error=%v", err)
	}
	if err := pb.RegisterUserGradeHandlerServer(ctx, mux, &ugserver.UgGradeServer{}); err != nil {
		log.Printf("Faile to RegisterUserGradeHandlerServer error=%v", err)
	}

	//httpMux := http.NewServeMux()
	//httpMux.Handle("/v2/UserGrowth", mux) // 注册grpc-gateway的路由，将/v1/UserGrowth的请求转发到grpc服务
	// 配置http服务
	server := &http.Server{
		Addr: ":8081",
		Handler: h2c.NewHandler(
			corsMiddleware(mux), // 包裹跨域中间件
			&http2.Server{}),
	}

	// 启动http服务
	log.Printf("server.ListenAdnServe(%s)", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("ListenAndServe error=%v", err)
	}
	//=========================================================

	// 启动服务
	log.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
