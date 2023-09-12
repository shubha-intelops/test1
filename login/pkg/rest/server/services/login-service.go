package services

import (
	"github.com/shubha-intelops/test1/login/pkg/rest/server/daos"
	"github.com/shubha-intelops/test1/login/pkg/rest/server/models"
)

type LoginService struct {
	loginDao *daos.LoginDao
}

func NewLoginService() (*LoginService, error) {
	loginDao, err := daos.NewLoginDao()
	if err != nil {
		return nil, err
	}
	return &LoginService{
		loginDao: loginDao,
	}, nil
}

func (loginService *LoginService) CreateLogin(login *models.Login) (*models.Login, error) {
	return loginService.loginDao.CreateLogin(login)
}

func (loginService *LoginService) UpdateLogin(id int64, login *models.Login) (*models.Login, error) {
	return loginService.loginDao.UpdateLogin(id, login)
}

func (loginService *LoginService) DeleteLogin(id int64) error {
	return loginService.loginDao.DeleteLogin(id)
}

func (loginService *LoginService) ListLogins() ([]*models.Login, error) {
	return loginService.loginDao.ListLogins()
}

func (loginService *LoginService) GetLogin(id int64) (*models.Login, error) {
	return loginService.loginDao.GetLogin(id)
}
