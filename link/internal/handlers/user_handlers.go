package handlers

import (
	"link/internal/config"
	"link/internal/models"
	"link/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
	Config *config.Config
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return utils.NewBadRequestError("Invalid request payload")
	}

	if user.Username == "" || user.Password == "" {
		return utils.NewBadRequestError("Username and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewInternalServerError("Failed to hash password")
	}

	user.Password = string(hashedPassword)

	if err := h.DB.Create(user).Error; err != nil {
		return utils.NewInternalServerError("Failed to register user")
	}

	user.Password = ""
	return c.JSON(201, user)
}

func (h *UserHandler) Login(c echo.Context) error {
	loginRequest := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.Bind(&loginRequest); err != nil {
		return utils.NewBadRequestError("Invalid request payload")
	}

	var user models.User
	if err := h.DB.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
		return utils.NewUnauthorizedError("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return utils.NewUnauthorizedError("Invalid credentials")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(h.Config.Auth.SecretKey))
	if err != nil {
		return utils.NewInternalServerError("Failed to generate token")
	}

	return c.JSON(200, map[string]string{
		"token": tokenString,
		"message": "Login successful",
	})
}