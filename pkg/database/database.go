package database

import (
	"database/sql"
	"github.com/ahmadrezamusthafa/logwatcher/common/logger"
	"github.com/ahmadrezamusthafa/logwatcher/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/segmentio/go-athena"
)

type Result struct {
	Data  interface{}
	Error error
}

type Block func(db *sqlx.Tx, c chan Result)

type Database struct {
	Config config.Config `inject:"config"`
	db     *sql.DB
}

func (m Database) GetDB() *sql.DB {
	return m.db
}

func (m Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.db.QueryRow(query, args)
}

func (m Database) In(query string, params map[string]interface{}) (string, []interface{}, error) {
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return query, args, err
	}
	return sqlx.In(query, args...)
}

func (m Database) IsReady() bool {
	if m.db == nil {
		return false
	}
	if err := m.db.Ping(); err != nil {
		logger.Err(err.Error())
		return false
	}
	return true
}

type EngineDatabase struct {
	Config config.Config `inject:"config"`
	Database
}

func (m *EngineDatabase) Shutdown() {
	if m.db != nil {
		logger.Info("Closing database connection...")
		m.db.Close()
	}
}

func (m *EngineDatabase) StartUp() {
	logger.Info("Init database connection...")

	info := "asipcnt-logging-775451169198-dc88e4c4e4897c39"
	db, err := sql.Open("athena", "db=s3log&region=ap-southeast-1&output_location=s3://"+info+"/QUERY_OUTPUT/")
	if err != nil {
		logger.Warn("Failed to connect [%s]", info)
	} else if err := db.Ping(); err != nil {
		logger.Err("Error while connecting to [%s]", info)
	} else {
		logger.Info("Successfully connected to [%s]", info)
	}
	m.db = db
}
