package client

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	motan "github.com/weibocom/motan-go"
	"github.com/weibocom/motan-go/cluster"
	"github.com/weibocom/motan-go/core"
	"github.com/weibocom/motan-go/endpoint"
	"github.com/weibocom/motan-go/protocol"
)

// MotanTiming Motan请求时间记录结构
// 注意：由于Motan通过Mesh代理调用，BaseCall是同步阻塞调用，无法精确追踪各个网络阶段
// 因此只记录总耗时
type MotanTiming struct {
	TotalTime time.Duration // 总耗时
}

const meshDirectRegistryKey = "mesh-registry"

var (
	customMeshClient  *CustomMeshClient
	meshClientOnce    sync.Once
	meshClientInitErr error
	errMeshClientErr  = errors.New("mesh client error")
	// 保存配置用于超时范围验证
	motanConfig *MotanConfig
)

// MotanConfig Motan配置
type MotanConfig struct {
	MeshAddress       string        // Mesh代理地址，默认 "127.0.0.1:9981"
	Application       string        // 调用服务的业务方
	Serialization     string        // 序列化方式，默认 "simple"
	RequestTimeout    time.Duration // 请求超时时间，默认 1秒
	MinRequestTimeout time.Duration // 最小请求超时时间，默认为 RequestTimeout/2
	MaxRequestTimeout time.Duration // 最大请求超时时间，默认为 RequestTimeout*2
}

// CustomMeshClient 自定义MeshClient，支持设置minRequestTimeout和maxRequestTimeout
type CustomMeshClient struct {
	requestTimeout    time.Duration
	minRequestTimeout time.Duration
	maxRequestTimeout time.Duration
	application       string
	address           string
	serialization     string
	cluster           *cluster.MotanCluster
}

// NewCustomMeshClient 创建自定义MeshClient
func NewCustomMeshClient() *CustomMeshClient {
	return &CustomMeshClient{}
}

// SetAddress 设置Mesh代理地址
func (c *CustomMeshClient) SetAddress(address string) {
	c.address = address
}

// SetRequestTimeout 设置请求超时时间
func (c *CustomMeshClient) SetRequestTimeout(requestTimeout time.Duration) {
	c.requestTimeout = requestTimeout
}

// SetMinRequestTimeout 设置最小请求超时时间
func (c *CustomMeshClient) SetMinRequestTimeout(minRequestTimeout time.Duration) {
	c.minRequestTimeout = minRequestTimeout
}

// SetMaxRequestTimeout 设置最大请求超时时间
func (c *CustomMeshClient) SetMaxRequestTimeout(maxRequestTimeout time.Duration) {
	c.maxRequestTimeout = maxRequestTimeout
}

// SetSerialization 设置序列化方式
func (c *CustomMeshClient) SetSerialization(serialization string) {
	c.serialization = serialization
}

// SetApplication 设置应用名称
func (c *CustomMeshClient) SetApplication(application string) {
	c.application = application
}

