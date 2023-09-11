package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shubha-intelops/test1/login/pkg/rest/server/models"
	"github.com/shubha-intelops/test1/login/pkg/rest/server/services"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type LoginController struct {
	loginService *services.LoginService
}

func NewLoginController() (*LoginController, error) {
	loginService, err := services.NewLoginService()
	if err != nil {
		return nil, err
	}
	return &LoginController{
		loginService: loginService,
	}, nil
}

func (loginController *LoginController) CreateLogin(context *gin.Context) {
	// validate input
	var input models.Login
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger login creation
	if _, err := loginController.loginService.CreateLogin(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Login created successfully"})
}

func (loginController *LoginController) UpdateLogin(context *gin.Context) {
	// validate input
	var input models.Login
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger login update
	if _, err := loginController.loginService.UpdateLogin(id, &input); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login updated successfully"})
}

func (loginController *LoginController) FetchLogin(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger login fetching
	login, err := loginController.loginService.GetLogin(id)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, login)
}

func (loginController *LoginController) DeleteLogin(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger login deletion
	if err := loginController.loginService.DeleteLogin(id); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Login deleted successfully",
	})
}

func (loginController *LoginController) ListLogins(context *gin.Context) {
	// trigger all logins fetching
	logins, err := loginController.loginService.ListLogins()
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, logins)
}

func (*LoginController) PatchLogin(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "PATCH",
	})
}

func (*LoginController) OptionsLogin(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "OPTIONS",
	})
}

func (*LoginController) HeadLogin(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "HEAD",
	})
}
