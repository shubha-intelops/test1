package daos

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/shubha-intelops/test1/prpfile_data/pkg/rest/server/daos/clients/sqls"
	"github.com/shubha-intelops/test1/prpfile_data/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
)

type ProfileDao struct {
	sqlClient *sqls.MySQLClient
}

func migrateProfiles(r *sqls.MySQLClient) error {
	query := `
	CREATE TABLE IF NOT EXISTS profiles(
		ID int NOT NULL AUTO_INCREMENT,
        
		Address VARCHAR(100) NOT NULL,
		Name VARCHAR(100) NOT NULL,
		Age INT NOT NULL,
	    PRIMARY KEY (ID)
	);
	`
	_, err := r.DB.Exec(query)
	return err
}

func NewProfileDao() (*ProfileDao, error) {
	sqlClient, err := sqls.InitMySQLDB()
	if err != nil {
		return nil, err
	}
	err = migrateProfiles(sqlClient)
	if err != nil {
		return nil, err
	}
	return &ProfileDao{
		sqlClient,
	}, nil
}

func (profileDao *ProfileDao) CreateProfile(m *models.Profile) (*models.Profile, error) {
	insertQuery := "INSERT INTO profiles(Address, Name, Age) values(?, ?, ?)"
	res, err := profileDao.sqlClient.DB.Exec(insertQuery, m.Address, m.Name, m.Age)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return nil, sqls.ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.Id = id
	log.Debugf("profile created")
	return m, nil
}

func (profileDao *ProfileDao) UpdateProfile(id int64, m *models.Profile) (*models.Profile, error) {
	if id == 0 {
		return nil, errors.New("invalid profile ID")
	}
	if id != m.Id {
		return nil, errors.New("id and payload don't match")
	}

	profile, err := profileDao.GetProfile(id)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, sql.ErrNoRows
	}

	updateQuery := "UPDATE profiles SET Address = ?, Name = ?, Age = ? WHERE Id = ?"
	res, err := profileDao.sqlClient.DB.Exec(updateQuery, m.Address, m.Name, m.Age, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sqls.ErrUpdateFailed
	}

	log.Debugf("profile updated")
	return m, nil
}

func (profileDao *ProfileDao) DeleteProfile(id int64) error {
	deleteQuery := "DELETE FROM profiles WHERE Id = ?"
	res, err := profileDao.sqlClient.DB.Exec(deleteQuery, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sqls.ErrDeleteFailed
	}

	log.Debugf("profile deleted")
	return nil
}

func (profileDao *ProfileDao) ListProfiles() ([]*models.Profile, error) {
	selectQuery := "SELECT * FROM profiles"
	rows, err := profileDao.sqlClient.DB.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var profiles []*models.Profile
	for rows.Next() {
		m := models.Profile{}
		if err = rows.Scan(&m.Id, &m.Address, &m.Name, &m.Age); err != nil {
			return nil, err
		}
		profiles = append(profiles, &m)
	}
	if profiles == nil {
		profiles = []*models.Profile{}
	}
	log.Debugf("profile listed")
	return profiles, nil
}

func (profileDao *ProfileDao) GetProfile(id int64) (*models.Profile, error) {
	selectQuery := "SELECT * FROM profiles WHERE Id = ?"
	row := profileDao.sqlClient.DB.QueryRow(selectQuery, id)

	m := models.Profile{}
	if err := row.Scan(&m.Id, &m.Address, &m.Name, &m.Age); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}
	log.Debugf("profile retrieved")
	return &m, nil
}
