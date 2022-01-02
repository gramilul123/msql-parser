package mysql

import (
	"context"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"mysql-parser/config"
)

// NewMysql returns mysql instance and defer function or error
func NewMysql(db *config.Db) (*sqlx.DB, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}

	mysqlCfg := mysql.NewConfig()
	mysqlCfg.Addr = db.Host + ":" + db.Port
	mysqlCfg.Net = "tcp"
	mysqlCfg.User = db.User
	mysqlCfg.Passwd = db.Password
	mysqlCfg.DBName = db.Name
	mysqlCfg.ParseTime = true
	mysqlCfg.Params = db.Params
	mysqlCfg.Loc = loc

	dbConnection, err := sqlx.Open("mysql", mysqlCfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	dbConnection.SetMaxOpenConns(db.MaxOpenConnections)
	dbConnection.SetMaxIdleConns(db.MaxIdleConnections)
	dbConnection.SetConnMaxLifetime(time.Minute * time.Duration(db.MaxLifetime))
	dbConnection.SetConnMaxIdleTime(time.Minute * time.Duration(db.MaxIdleLifetime))

	var bgCtx = context.Background()
	var ctxPingTimeout, pingCancelFunc = context.WithTimeout(bgCtx, time.Second*time.Duration(db.MaxPingTimeout))
	defer pingCancelFunc()

	if err := dbConnection.PingContext(ctxPingTimeout); err != nil {
		if err := dbConnection.Close(); err != nil {
			return nil, err
		}
		return nil, err
	}

	return dbConnection, nil
}
