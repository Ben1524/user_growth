package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
	"user_growth/maingin/middleware"
	"user_growth/pb"
)

func init() {
	time.Local = time.UTC // 设置默认时区为UTC
}

func main() {
	// 连接到grpc服务的客户端
	conn, err := grpc.Dial("localhost:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	clientCoin := pb.NewUserCoinClient(conn)
	clientGrade := pb.NewUserGradeClient(conn)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1Group := router.Group("/v1")
	v1Group.Use(middleware.CrossMiddleware)            // 使用跨域中间件
	gUserCoin := v1Group.Group("/UserGrowth.UserCoin") // 子路由组
	gUserCoin.GET("/ListTasks", func(ctx *gin.Context) {
		out, err := clientCoin.ListTasks(ctx, &pb.ListTasksRequest{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    codes.Internal,
				"message": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, out)
		}
	})
	gUserCoin.POST("/UserCoinChange", func(ctx *gin.Context) {
		body := &pb.UserCoinChangeRequest{}
		err := ctx.BindJSON(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    codes.InvalidArgument,
				"message": "Invalid request body",
			})
			return
		}
		if body.Uid <= 0 {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    codes.InvalidArgument,
				"message": "Invalid parameters",
			})
			return
		}
		if out, err := clientCoin.UserCoinChange(ctx, body); err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    codes.Internal,
				"message": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, out)
		}
	})

	gUserGrade := v1Group.Group("/UserGrowth.UserGrade") // 子路由组
	gUserGrade.GET("/ListGrades", func(ctx *gin.Context) {
		out, err := clientGrade.ListGrades(ctx, &pb.ListGradesRequest{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    codes.Internal,
				"message": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, out)
		}
	})
	h2Handler := h2c.NewHandler(router, &http2.Server{})

	server := &http.Server{
		Addr:    ":8080",
		Handler: h2Handler,
	}
	log.Printf("Server is running on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe failed: %v", err)
	} else {
		log.Printf("Server stopped gracefully")
	}
}
