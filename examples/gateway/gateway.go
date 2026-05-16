package main

import (
	"fmt"
	"github.com/liangboceo/dependencyinjection"
	"github.com/liangboceo/yuanboot/abstractions/servicediscovery"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/pkg/configuration"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/middlewares"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// GatewayConfig 网关配置
type GatewayConfig struct {
	Routes []RouteConfig `yaml:"routes"`
}

// GetSection 获取配置节名称
func (c GatewayConfig) GetSection() string {
	return "yuanboot.cloud.gateway"
}

// RouteConfig 路由配置
type RouteConfig struct {
	ID         string   `yaml:"id"`
	URI        string   `yaml:"uri"`        // 如 lb://service-name 或 http://host:port
	Predicates []string `yaml:"predicates"` // 匹配条件，如 Path=/api/**
	Filters    []string `yaml:"filters"`    // 过滤器，如 StripPrefix=1
}

// Route 路由定义
type Route struct {
	Path    []string // 支持多个路径，如 [/assetservice/**, /inventoryservice/**]
	Service string   // 后端服务名
	Filters []string // 过滤器列表
	RawURI  string   // 原始 URI 配置
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
func NewGatewayHandler(config configuration.OptionsSnapshot[GatewayConfig], selector servicediscovery.ISelector, logger xlog.ILogger) *GatewayHandler {
	handler := &GatewayHandler{
		routes:        loadRoutesFromConfig(config),
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

// loadRoutesFromConfig 从配置对象加载路由
func loadRoutesFromConfig(config configuration.OptionsSnapshot[GatewayConfig]) []Route {
	if len(config.CurrentValue().Routes) <= 0 {
		return []Route{}
	}

	routes := make([]Route, 0, len(config.CurrentValue().Routes))
	for _, rc := range config.CurrentValue().Routes {
		route := Route{
			Service: extractServiceName(rc.URI),
			Filters: rc.Filters,
			RawURI:  rc.URI,
		}

		// 解析 predicates
		for _, predicate := range rc.Predicates {
			if strings.HasPrefix(predicate, "Path=") {
				// 解析多个路径，用逗号分隔
				pathsStr := strings.TrimPrefix(predicate, "Path=")
				paths := strings.Split(pathsStr, ",")
				for _, p := range paths {
					p = strings.TrimSpace(p)
					if p != "" {
						route.Path = append(route.Path, p)
					}
				}
			}
		}

		if len(route.Path) > 0 {
			routes = append(routes, route)
		}
	}

	fmt.Printf("Loaded %d routes from config\n", len(routes))
	return routes
}

// extractServiceName 从 URI 中提取服务名
func extractServiceName(uri string) string {
	// 支持 lb://service-name 格式
	if strings.HasPrefix(uri, "lb://") {
		return strings.TrimPrefix(uri, "lb://")
	}
	// http:// 或 https:// 格式，返回空，让路由直接转发
	return ""
}

// RegisterGatewayServices 注册网关服务
func RegisterGatewayServices(sc *dependencyinjection.ServiceCollection) {
	// 注册 GatewayConfig
	configuration.Configure[GatewayConfig](sc)

	// 注册 GatewayHandler
	sc.AddSingleton(func(config configuration.OptionsSnapshot[GatewayConfig], selector servicediscovery.ISelector) *GatewayHandler {
		return NewGatewayHandler(config, selector, middlewares.NewLogger().ALogger)
	})
}

// HandleProxyRequest 处理代理请求
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

	// 应用过滤器
	targetPath := h.applyFilters(path, route.Filters)

	// 如果是 lb:// 格式，需要通过服务发现
	if strings.HasPrefix(route.RawURI, "lb://") {
		instance, err := h.selector.Select(route.Service)
		if err != nil {
			h.logger.Errorf("Failed to select instance for service %s: %v", route.Service, err)
			ctx.JSON(http.StatusServiceUnavailable, context.H{
				"error":   "Service unavailable",
				"service": route.Service,
			})
			return
		}
		h.proxyRequest(ctx, instance, targetPath)
	} else {
		// 直接转发到指定地址
		h.proxyDirectRequest(ctx, route.RawURI, targetPath)
	}
}

// applyFilters 应用过滤器
func (h *GatewayHandler) applyFilters(path string, filters []string) string {
	result := path
	for _, filter := range filters {
		if strings.HasPrefix(filter, "StripPrefix=") {
			parts := strings.Split(filter, "=")
			if len(parts) == 2 {
				n, err := strconv.Atoi(parts[1])
				if err == nil {
					result = h.stripPrefix(result, n)
				}
			}
		}
		// 可以添加更多过滤器：AddRequestHeader, RemoveRequestHeader 等
	}
	return result
}

// stripPrefix 移除路径前缀
func (h *GatewayHandler) stripPrefix(path string, n int) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) <= n {
		return "/"
	}
	return "/" + strings.Join(parts[n:], "/")
}

// HandleRequest 处理请求
func (h *GatewayHandler) HandleRequest(ctx *context.HttpContext) {
	path := ctx.Input.Path()
	h.HandleProxyRequest(ctx, path)
}

// matchRoute 匹配路由
func (h *GatewayHandler) matchRoute(path string) *Route {
	for i := range h.routes {
		for _, routePath := range h.routes[i].Path {
			if h.matchPath(path, routePath) {
				return &h.routes[i]
			}
		}
	}
	return nil
}

// matchPath 匹配路径，支持通配符 **
func (h *GatewayHandler) matchPath(path, pattern string) bool {
	// 移除末尾的 / 保持一致
	path = strings.TrimSuffix(path, "/")
	pattern = strings.TrimSuffix(pattern, "/")

	// 处理通配符 **
	if strings.Contains(pattern, "**") {
		prefix := strings.TrimSuffix(pattern, "/**")
		prefix = strings.TrimSuffix(prefix, "/**")
		if strings.HasPrefix(path, prefix) {
			return true
		}
		return false
	}

	// 处理单层通配符 *
	if strings.Contains(pattern, "*") {
		// 将 * 转换为正则表达式
		parts := strings.Split(pattern, "*")
		if len(parts) == 2 {
			return strings.HasPrefix(path, parts[0]) && strings.HasSuffix(path, parts[1])
		}
	}

	// 精确前缀匹配
	return len(path) >= len(pattern) && path[:len(pattern)] == pattern
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
		if route.Service == "" {
			continue
		}
		instance, err := h.selector.Select(route.Service)
		if err == nil && instance != nil {
			h.instanceCache[route.Service] = []servicediscovery.ServiceInstance{instance}
			h.logger.Debugf("Refreshed cache for %s", route.Service)
		}
	}
}

