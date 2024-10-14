package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VikaPaz/task_tracker/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type AuthHandler struct {
	router   *gin.Engine
	service  AuthServise
	validate *validator.Validate
	log      *zerolog.Logger
}

func NewAuthHandler(svc AuthServise, log *zerolog.Logger) *AuthHandler {
	router := gin.Default()
	validate := validator.New()
	return &AuthHandler{
		router:   router,
		service:  svc,
		validate: validate,
		log:      log,
	}
}

type AuthServise interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GenerateToken(ctx context.Context, user models.User) (string, error)
}

func (h *AuthHandler) registerRoutes() {
	tasks := h.router.Group("/auth")
	{
		tasks.POST("/", h.signUp)
		tasks.GET("/", h.signIn)
	}
	h.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func RunAuthServ(server *AuthHandler, serverPort string) {
	server.registerRoutes()
	server.router.Run(":" + serverPort)
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body todo.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *AuthHandler) signUp(c *gin.Context) {
	var req models.User

	if err := c.BindJSON(&req); err != nil {
		Response(h.log, c, nil, http.StatusBadRequest, fmt.Errorf("failed to bind request JSON: %w", err))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		Response(h.log, c, nil, http.StatusBadRequest, fmt.Errorf("failed validate: %w", err))
		return
	}

	id, err := h.service.CreateUser(c, req)
	if err != nil {
		Response(h.log, c, nil, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInResponse struct {
	Login    string `form:"login" validate:"required"`
	Password string `form:"password" validate:"required"`
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *AuthHandler) signIn(c *gin.Context) {
	var req signInResponse

	if err := c.ShouldBindQuery(&req); err != nil {
		Response(h.log, c, nil, http.StatusBadRequest, fmt.Errorf("error: %w", err))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		Response(h.log, c, nil, http.StatusBadRequest, fmt.Errorf("failed validate: %w", err))
		return
	}

	user := models.User{
		Login:    req.Login,
		Password: req.Password,
	}

	token, err := h.service.GenerateToken(c, user)
	if err != nil {
		Response(h.log, c, nil, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
