package tracing

const (
	KeyTraceId               = "trace.id"
	keyRpcOperation          = "rpc.operation"
	keyRpcRequestBody        = "rpc.request.body"
	keyRpcResponseStatusCode = "rpc.response.status_code"
	keyRpcResponseBody       = "rpc.response.body"
	keyHttpRequestBody       = "http.request.body"
	keyHttpResponseBody      = "http.response.body"
	keyHttpRequestIp         = "http.request.ip"
	keyHttpRequestHeader     = "http.request.header"
	keySendMsgSize           = "send_msg.size"
	keyRecvMsgSize           = "recv_msg.size"
	keyEnv                   = "env"
	keyOK                    = "OK"
)

const (
	keyRpc  = "rpc"
	keyHttp = "http"
)

const (
	defaultTracerName = "kratos"
)

const serviceHeader = "x-md-service-name"