// Initialize 初始化MeshClient
func (c *CustomMeshClient) Initialize() error {
	// 双重检查：如果已经初始化，直接返回
	if c.cluster != nil {
		return nil
	}

	if c.requestTimeout == 0 {
		c.requestTimeout = 5 * time.Second
	}
	if c.address == "" {
		c.address = "127.0.0.1:9981"
	}
	if c.serialization == "" {
		c.serialization = "simple"
	}
	// 设置默认的超时范围
	if c.minRequestTimeout == 0 {
		c.minRequestTimeout = c.requestTimeout / 2
	}
	if c.maxRequestTimeout == 0 {
		c.maxRequestTimeout = c.requestTimeout * 2
	}

	clusterURL := &core.URL{}
	clusterURL.Protocol = endpoint.Motan2
	clusterURL.PutParam(core.TimeOutKey, strconv.Itoa(int(c.requestTimeout/time.Millisecond)))
	clusterURL.PutParam(core.ApplicationKey, c.application)
	clusterURL.PutParam(core.ErrorCountThresholdKey, "0")
	clusterURL.PutParam(core.RegistryKey, meshDirectRegistryKey)
	clusterURL.PutParam(core.ConnectRetryIntervalKey, "5000")
	clusterURL.PutParam(core.SerializationKey, c.serialization)
	clusterURL.PutParam(core.AsyncInitConnection, "false")
	// 设置最小和最大请求超时时间，使 interface.yaml 中配置的超时能够精确生效
	clusterURL.PutParam(core.MinTimeOutKey, strconv.Itoa(int(c.minRequestTimeout/time.Millisecond)))
	clusterURL.PutParam(core.MaxTimeOutKey, strconv.Itoa(int(c.maxRequestTimeout/time.Millisecond)))

	meshRegistryURL := &core.URL{}
	meshRegistryURL.Protocol = "direct"
	meshRegistryURL.PutParam(core.AddressKey, c.address)
	context := &core.Context{}
	context.RegistryURLs = make(map[string]*core.URL)
	context.RegistryURLs[meshDirectRegistryKey] = meshRegistryURL

	// 创建 cluster 并检查是否成功
	cluster := cluster.NewCluster(context, motan.GetDefaultExtFactory(), clusterURL, false)
	if cluster == nil {
		return errors.New("failed to create motan cluster: cluster.NewCluster returned nil")
	}

	c.cluster = cluster
	return nil
}

// Destroy 销毁MeshClient
func (c *CustomMeshClient) Destroy() {
	if c.cluster != nil {
		c.cluster.Destroy()
	}
}

// BuildRequest 构建请求
func (c *CustomMeshClient) BuildRequest(service string, method string, args []interface{}) core.Request {
	return c.BuildRequestWithGroup(service, method, args, "")
}

// BuildRequestWithGroup 构建带Group的请求
func (c *CustomMeshClient) BuildRequestWithGroup(service string, method string, args []interface{}, group string) core.Request {
	request := &core.MotanRequest{Method: method, ServiceName: service, Arguments: args, Attachment: core.NewStringMap(core.DefaultAttachmentSize)}
	request.RequestID = endpoint.GenerateRequestID()
	request.SetAttachment(protocol.MSource, c.application)
	request.SetAttachment(protocol.MGroup, group)
	request.SetAttachment(protocol.MPath, request.GetServiceName())
	return request
}

// BaseCall 执行请求
func (c *CustomMeshClient) BaseCall(request core.Request, reply interface{}) core.Response {
	rc := request.GetRPCContext(true)
	rc.Reply = reply
	return c.cluster.Call(request)
}

// InitMotanMeshClient 初始化MeshClient
// MeshClient主要用来调用Mesh进程，服务发现/治理等能力都是由Mesh进程提供的
func InitMotanMeshClient(config *MotanConfig) error {
	meshClientOnce.Do(func() {
		// 验证配置
		if config == nil {
			meshClientInitErr = errors.New("motan config is nil")
			return
		}
		if config.MeshAddress == "" {
			meshClientInitErr = errors.New("motan mesh_address is required")
			return
		}
		if config.Application == "" {
			meshClientInitErr = errors.New("motan application is required")
			return
		}
		if config.Serialization == "" {
			meshClientInitErr = errors.New("motan serialization is required")
			return
		}
		if config.RequestTimeout == 0 {
			meshClientInitErr = errors.New("motan request_timeout is required")
			return
		}

		// 设置默认的超时范围
		if config.MinRequestTimeout == 0 {
			config.MinRequestTimeout = config.RequestTimeout / 2
		}
		if config.MaxRequestTimeout == 0 {
			config.MaxRequestTimeout = config.RequestTimeout * 2
		}

		// 保存配置
		motanConfig = config

		// 创建并初始化自定义MeshClient（支持设置minRequestTimeout和maxRequestTimeout）
		customMeshClient = NewCustomMeshClient()
		customMeshClient.SetAddress(config.MeshAddress)
		customMeshClient.SetApplication(config.Application)
		customMeshClient.SetSerialization(config.Serialization)
		customMeshClient.SetRequestTimeout(config.RequestTimeout)
		customMeshClient.SetMinRequestTimeout(config.MinRequestTimeout)
		customMeshClient.SetMaxRequestTimeout(config.MaxRequestTimeout)

		// 检查 Initialize 的返回值
		if initErr := customMeshClient.Initialize(); initErr != nil {
			meshClientInitErr = fmt.Errorf("failed to initialize motan mesh client: %w", initErr)
			return
		}

		if customMeshClient == nil {
			meshClientInitErr = errMeshClientErr
		}
	})
	return meshClientInitErr
}

