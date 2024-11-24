package app

import (
	logging "backend/internal"
	"backend/internal/middleware"
	"backend/internal/monitoring"
	"backend/internal/transport"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"time"
)

func Run() {

	logging.InitLogger()
	logging.Logger.Info("Инициализация приложения")
	router := gin.Default()
	router.Use(middleware.RequestLogging())
	monitoring.InitMonitoring(router)
	go monitoring.CollectExpvarMetrics(router)

	for i := 0; i < 10000; i++ {
		logging.Logger.Info("Тестовая запись лога: ", i)
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	csrfSecret := "32-byte-long-auth-key"

	router.Use(csrf.Middleware(csrf.Options{
		Secret: csrfSecret,
		ErrorFunc: func(c *gin.Context) {
			logging.Logger.Warn("CSRF token mismatch")
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	transport.SetupRouter(router)

	router.Use(middleware.SecurityHeaders())

	err := router.RunTLS(":8080", "cert.pem", "key.pem")
	if err != nil {
		logging.Logger.Fatal("Не удалось запустить сервер с TLS: ", err)
	}
}
