package main

import (
	"fmt"
	"github.com/liangboceo/dependencyinjection"
	"github.com/liangboceo/yuanboot/abstractions/servicediscovery"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/web/context"
	"io"
	"net/http"
	"sync"
	"time"
)

// Route 路由配置
type Route struct {
	Path    string // 路由路径，如 /api/, /user/
	Service string // 后端服务名
	Rewrite string // 路径重写规则
}

// GatewayHandler 网关处理器
type GatewayHandler struct {
	routes        []Route
	selector      servicediscovery.ISelector
	logger        xlog.ILogger
	stopChan      chan struct{}
	instanceCache map[string][]servicediscovery.ServiceInstance
	cacheLock     sync.RWMutex
	httpClient    *http.Client
}

// NewGatewayHandler 创建网关处理器
func NewGatewayHandler(selector servicediscovery.ISelector, logger xlog.ILogger) *GatewayHandler {
	handler := &GatewayHandler{
		routes:        getDefaultRoutes(),
		selector:      selector,
		logger:        logger,
		stopChan:      make(chan struct{}),
		instanceCache: make(map[string][]servicediscovery.ServiceInstance),
		httpClient:    &http.Client{Timeout: 30 * time.Second},
	}

	// 启动实例缓存刷新
	go handler.startCacheRefresh()

	return handler
}

// getDefaultRoutes 获取默认路由配置
func getDefaultRoutes() []Route {
	return []Route{
		{Path: "/api/", Service: "yuanboot-api", Rewrite: "/"},
		{Path: "/user/", Service: "user-service", Rewrite: "/"},
		{Path: "/order/", Service: "order-service", Rewrite: "/"},
		{Path: "/product/", Service: "product-service", Rewrite: "/"},
	}
}

// RegisterGatewayServices 注册网关服务
func RegisterGatewayServices(sc *dependencyinjection.ServiceCollection) {
	sc.AddSingletonByImplements(
		func(selector servicediscovery.ISelector, logger xlog.ILogger) *GatewayHandler {
			return NewGatewayHandler(selector, logger)
		},
		new(*GatewayHandler),
	)
}

// HandleProxyRequest 处理代理请求（带完整路径）
func (h *GatewayHandler) HandleProxyRequest(ctx *context.HttpContext, path string) {
	h.logger.Debugf("Gateway received request: %s %s", ctx.Input.Method(), path)

	// 查找匹配的路由
	route := h.matchRoute(path)
	if route == nil {
		ctx.JSON(http.StatusNotFound, context.H{
			"error": "Route not found",
			"path":  path,
		})
		return
	}

	// 使用 selector 选择服务实例
	instance, err := h.selector.Select(route.Service)
	if err != nil {
		h.logger.Errorf("Failed to select instance for service %s: %v", route.Service, err)
		ctx.JSON(http.StatusServiceUnavailable, context.H{
			"error":   "Service unavailable",
			"service": route.Service,
		})
		return
	}

	// 重写路径
	targetPath := h.rewritePath(path, route.Rewrite)

	// 转发请求
	h.proxyRequest(ctx, instance, targetPath)
}

// HandleRequest 处理请求
func (h *GatewayHandler) HandleRequest(ctx *context.HttpContext) {
	path := ctx.Input.Path()
	h.HandleProxyRequest(ctx, path)
}

// matchRoute 匹配路由
func (h *GatewayHandler) matchRoute(path string) *Route {
	for i := range h.routes {
		if len(path) >= len(h.routes[i].Path) && path[:len(h.routes[i].Path)] == h.routes[i].Path {
			return &h.routes[i]
		}
	}
	return nil
}

// rewritePath 重写路径
func (h *GatewayHandler) rewritePath(originalPath, rewrite string) string {
	if rewrite == "/" {
		return originalPath
	}
	// 简单实现：替换前缀
	for _, route := range h.routes {
		if len(originalPath) >= len(route.Path) && originalPath[:len(route.Path)] == route.Path {
			return rewrite + originalPath[len(route.Path):]
		}
	}
	return originalPath
}

// startCacheRefresh 启动缓存刷新
func (h *GatewayHandler) startCacheRefresh() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.refreshCache()
		case <-h.stopChan:
			return
		}
	}
}

// refreshCache 刷新缓存
func (h *GatewayHandler) refreshCache() {
	h.cacheLock.Lock()
	defer h.cacheLock.Unlock()

	for _, route := range h.routes {
		// 使用 selector 获取服务
		instance, err := h.selector.Select(route.Service)
		if err == nil && instance != nil {
			h.instanceCache[route.Service] = []servicediscovery.ServiceInstance{instance}
			h.logger.Debugf("Refreshed cache for %s", route.Service)
		}
	}
}

// proxyRequest 代理请求
func (h *GatewayHandler) proxyRequest(ctx *context.HttpContext, instance servicediscovery.ServiceInstance, targetPath string) {
	// 构建目标 URL
	targetURL := fmt.Sprintf("http://%s:%d%s", instance.GetHost(), instance.GetPort(), targetPath)

	h.logger.Debugf("Proxying to: %s", targetURL)

	// 创建代理请求
	req, err := http.NewRequest(ctx.Input.Method(), targetURL, nil)
	if err != nil {
		h.logger.Errorf("Failed to create request: %v", err)
		ctx.JSON(http.StatusBadGateway, context.H{
			"error": "Failed to create request",
		})
		return
	}

	// 复制请求头
	for key := range ctx.Input.Request.Header {
		req.Header.Set(key, ctx.Input.Header(key))
	}

	// 添加转发头
	req.Header.Set("X-Gateway", "yuanboot")
	req.Header.Set("X-Forwarded-For", ctx.Input.RemoteIP())
	req.Header.Set("X-Real-IP", ctx.Input.RealIP())

	// 执行请求
	resp, err := h.httpClient.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to proxy request: %v", err)
		ctx.JSON(http.StatusBadGateway, context.H{
			"error": "Failed to reach backend service",
		})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			ctx.Output.Header(key, value)
		}
	}

	// 设置状态码
	ctx.Output.SetStatusCode(resp.StatusCode)

	// 复制响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Errorf("Failed to read response body: %v", err)
		ctx.JSON(http.StatusBadGateway, context.H{
			"error": "Failed to read response",
		})
		return
	}

	ctx.Output.Write(body)
}

// Stop 停止网关
func (h *GatewayHandler) Stop() {
	close(h.stopChan)
}

// GetRegisteredServices 获取已注册的服务列表
func (h *GatewayHandler) GetRegisteredServices() []string {
	services := make([]string, len(h.routes))
	for i, route := range h.routes {
		services[i] = route.Service
	}
	return services
}

// GetRoutes 获取路由配置
func (h *GatewayHandler) GetRoutes() []map[string]string {
	routes := make([]map[string]string, len(h.routes))
	for i, route := range h.routes {
		routes[i] = map[string]string{
			"path":    route.Path,
			"service": route.Service,
			"rewrite": route.Rewrite,
		}
	}
	return routes
}

// GetServiceInstances 获取服务实例
func (h *GatewayHandler) GetServiceInstances(serviceName string) []servicediscovery.ServiceInstance {
	h.cacheLock.RLock()
	defer h.cacheLock.RUnlock()
	return h.instanceCache[serviceName]
}

// ForceRefresh 强制刷新缓存
func (h *GatewayHandler) ForceRefresh() {
	h.refreshCache()
}
