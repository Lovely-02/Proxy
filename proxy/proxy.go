package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"proxy/config"
	"proxy/utils"
)

// ProxyServer 代理服务器
type ProxyServer struct {
	config       *config.Config
	reverseProxy *httputil.ReverseProxy
}

// transportWrapper 包装Transport，用于修改响应
type transportWrapper struct {
	original   http.RoundTripper
	proxyHost  string
	targetHost string
}

// RoundTrip 实现http.RoundTripper接口
func (t *transportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	// 发送请求
	var resp *http.Response
	var err error

	if t.original != nil {
		resp, err = t.original.RoundTrip(req)
	} else {
		// 如果original为nil，使用默认的http.Transport
		resp, err = http.DefaultTransport.RoundTrip(req)
	}

	if err != nil {
		return nil, err
	}

	// 修改重定向响应的Location头
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		if location := resp.Header.Get("Location"); location != "" {
			// 解析Location URL
			locURL, err := url.Parse(location)
			if err == nil && locURL.Host == t.targetHost {
				locURL.Host = t.proxyHost
				resp.Header.Set("Location", locURL.String())
			}
		}
	}

	return resp, nil
}

// NewProxyServer 创建新的代理服务器实例
func NewProxyServer(cfg *config.Config) *ProxyServer {
	// 解析目标URL，使用用户传入的完整URL
	targetURL, err := url.Parse(cfg.Target)
	if err != nil {
		utils.Error(fmt.Sprintf("解析目标URL失败: %v", err))
		panic(err)
	}

	// 创建反向代理
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 自定义请求处理
	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		// 修改请求头，添加X-Forwarded-For
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		// 设置Host头为目标服务器主机名
		req.Host = targetURL.Host
	}

	// 自定义响应处理，修改重定向Location头
	originalTransport := reverseProxy.Transport
	reverseProxy.Transport = &transportWrapper{
		original:   originalTransport,
		proxyHost:  cfg.Listen,
		targetHost: targetURL.Host,
	}

	// 自定义错误处理
	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		utils.Error(fmt.Sprintf("代理错误: %v", err))
		http.Error(w, "代理错误", http.StatusInternalServerError)
	}

	return &ProxyServer{
		config:       cfg,
		reverseProxy: reverseProxy,
	}
}

// Start 启动代理服务器
func (s *ProxyServer) Start() error {
	// 注册处理函数
	http.HandleFunc("/", s.handleRequest)

	// 启动服务器
	utils.Info(fmt.Sprintf("代理服务器已启动在 %s", s.config.Listen))
	return http.ListenAndServe(s.config.Listen, nil)
}

// handleRequest 处理HTTP请求
func (s *ProxyServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	// 记录请求信息
	utils.Info(fmt.Sprintf("请求: %s %s", r.Method, r.URL.Path))

	// 处理WebSocket连接
	if r.Header.Get("Upgrade") == "websocket" {
		s.handleWebSocket(w, r)
		return
	}

	// 转发HTTP请求
	s.reverseProxy.ServeHTTP(w, r)
}

// handleWebSocket 处理WebSocket连接
func (s *ProxyServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	utils.Info("WebSocket连接请求")
	s.reverseProxy.ServeHTTP(w, r)
}
