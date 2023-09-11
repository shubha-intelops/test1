package daos

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/shubha-intelops/test1/login/pkg/rest/server/daos/clients/sqls"
	"github.com/shubha-intelops/test1/login/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
)

type LoginDao struct {
	sqlClient *sqls.MySQLClient
}

func migrateLogins(r *sqls.MySQLClient) error {
	query := `
	CREATE TABLE IF NOT EXISTS logins(
		ID int NOT NULL AUTO_INCREMENT,
        
		Password VARCHAR(100) NOT NULL,
		Username VARCHAR(100) NOT NULL,
	    PRIMARY KEY (ID)
	);
	`
	_, err := r.DB.Exec(query)
	return err
}

func NewLoginDao() (*LoginDao, error) {
	sqlClient, err := sqls.InitMySQLDB()
	if err != nil {
		return nil, err
	}
	err = migrateLogins(sqlClient)
	if err != nil {
		return nil, err
	}
	return &LoginDao{
		sqlClient,
	}, nil
}

func (loginDao *LoginDao) CreateLogin(m *models.Login) (*models.Login, error) {
	insertQuery := "INSERT INTO logins(Password, Username) values(?, ?)"
	res, err := loginDao.sqlClient.DB.Exec(insertQuery, m.Password, m.Username)
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
	log.Debugf("login created")
	return m, nil
}

func (loginDao *LoginDao) UpdateLogin(id int64, m *models.Login) (*models.Login, error) {
	if id == 0 {
		return nil, errors.New("invalid login ID")
	}
	if id != m.Id {
		return nil, errors.New("id and payload don't match")
	}

	login, err := loginDao.GetLogin(id)
	if err != nil {
		return nil, err
	}
	if login == nil {
		return nil, sql.ErrNoRows
	}

	updateQuery := "UPDATE logins SET Password = ?, Username = ? WHERE Id = ?"
	res, err := loginDao.sqlClient.DB.Exec(updateQuery, m.Password, m.Username, id)
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

	log.Debugf("login updated")
	return m, nil
}

func (loginDao *LoginDao) DeleteLogin(id int64) error {
	deleteQuery := "DELETE FROM logins WHERE Id = ?"
	res, err := loginDao.sqlClient.DB.Exec(deleteQuery, id)
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

	log.Debugf("login deleted")
	return nil
}

func (loginDao *LoginDao) ListLogins() ([]*models.Login, error) {
	selectQuery := "SELECT * FROM logins"
	rows, err := loginDao.sqlClient.DB.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var logins []*models.Login
	for rows.Next() {
		m := models.Login{}
		if err = rows.Scan(&m.Id, &m.Password, &m.Username); err != nil {
			return nil, err
		}
		logins = append(logins, &m)
	}
	if logins == nil {
		logins = []*models.Login{}
	}
	log.Debugf("login listed")
	return logins, nil
}

func (loginDao *LoginDao) GetLogin(id int64) (*models.Login, error) {
	selectQuery := "SELECT * FROM logins WHERE Id = ?"
	row := loginDao.sqlClient.DB.QueryRow(selectQuery, id)

	m := models.Login{}
	if err := row.Scan(&m.Id, &m.Password, &m.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}
	log.Debugf("login retrieved")
	return &m, nil
}
