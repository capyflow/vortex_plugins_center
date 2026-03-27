package gateway

import (
	"net/http"

	"plugin-platform/internal/center"
	"plugin-platform/pkg/models"

	"github.com/capyflow/vortexv3/pkg"
	"github.com/capyflow/vortexv3/server/http"
)

// Gateway API 网关
type Gateway struct {
	center *center.PluginCenter
	engine *http.VortexHttpRouterGroup
}

// New 创建网关
func New(pc *center.PluginCenter) *Gateway {
	return &Gateway{
		center: pc,
	}
}

// Start 启动网关
func (g *Gateway) Start(addr string) error {
	// 创建 Vortex HTTP 引擎
	root := http.NewRootGroup("/api/v1")

	// 注册路由
	g.registerRoutes(root)

	// 启动服务
	return http.NewHttpServer(addr, root).Start()
}

// registerRoutes 注册路由
func (g *Gateway) registerRoutes(root *http.VortexHttpRouterGroup) {
	// 插件管理
	pluginGroup := root.AddGroup("/plugins")
	
	// 注册插件
	pluginGroup.AddRouter([]string{http.MethodPost}, "/register", g.handleRegister)
	
	// 注销插件
	pluginGroup.AddRouter([]string{http.MethodDelete}, "/:id", g.handleUnregister)
	
	// 获取插件列表
	pluginGroup.AddRouter([]string{http.MethodGet}, "/list", g.handleList)
	
	// 获取插件详情
	pluginGroup.AddRouter([]string{http.MethodGet}, "/:id", g.handleGet)
	
	// 执行插件方法
	pluginGroup.AddRouter([]string{http.MethodPost}, "/:id/execute/:method", g.handleExecute)
	
	// 健康检查
	pluginGroup.AddRouter([]string{http.MethodGet}, "/:id/health", g.handleHealth)
}

// handleRegister 处理注册请求
func (g *Gateway) handleRegister(ctx *http.Context) error {
	var req models.RegisterRequest
	if err := ctx.UnmarshalBody(&req); err != nil {
		return ctx.JsonResponse(http.Codes.BadRequest, pkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	resp, err := g.center.RegisterPlugin(ctx.Context(), &req)
	if err != nil {
		return ctx.JsonResponse(http.Codes.InternalError, pkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	return ctx.JsonResponse(http.Codes.Success, resp)
}

// handleUnregister 处理注销请求
func (g *Gateway) handleUnregister(ctx *http.Context) error {
	pluginID := ctx.Param("id")
	if err := g.center.UnregisterPlugin(ctx.Context(), pluginID); err != nil {
		return ctx.JsonResponse(http.Codes.InternalError, pkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	return ctx.JsonResponse(http.Codes.Success, pkg.VMsgResponse{
		"message": "Plugin unregistered successfully",
	})
}

// handleList 处理列表请求
func (g *Gateway) handleList(ctx *http.Context) error {
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")
	page := 1
	limit := 20

	// TODO: 解析 page 和 limit

	resp, err := g.center.ListPlugins(ctx.Context(), keyword, status, page, limit)
	if err != nil {
		return ctx.JsonResponse(http.Codes.InternalError, pkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	return ctx.JsonResponse(http.Codes.Success, resp)
}

// handleGet 处理获取详情请求
func (g *Gateway) handleGet(ctx *http.Context) error {
	pluginID := ctx.Param("id")
	
	plugin, err := g.center.GetPlugin(ctx.Context(), pluginID)
	if err != nil {
		return ctx.JsonResponse(http.Codes.InternalError, pkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	if plugin == nil {
		return ctx.JsonResponse(http.Codes.NotFound, pkg.VMsgResponse{
			"error": "Plugin not found",
		})
	}

	return ctx.JsonResponse(http.Codes.Success, plugin)
}

// handleExecute 处理执行请求
func (g *Gateway) handleExecute(ctx *http.Context) error {
	pluginID := ctx.Param("id")
	methodName := ctx.Param("method")

	var params map[string]interface{}
	if err := ctx.UnmarshalBody(&params); err != nil {
		// 允许空参数
		params = nil
	}

	resp, err := g.center.Execute(ctx.Context(), pluginID, methodName, params)
	if err != nil {
		return ctx.JsonResponse(http.Codes.InternalError, pkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	return ctx.JsonResponse(http.Codes.Success, resp)
}

// handleHealth 处理健康检查请求
func (g *Gateway) handleHealth(ctx *http.Context) error {
	pluginID := ctx.Param("id")
	
	health, err := g.center.HealthCheck(ctx.Context(), pluginID)
	if err != nil {
		return ctx.JsonResponse(http.Codes.InternalError, pkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	return ctx.JsonResponse(http.Codes.Success, health)
}
