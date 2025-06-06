package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"goblog/ent"
	"goblog/internal/config"
	"goblog/internal/domain"
	"goblog/internal/handler"
	"goblog/internal/middleware"
	"goblog/internal/pkg/logger"
	"goblog/internal/repository"
	"goblog/internal/service"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// 检查命令行参数
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--health-check":
			healthCheck()
			return
		case "--migrate-only":
			runMigrationOnly()
			return
		case "--version":
			fmt.Println("goblog version 1.0.0")
			return
		}
	}

	// 初始化日志
	logger.Init()

	// 加载配置
	cfg := config.Load()

	// 创建Ent客户端
	client, err := ent.Open(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed opening connection to database: %v", err)
	}
	defer client.Close()

	// 运行自动迁移
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// 初始化仓储层
	articleRepo := repository.NewArticleRepository(client)
	categoryRepo := repository.NewCategoryRepository(client)
	tagRepo := repository.NewTagRepository(client)

	// 初始化服务层
	authService := service.NewAuthService(cfg)
	articleService := service.NewArticleService(articleRepo, categoryRepo, tagRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)

	// 初始化中间件
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// 初始化处理器
	articleHandler := handler.NewArticleHandler(articleService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)

	// 创建Echo实例
	e := echo.New()

	// 全局中间件
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// 路由组
	api := e.Group("/api")

	// 公开路由（读操作）
	setupPublicRoutes(api, articleHandler, categoryHandler, tagHandler)

	// 需要认证的路由（写操作）
	setupAuthRoutes(api, authMiddleware, articleHandler, categoryHandler, tagHandler)

	// 认证路由
	setupAuthEndpoints(e, authService)

	// 健康检查端点
	setupHealthCheck(e)

	// 启动服务器
	logger.Info("博客服务器启动", "port", cfg.Server.Port)
	log.Fatal(e.Start(cfg.Server.Port))
}

// setupPublicRoutes 设置公开路由
func setupPublicRoutes(api *echo.Group, articleHandler *handler.ArticleHandler, categoryHandler *handler.CategoryHandler, tagHandler *handler.TagHandler) {
	// 文章路由
	api.GET("/articles", articleHandler.List)
	api.GET("/articles/:id", articleHandler.GetByID)
	api.GET("/articles/category/:categoryId", articleHandler.ListByCategory)
	api.GET("/articles/tag/:tagId", articleHandler.ListByTag)

	// 分类路由
	api.GET("/categories", categoryHandler.List)
	api.GET("/categories/:id", categoryHandler.GetByID)

	// 标签路由
	api.GET("/tags", tagHandler.List)
	api.GET("/tags/:id", tagHandler.GetByID)
}

// setupAuthRoutes 设置需要认证的路由
func setupAuthRoutes(api *echo.Group, authMiddleware *middleware.AuthMiddleware, articleHandler *handler.ArticleHandler, categoryHandler *handler.CategoryHandler, tagHandler *handler.TagHandler) {
	authGroup := api.Group("", authMiddleware.RequireAuth())

	// 文章管理
	authGroup.POST("/articles", articleHandler.Create)
	authGroup.PUT("/articles/:id", articleHandler.Update)
	authGroup.DELETE("/articles/:id", articleHandler.Delete)

	// 文章备份
	authGroup.GET("/articles/backup", articleHandler.Backup)

	// 分类管理
	authGroup.POST("/categories", categoryHandler.Create)
	authGroup.PUT("/categories/:id", categoryHandler.Update)
	authGroup.DELETE("/categories/:id", categoryHandler.Delete)

	// 标签管理
	authGroup.POST("/tags", tagHandler.Create)
	authGroup.PUT("/tags/:id", tagHandler.Update)
	authGroup.DELETE("/tags/:id", tagHandler.Delete)
}

// setupAuthEndpoints 设置认证端点
func setupAuthEndpoints(e *echo.Echo, authService domain.AuthService) {
	// 直接处理登录，避免循环依赖
	e.POST("/auth/login", func(c echo.Context) error {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, map[string]string{"error": "无效的请求参数"})
		}

		loginReq := &domain.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		}

		resp, err := authService.Login(c.Request().Context(), loginReq)
		if err != nil {
			return c.JSON(401, map[string]string{"error": "用户名或密码错误"})
		}

		return c.JSON(200, resp)
	})
}

// setupHealthCheck 设置健康检查端点
func setupHealthCheck(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "goblog",
			"version":   "1.0.0",
		})
	})
}

// runMigrationOnly 仅运行数据库迁移
func runMigrationOnly() {
	log.Println("运行数据库迁移...")

	// 初始化日志
	logger.Init()

	// 加载配置
	cfg := config.Load()

	// 创建Ent客户端
	client, err := ent.Open(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed opening connection to database: %v", err)
	}
	defer client.Close()

	// 运行自动迁移
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("数据库迁移完成")
}

// healthCheck 执行健康检查
func healthCheck() {
	// 获取服务端口，默认为8080
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	// 创建HTTP客户端，设置超时
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	// 尝试访问健康检查端点
	healthURL := fmt.Sprintf("http://localhost%s/health", port)
	resp, err := client.Get(healthURL)
	if err != nil {
		log.Printf("健康检查失败: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		log.Printf("健康检查失败: HTTP状态码 %d", resp.StatusCode)
		os.Exit(1)
	}

	log.Println("健康检查通过")
	os.Exit(0)
}
