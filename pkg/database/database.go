package database

import (
	"database/sql"
	"github.com/ahmadrezamusthafa/logwatcher/config"
	"github.com/ahmadrezamusthafa/logwatcher/pkg/logger"
	"github.com/ahmadrezamusthafa/logwatcher/shared"
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
	dbMap  map[shared.ServiceName]*sql.DB
}

func (m Database) GetDB(svcName shared.ServiceName) *sql.DB {
	if _, ok := m.dbMap[svcName]; ok {
		return m.dbMap[svcName]
	}
	return nil
}

func (m Database) QueryRow(svcName shared.ServiceName, query string, args ...interface{}) *sql.Row {
	return m.GetDB(svcName).QueryRow(query, args)
}

func (m Database) In(query string, params map[string]interface{}) (string, []interface{}, error) {
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return query, args, err
	}
	return sqlx.In(query, args...)
}

type EngineDatabase struct {
	Config config.Config `inject:"config"`
	Database
}

func (m *EngineDatabase) Shutdown() {
	if m.dbMap != nil {
		for _, db := range m.dbMap {
			if db != nil {
				logger.Info("Closing database connection...")
				db.Close()
			}
		}
	}
}

func (m *EngineDatabase) StartUp() {
	logger.Info("Init database connection...")

	m.dbMap = make(map[shared.ServiceName]*sql.DB)
	for svcName, bucket := range shared.MapS3Bucket {
		db, err := sql.Open("athena", "db=s3log&region=ap-southeast-1&output_location=s3://"+bucket+"/QUERY_OUTPUT/")
		if err != nil {
			logger.Warn("Failed to connect [%s]", bucket)
		} else if err := db.Ping(); err != nil {
			logger.Err("Error while connecting to [%s]", bucket)
		} else {
			logger.Info("Successfully connected to [%s]", bucket)
		}
		if _, ok := m.dbMap[svcName]; !ok {
			m.dbMap[svcName] = db
		}
	}
}
