# kratos-middleware

### 介绍
kratos框架中间件

### tracing
Tracing 中间件使用 OpenTelemetry 实现了链路追踪。
> #### 引用
> ```go
> import "kratos-middleware/tracing"
> ```
> #### 设置全局跟踪提供程序示例
> ```go
> 	err := tracing.SetTracerProvider(traceUrl, name, "dev", 1.0)
> 	if err != nil {
> 		log.Error(err)
> 	}
> ```
> #### 服务端--http接口增加中间件(http.Middleware()方法中)
> ```go
> tracing.Server("http"),
> ```
> #### 服务端--rpc接口增加中间件(grpc.Middleware()方法中)
> ```go
> tracing.Server("rpc"),
> ```
> #### 客户端--增加参数
> ```go
>   import grpcx "google.golang.org/grpc"
>	
>   grpc.WithMiddleware(
>       recovery.Recovery(),
>       tracing.Client()),
>   grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
> ```
> #### 特别注意
> - logging.Server(logger)要写在tracing中间件引用的后面，否则会造成trace_id和span_id为空的问题

