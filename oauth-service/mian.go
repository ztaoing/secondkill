/**
* @Author:zhoutao
* @Date:2020/7/5 上午7:31
 */

package main

import (
	"context"
	"flag"
	"fmt"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"os"
	"os/signal"
	localConfig "secondkill/oauth-service/config"
	"secondkill/oauth-service/endpoint"
	"secondkill/oauth-service/plugins"
	"secondkill/oauth-service/service"
	"secondkill/oauth-service/transports"
	"secondkill/pb"
	"secondkill/pkg/bootstrap"
	"secondkill/pkg/config"
	register "secondkill/pkg/discover"
	"secondkill/pkg/mysql"
	"syscall"
	"time"
)

func main() {
	var (
		servicePort = flag.String("service.port", bootstrap.HttpConfig.Port, "service port")
		grpcAddr    = flag.String("grpc", bootstrap.RpcConfig.Port, "grpc listen address")
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	//限流
	rateBucket := rate.NewLimiter(rate.Every(time.Second*1), 100)

	var tokenService service.TokenService
	var tokenGranter service.TokenGranter
	var tokenEnhancer service.TokenEnhancer
	var tokenStore service.TokenStore
	var userDetailsService service.UserDetailsService
	var clientDetailsService service.ClientDetailsService
	var svc service.Service

	//token with jwt
	tokenEnhancer = service.NewJwtTokenEnhancer("secret")
	tokenStore = service.NewJwtTokenStore(tokenEnhancer.(*service.JwtTokenEnhancer))

	tokenService = service.NewTokenService(tokenStore, tokenEnhancer)

	userDetailsService = service.NewRemoteUserDetailsService()

	clientDetailsService = service.NewMysqlClientDetailsService()

	svc = service.NewCommonService()

	tokenGranter = service.NewComposeTokenGranter(map[string]service.TokenGranter{
		"password":      service.NewUsernamePasswordTokenGranter("password", userDetailsService, tokenService),
		"refresh_token": service.NewRefreshTokenGrant("refresh_token", userDetailsService, tokenService),
	})

	//tokenEndpoint
	tokenEndpoint := endpoint.MakeTokenEndpoint(tokenGranter, clientDetailsService)
	//在进入Endpoint之前统一验证context中的OAuthDetails是否存在
	tokenEndpoint = endpoint.MakeClientAuthorizationMiddleware(localConfig.Logger)(tokenEndpoint)
	//在进入endpoint之前，限流（此处使用的time/rate）限流组件
	tokenEndpoint = plugins.NewTokenBucketLimitterWithBuildIn(rateBucket)(tokenEndpoint)
	//添加 名称为token-endpoint的链路追踪器
	tokenEndpoint = kitzipkin.TraceEndpoint(localConfig.ZipKinTracer, "token-endpoint")(tokenEndpoint)

	//checkEndpoint
	checkTokenEndpoint := endpoint.MakeCheckTokenEndpoint(tokenService)
	//校验客户端信息和用户信息是否存在
	checkTokenEndpoint = endpoint.MakeClientAuthorizationMiddleware(localConfig.Logger)(checkTokenEndpoint)
	//进入checkTokenEndpoint之前进行限流
	checkTokenEndpoint = plugins.NewTokenBucketLimitterWithBuildIn(rateBucket)(checkTokenEndpoint)
	//添加 名称为check-endpoint的链路追踪器
	checkTokenEndpoint = kitzipkin.TraceEndpoint(localConfig.ZipKinTracer, "check-endpoint")(checkTokenEndpoint)

	//grpcCheckToken
	grpcCheckTokenEndpoint := endpoint.MakeCheckTokenEndpoint(tokenService)
	//在进入endpoint之前，限流（此处使用的time/rate）限流组件
	grpcCheckTokenEndpoint = plugins.NewTokenBucketLimitterWithBuildIn(rateBucket)(grpcCheckTokenEndpoint)
	//添加 名称为grpc-check-endpoint的链路追踪器
	grpcCheckTokenEndpoint = kitzipkin.TraceEndpoint(localConfig.ZipKinTracer, "grpc-check-endpoint")(grpcCheckTokenEndpoint)

	//healthEndpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	healthEndpoint = kitzipkin.TraceEndpoint(localConfig.ZipKinTracer, "health-endpoint")(healthEndpoint)

	endpts := endpoint.OAuthEndpoints{
		TokenEndpoint:          tokenEndpoint,
		CheckTokenEndpoint:     checkTokenEndpoint,
		GRPCCheckTokenEndpoint: grpcCheckTokenEndpoint,
		HealthCheckEndpoint:    healthEndpoint,
	}

	//transport with http
	//创建http.Handler
	r := transports.MakeHttpHandler(ctx, endpts, tokenService, clientDetailsService, localConfig.ZipKinTracer, localConfig.Logger)

	//http-server
	go func() {
		fmt.Println("http server start at port:" + *servicePort)

		//初始化MySQL
		mysql.InitMysql(config.MysqlConfig.Host, config.MysqlConfig.Port, config.MysqlConfig.User, config.MysqlConfig.Pwd, config.MysqlConfig.Db)

		//将服务注册到注册中心
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()

	//grpc-server
	go func() {
		fmt.Println("grpc server start at port:" + *grpcAddr)
		listenner, err := net.Listen("tcp", ":"+*grpcAddr)

		if err != nil {
			errChan <- err
			return
		}

		serverTracer := kitzipkin.GRPCServerTrace(localConfig.ZipKinTracer, kitzipkin.Name("grpc-transport"))
		grpcTracer := localConfig.ZipKinTracer
		md := metadata.MD{}
		//创建root span
		rootSpan := grpcTracer.StartSpan("grpc-root-span")
		//InjectGRPC will inject a span.Context into gRPC metadata.
		b3.InjectGRPC(&md)(rootSpan.Context())
		//NewIncomingContext creates a new context with incoming md attached.
		ctx := metadata.NewIncomingContext(context.Background(), md)
		handler := transports.NewGRPCServer(ctx, endpts, serverTracer)

		//todo ?
		gRPCserver := grpc.NewServer()
		pb.RegisterOAuthServiceServer(gRPCserver, handler)
		errChan <- gRPCserver.Serve(listenner)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	//等待错误
	error := <-errChan
	//取消注册
	register.Deregister()
	fmt.Println(error)

}
