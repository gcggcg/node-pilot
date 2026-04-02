package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"node-pilot/internal/config"
	"node-pilot/internal/handler"
	"node-pilot/internal/logger"
	"node-pilot/internal/middleware"
	"node-pilot/internal/repository"
	"node-pilot/internal/service"
	"node-pilot/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	dbPath := flag.String("db", "./data/servers.db", "SQLite database path")
	listen := flag.String("listen", ":8080", "Listen address")
	debug := flag.Bool("debug", false, "Enable debug logging")
	logFile := flag.String("log", "", "Log file path (default: stdout)")
	flag.Parse()

	// Initialize logger
	logger.Init(*debug)
	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		logger.SetOutput(f)
	}

	dbDir := filepath.Dir(*dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	cfg := config.Config{
		DBPath: *dbPath,
		Listen: *listen,
		Debug:  *debug,
	}

	logger.Info("==========================================")
	logger.Info("  NodePilot 批量服务器管理平台")
	logger.Info("==========================================")
	logger.Info("Debug模式: %v", *debug)
	logger.Info("数据库: %s", *dbPath)
	logger.Info("监听地址: %s", *listen)

	db, err := repository.NewDB(cfg.DBPath)
	if err != nil {
		logger.Error("Failed to open database: %v", err)
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	sshPool := service.NewSSHPool()
	wsHub := websocket.NewHub()
	go wsHub.Run()
	taskSvc := service.NewTaskExecutor(repo, sshPool, wsHub, *debug)

	h := handler.NewHandler(repo, sshPool, wsHub, taskSvc)

	jwtSecret := "node-pilot-jwt-secret-key-32bytes!" // 32 bytes for HS256
	authHandler := handler.NewAuthHandler(repo, jwtSecret)
	userHandler := handler.NewUserHandler(repo, jwtSecret)

	fileUploadService := service.NewFileUploadService(repo, sshPool, "./data/files")
	fileUploadHandler := handler.NewFileUploadHandler(repo, fileUploadService, "./data/files")

	if *debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api")
	{
		auth := api.Group("/v1/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.JWTAuth(jwtSecret), authHandler.Me)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.PUT("/profile", middleware.JWTAuth(jwtSecret), authHandler.UpdateProfile)
			auth.PUT("/password", middleware.JWTAuth(jwtSecret), authHandler.ChangePassword)
		}

		admin := api.Group("/v1/admin")
		admin.Use(middleware.JWTAuth(jwtSecret))
		admin.Use(middleware.RequireRole("ROLE_ADMIN"))
		{
			admin.GET("/users", userHandler.ListUsers)
			admin.POST("/users", userHandler.CreateUser)
			admin.DELETE("/users/:id", userHandler.DeleteUsers)
			admin.POST("/users/batch-delete", userHandler.DeleteUsers)
		}

		servers := api.Group("/servers")
		{
			servers.GET("", h.ListServers)
			servers.GET("/:id", h.GetServer)
			servers.POST("", h.CreateServer)
			servers.PUT("/:id", h.UpdateServer)
			servers.DELETE("/:id", h.DeleteServer)
			servers.POST("/batch-delete", h.DeleteServers)
			servers.POST("/:id/test", h.TestServerConnection)
		}

		scripts := api.Group("/scripts")
		{
			scripts.GET("", h.ListScripts)
			scripts.GET("/:id", h.GetScript)
			scripts.POST("", h.CreateScript)
			scripts.PUT("/:id", h.UpdateScript)
			scripts.DELETE("/:id", h.DeleteScript)
			scripts.POST("/batch-delete", h.DeleteScripts)
		}

		tasks := api.Group("/tasks")
		{
			tasks.GET("", h.ListTasks)
			tasks.GET("/:id", h.GetTask)
			tasks.POST("", h.CreateTask)
			tasks.DELETE("/:id", h.CancelTask)
			tasks.POST("/batch-delete", h.DeleteTasks)
			tasks.GET("/:id/output", h.GetTaskOutput)
		}

		api.POST("/upload", h.UploadFile)
		api.POST("/deploy", h.DeployFile)

		fileUploads := api.Group("/v1/file-uploads")
		fileUploads.Use(middleware.JWTAuth(jwtSecret))
		{
			fileUploads.GET("", fileUploadHandler.ListFileUploads)
			fileUploads.POST("", fileUploadHandler.CreateFileUpload)
			fileUploads.PUT("/:id", fileUploadHandler.UpdateFileUpload)
			fileUploads.DELETE("", fileUploadHandler.DeleteFileUploads)
			fileUploads.POST("/:id/execute", fileUploadHandler.ExecuteFileUpload)
			fileUploads.GET("/:id/results", fileUploadHandler.GetFileUploadResults)
		}

		api.POST("/v1/file-uploads/upload-file", fileUploadHandler.UploadFileToStorage)
	}

	r.GET("/ws", h.WebSocketHandler)

	exePath, _ := os.Executable()
	webRoot := filepath.Dir(exePath)
	assetsDir := filepath.Join(webRoot, "web", "assets")
	indexFile := filepath.Join(webRoot, "web", "index.html")
	logger.Info("Web根目录: %s", webRoot)
	logger.Info("静态资源: %s", assetsDir)
	r.StaticFS("/assets", gin.Dir(assetsDir, false))

	r.NoRoute(func(c *gin.Context) {
		logger.Debug("NoRoute: %s", c.Request.URL.Path)
		c.File(indexFile)
	})

	logger.Info("NodePilot starting on %s", *listen)
	if err := r.Run(*listen); err != nil {
		logger.Error("Failed to start server: %v", err)
		log.Fatalf("Failed to start server: %v", err)
	}
}
