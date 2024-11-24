package transport

import (
	logging "backend/internal"
	"backend/internal/models"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("secret_key")
var users = make(map[string]*models.User)
var mu sync.Mutex
var userCount uint = 0
var validate = validator.New()

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logging.Logger.Warn("Ошибка привязки JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation
	if err := validate.Struct(user); err != nil {
		logging.Logger.Error("Validation error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logging.Logger.Error("Ошибка хеширования пароля: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}
	user.Password = string(hashedPassword)

	mu.Lock()
	if _, exists := users[user.Email]; exists {
		mu.Unlock()
		logging.Logger.Warn("Пользователь с таким email уже существует: ", user.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким email уже существует"})
		return
	}

	users[user.Email] = &user
	mu.Unlock()

	logging.Logger.Info("Пользователь зарегистрирован: ", user.Email)
	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован"})
}

func LoginUser(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
		return
	}

	mu.Lock()
	user, exists := users[loginRequest.Email]
	mu.Unlock()

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	fmt.Println("Захешированный пароль в памяти:", user.Password)
	fmt.Println("Введенный пароль для входа:", loginRequest.Password)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		fmt.Println("Ошибка сравнения пароля:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
		return
	}

	logging.Logger.Info("Пользователь залогинился: ", user.Email)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