// GetCustomMeshClient 获取自定义MeshClient实例
func GetCustomMeshClient() *CustomMeshClient {
	return customMeshClient
}

// MotanCallWithTiming 调用Motan服务（通过MeshClient）并返回时间记录
// service: 要调用的RPC服务名称（如 "com.weibo.api.FeedCoreStatusHttpService"）
// method: 要调用rpc的方法，在v4 RPC服务中，即http请求时的url路径（如 "/2/statuses/show_batch.json"）
// args: 参数（map[string]string 类型）
// headers: 请求头（通过SetAttachment设置，如Authorization）
// timeout: 超时时间（可选，如果为0则使用默认超时）
func MotanCallWithTiming(service, method string, args interface{}, headers map[string]string, timeout time.Duration) (resp []byte, reqID uint64, timing *MotanTiming, err error) {
	if customMeshClient == nil {
		return nil, 0, nil, errMeshClientErr
	}

	// 检查 cluster 是否已正确初始化
	if customMeshClient.cluster == nil {
		return nil, 0, nil, errors.New("motan cluster is not initialized, please call InitMotanMeshClient first")
	}

	timing = &MotanTiming{}

	// 记录开始时间
	startTime := time.Now()

	// 构造请求Request，customMeshClient会自行填充service等必要信息
	req := customMeshClient.BuildRequest(service, method, []interface{}{args})

	// 通过motan的attachment传递请求头（如认证信息）
	for k, v := range headers {
		req.SetAttachment(k, v)
	}

	// 设置超时时间（如果指定）
	// 由于已在初始化时设置了 minRequestTimeout 和 maxRequestTimeout，
	// interface.yaml 中配置的超时时间只要在该范围内就能精确生效
	if timeout > 0 {
		// 将超时时间转换为毫秒并通过 M_tmo attachment 传递
		timeoutMs := int64(timeout / time.Millisecond)
		req.SetAttachment("M_tmo", fmt.Sprintf("%d", timeoutMs))
	}

	// 使用BaseCall方法调用服务
	var reply []byte
	response := customMeshClient.BaseCall(req, &reply)

	// 检查是否有异常
	if response.GetException() != nil {
		timing.TotalTime = time.Since(startTime)
		return nil, req.GetRequestID(), timing, fmt.Errorf("motan call error: %s", response.GetException().ErrMsg)
	}

	// 计算总耗时
	timing.TotalTime = time.Since(startTime)

	return reply, req.GetRequestID(), timing, nil
}

// MotanCall 调用Motan服务（通过MeshClient）
// service: 要调用的RPC服务名称（如 "com.weibo.api.FeedCoreStatusHttpService"）
// method: 要调用rpc的方法，在v4 RPC服务中，即http请求时的url路径（如 "/2/statuses/show_batch.json"）
// args: 参数（map[string]string 类型）
// headers: 请求头（通过SetAttachment设置，如Authorization）
// timeout: 超时时间（可选，如果为0则使用默认超时）
func MotanCall(service, method string, args interface{}, headers map[string]string, timeout time.Duration) (resp []byte, reqID uint64, err error) {
	resp, reqID, _, err = MotanCallWithTiming(service, method, args, headers, timeout)
	return resp, reqID, err
}

// IsMotanInitialized 检查MeshClient是否已初始化
func IsMotanInitialized() bool {
	return customMeshClient != nil && customMeshClient.cluster != nil
}
