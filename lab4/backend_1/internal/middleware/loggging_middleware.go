package middleware

import (
	logging "backend/internal"
	"backend/internal/monitoring"
	"backend/internal/notifications"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		monitoring.RequestCount.WithLabelValues(c.Request.Method, c.FullPath(), strconv.Itoa(status)).Inc()
		monitoring.RequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)

		monitoring.UpdateMetrics(time.Since(start), status)

		webhookURL := "WEBHOOK"
		chatID := "CHAT_ID"

		if status >= 400 {
			logging.Logger.WithFields(map[string]interface{}{
				"method":  c.Request.Method,
				"path":    c.Request.URL.Path,
				"status":  status,
				"latency": duration,
			}).Error("Сервер вернул ошибку")

			message := "Критическая ошибка на маршруте " + c.Request.URL.Path +
				"\nСтатус: " + strconv.Itoa(status) +
				"\nВремя ответа: " + strconv.FormatFloat(duration, 'f', 2, 64) + " секунд"
			_ = notifications.SendTelegramNotification(webhookURL, chatID, message)
		} else if duration > 2 {
			logging.Logger.WithFields(map[string]interface{}{
				"method":  c.Request.Method,
				"path":    c.Request.URL.Path,
				"status":  status,
				"latency": duration,
			}).Warn("Медленный ответ")

			message := "Медленный ответ на маршруте " + c.Request.URL.Path +
				"\nСтатус: " + strconv.Itoa(status) +
				"\nВремя ответа: " + strconv.FormatFloat(duration, 'f', 2, 64) + " секунд"
			_ = notifications.SendTelegramNotification(webhookURL, chatID, message)
		}

		logging.Logger.WithFields(map[string]interface{}{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  status,
			"latency": duration,
		}).Info("HTTP запрос обработан")

	}
}
