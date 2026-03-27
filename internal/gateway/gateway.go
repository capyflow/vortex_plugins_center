package gateway

import (
	"context"
	"net/http"

	"plugin-platform/internal/center"
	"plugin-platform/pkg/models"

	vortex "github.com/capyflow/vortexv3"
	vpkg "github.com/capyflow/vortexv3/pkg"
	vhttp "github.com/capyflow/vortexv3/server/http"
	"github.com/gin-gonic/gin"
)

// Gateway API 网关
type Gateway struct {
	center *center.PluginCenter
}

// New 创建网关
func New(pc *center.PluginCenter) *Gateway {
	return &Gateway{
		center: pc,
	}
}

// Start 启动网关
func (g *Gateway) Start(ctx context.Context, port int) error {
	root := vhttp.NewRootGroup("/v1/api")

	// 添加全局 OPTIONS 路由处理 CORS 预检
	root.AddRouter([]string{http.MethodOptions}, "/*path", func(ctx *vhttp.Context) error {
		c := ctx.GinContext()
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.AbortWithStatus(http.StatusNoContent)
		return nil
	})

	g.registerRoutes(root)

	// 创建 Vortex HTTP 引擎
	e := vortex.NewVortexEngine(ctx,
		vortex.WithEnableProtocol([]string{vpkg.HTTP}),
		vortex.WithPort(port),
		vortex.WithHttpRouterRootGroup(root),
		vortex.WithHttpServerOptions(func(e *gin.Engine) {
			e.Use(func(c *gin.Context) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				c.Writer.Header().Set("Access-Control-Max-Age", "86400")

				if c.Request.Method == "OPTIONS" {
					c.AbortWithStatus(http.StatusNoContent)
					return
				}
				c.Next()
			})
		}),
	)
	// 启动服务
	e.Start()
	return nil
}

// 注意：路由按注册顺序匹配，所以 /list 必须在 /:id 之前注册
func (g *Gateway) registerRoutes(root *vhttp.VortexHttpRouterGroup) {
	// 插件管理 - 按特异性排序（最具体的路由先注册）
	pluginGroup := root.AddGroup("/plugins")

	// 1. 获取插件列表（最具体，无路径参数）
	pluginGroup.AddRouter([]string{http.MethodGet}, "/list", g.handleList)

	// 2. 注册插件
	pluginGroup.AddRouter([]string{http.MethodPost}, "/register", g.handleRegister)

	// 3. 执行插件方法（二级路径）
	pluginGroup.AddRouter([]string{http.MethodPost}, "/:id/execute/:method", g.handleExecute)

	// 4. 健康检查（二级路径）
	pluginGroup.AddRouter([]string{http.MethodGet}, "/:id/health", g.handleHealth)

	// 5. 获取插件详情（一级路径参数）
	pluginGroup.AddRouter([]string{http.MethodGet}, "/:id", g.handleGet)

	// 6. 注销插件（一级路径参数）
	pluginGroup.AddRouter([]string{http.MethodDelete}, "/:id", g.handleUnregister)
}

// handleRegister 处理注册请求
func (g *Gateway) handleRegister(ctx *vhttp.Context) error {
	var req models.RegisterRequest
	if err := ctx.UnmarshalBody(&req); err != nil {
		return ctx.JsonResponse(vhttp.Codes.BadRequest, vpkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	resp, err := g.center.RegisterPlugin(ctx.Context(), &req)
	if err != nil {
		return ctx.JsonResponse(vhttp.Codes.InternalError, vpkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, resp)
}

// handleUnregister 处理注销请求
func (g *Gateway) handleUnregister(ctx *vhttp.Context) error {
	pluginID := ctx.GinContext().Param("id")
	if err := g.center.UnregisterPlugin(ctx.Context(), pluginID); err != nil {
		return ctx.JsonResponse(vhttp.Codes.InternalError, vpkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, vpkg.VMsgResponse{
		"message": "Plugin unregistered successfully",
	})
}

// handleList 处理列表请求
func (g *Gateway) handleList(ctx *vhttp.Context) error {
	keyword := ctx.GinContext().Query("keyword")
	status := ctx.GinContext().Query("status")
	page := 1
	limit := 20

	// TODO: 解析 page 和 limit

	resp, err := g.center.ListPlugins(ctx.Context(), keyword, status, page, limit)
	if err != nil {
		return ctx.JsonResponse(vhttp.Codes.InternalError, vpkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, resp)
}

// handleGet 处理获取详情请求
func (g *Gateway) handleGet(ctx *vhttp.Context) error {
	pluginID := ctx.GinContext().Param("id")

	plugin, err := g.center.GetPlugin(ctx.Context(), pluginID)
	if err != nil {
		return ctx.JsonResponse(vhttp.Codes.InternalError, vpkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	if plugin == nil {
		return ctx.JsonResponse(vhttp.Codes.BadRequest, vpkg.VMsgResponse{
			"error": "Plugin not found",
		})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, plugin)
}

// handleExecute 处理执行请求
func (g *Gateway) handleExecute(ctx *vhttp.Context) error {
	pluginID := ctx.GinContext().Param("id")
	methodName := ctx.GinContext().Param("method")

	var params map[string]interface{}
	if err := ctx.UnmarshalBody(&params); err != nil {
		// 允许空参数
		params = nil
	}

	resp, err := g.center.Execute(ctx.Context(), pluginID, methodName, params)
	if err != nil {
		return ctx.JsonResponse(vhttp.Codes.InternalError, vpkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, resp)
}

// handleHealth 处理健康检查请求
func (g *Gateway) handleHealth(ctx *vhttp.Context) error {
	pluginID := ctx.GinContext().Param("id")

	health, err := g.center.HealthCheck(ctx.Context(), pluginID)
	if err != nil {
		return ctx.JsonResponse(vhttp.Codes.InternalError, vpkg.VMsgResponse{
			"error": err.Error(),
		})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, health)
}
