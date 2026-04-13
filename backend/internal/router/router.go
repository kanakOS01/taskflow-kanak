package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"taskflow/internal/auth"
	"taskflow/internal/config"
	"taskflow/internal/middleware"
	"taskflow/internal/projects"
	"taskflow/internal/tasks"
	"taskflow/internal/users"
)

func NewRouter(db *pgxpool.Pool, cfg *config.Config) *gin.Engine {
	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Dependencies
	userRepo := users.NewRepository(db)
	projectRepo := projects.NewRepository(db)
	taskRepo := tasks.NewRepository(db)

	authService := auth.NewService(userRepo, cfg.JWTSecret)
	userService := users.NewService(userRepo)
	projectService := projects.NewService(projectRepo)
	taskService := tasks.NewService(taskRepo, projectRepo, userRepo)

	authHandler := auth.NewHandler(authService)
	usersHandler := users.NewHandler(userService)
	projectHandler := projects.NewHandler(projectService)
	tasksHandler := tasks.NewHandler(taskService)

	// Routes
	authGroup := r.Group("/auth")
	authHandler.RegisterRoutes(authGroup)

	// Protected routes
	api := r.Group("/")
	api.Use(middleware.RequireAuth(cfg.JWTSecret))

	usersGroup := api.Group("/users")
	usersHandler.RegisterRoutes(usersGroup)

	projectsGroup := api.Group("/projects")
	projectHandler.RegisterRoutes(projectsGroup)

	tasksHandler.RegisterRoutes(api)

	return r
}
