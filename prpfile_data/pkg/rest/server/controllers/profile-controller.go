package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shubha-intelops/test1/prpfile_data/pkg/rest/server/models"
	"github.com/shubha-intelops/test1/prpfile_data/pkg/rest/server/services"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ProfileController struct {
	profileService *services.ProfileService
}

func NewProfileController() (*ProfileController, error) {
	profileService, err := services.NewProfileService()
	if err != nil {
		return nil, err
	}
	return &ProfileController{
		profileService: profileService,
	}, nil
}

func (profileController *ProfileController) CreateProfile(context *gin.Context) {
	// validate input
	var input models.Profile
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger profile creation
	if _, err := profileController.profileService.CreateProfile(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Profile created successfully"})
}

func (profileController *ProfileController) UpdateProfile(context *gin.Context) {
	// validate input
	var input models.Profile
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

	// trigger profile update
	if _, err := profileController.profileService.UpdateProfile(id, &input); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (profileController *ProfileController) FetchProfile(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger profile fetching
	profile, err := profileController.profileService.GetProfile(id)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, profile)
}

func (profileController *ProfileController) DeleteProfile(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger profile deletion
	if err := profileController.profileService.DeleteProfile(id); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Profile deleted successfully",
	})
}

func (profileController *ProfileController) ListProfiles(context *gin.Context) {
	// trigger all profiles fetching
	profiles, err := profileController.profileService.ListProfiles()
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, profiles)
}

func (*ProfileController) PatchProfile(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "PATCH",
	})
}

func (*ProfileController) OptionsProfile(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "OPTIONS",
	})
}

func (*ProfileController) HeadProfile(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "HEAD",
	})
}
