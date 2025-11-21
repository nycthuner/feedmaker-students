package handler

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"

	"github.com/anaclaraddias/feedmaker-students/src/domain/entity"
	"github.com/anaclaraddias/feedmaker-students/src/domain/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginHandler struct {
	db *gorm.DB
}

func NewLoginHandler(db *gorm.DB) HandlerInterface {
	return &LoginHandler{db: db}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (handler *LoginHandler) Handle(context *gin.Context) {
	context.Writer.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashPassword := sha512.Sum512([]byte(req.Password))
	password := hex.EncodeToString(hashPassword[:])

	loginUser := entity.User{
		Username: req.Username,
		Password: password,
	}

	jwt, err := service.NewLoginService(handler.db).Execute(loginUser)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, map[string]any{"token": jwt})
}
