package transport

import (
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func SetupRouter(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":    "Welcome to the API",
			"csrf_token": csrf.GetToken(c),
		})
	})

	r.POST("/users/register", RegisterUser)
	r.POST("/users/login", LoginUser)

	taskRouter := r.Group("/tasks")
	taskRouter.Use(middleware.AuthMiddleware())
	{
		taskRouter.GET("/", GetTasks)
		taskRouter.POST("/", CreateTask)
		taskRouter.PUT("/:id", UpdateTask)
		taskRouter.DELETE("/:id", DeleteTask)
		taskRouter.GET("/:id", GetTaskDetail)
	}

	adminRouter := r.Group("/admin")
	adminRouter.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		adminRouter.GET("/stats", GetStats)
	}
}
