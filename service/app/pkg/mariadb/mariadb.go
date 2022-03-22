package mariadb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vladazn/dhq/service/config"
)

type MariaDB interface {
	GetClient() *sql.DB
}

type MariaDBModel struct {
	client *sql.DB
}

func NewMariadbConnection(conf config.MariaDBConfigs) (MariaDB, error) {
	url := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true",
		conf.User, conf.Pass, conf.Host, conf.Port, conf.Db,
	)
	conn, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	return &MariaDBModel{client: conn}, nil
}

func (m MariaDBModel) GetClient() *sql.DB {
	return m.client
}