// proxyRequest 代理请求（通过服务发现）
func (h *GatewayHandler) proxyRequest(ctx *context.HttpContext, instance servicediscovery.ServiceInstance, targetPath string) {
	targetURL := fmt.Sprintf("http://%s:%d%s", instance.GetHost(), instance.GetPort(), targetPath)
	h.proxyDirectRequest(ctx, targetURL, targetPath)
}

// proxyDirectRequest 直接代理请求
func (h *GatewayHandler) proxyDirectRequest(ctx *context.HttpContext, targetBase, targetPath string) {
	targetURL := targetBase + targetPath

	h.logger.Debugf("Proxying to: %s", targetURL)

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
	services := make([]string, 0)
	seen := make(map[string]bool)
	for _, route := range h.routes {
		if route.Service != "" && !seen[route.Service] {
			services = append(services, route.Service)
			seen[route.Service] = true
		}
	}
	return services
}

// GetRoutes 获取路由配置
func (h *GatewayHandler) GetRoutes() []map[string]interface{} {
	routes := make([]map[string]interface{}, len(h.routes))
	for i, route := range h.routes {
		routes[i] = map[string]interface{}{
			"id":      route.RawURI,
			"uri":     route.RawURI,
			"paths":   route.Path,
			"service": route.Service,
			"filters": route.Filters,
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
