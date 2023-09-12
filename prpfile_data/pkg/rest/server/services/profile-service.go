package services

import (
	"github.com/shubha-intelops/test1/prpfile_data/pkg/rest/server/daos"
	"github.com/shubha-intelops/test1/prpfile_data/pkg/rest/server/models"
)

type ProfileService struct {
	profileDao *daos.ProfileDao
}

func NewProfileService() (*ProfileService, error) {
	profileDao, err := daos.NewProfileDao()
	if err != nil {
		return nil, err
	}
	return &ProfileService{
		profileDao: profileDao,
	}, nil
}

func (profileService *ProfileService) CreateProfile(profile *models.Profile) (*models.Profile, error) {
	return profileService.profileDao.CreateProfile(profile)
}

func (profileService *ProfileService) UpdateProfile(id int64, profile *models.Profile) (*models.Profile, error) {
	return profileService.profileDao.UpdateProfile(id, profile)
}

func (profileService *ProfileService) DeleteProfile(id int64) error {
	return profileService.profileDao.DeleteProfile(id)
}

func (profileService *ProfileService) ListProfiles() ([]*models.Profile, error) {
	return profileService.profileDao.ListProfiles()
}

func (profileService *ProfileService) GetProfile(id int64) (*models.Profile, error) {
	return profileService.profileDao.GetProfile(id)
}
