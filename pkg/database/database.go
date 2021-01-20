package database

import (
	"database/sql"
	"fmt"
	"github.com/ahmadrezamusthafa/logwatcher/common/logger"
	"github.com/ahmadrezamusthafa/logwatcher/config"
	"github.com/jmoiron/sqlx"
	athenadriver "github.com/uber/athenadriver/go"
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

/*func (m Database) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	return m.db.BindNamed(query, arg)
}*/

func (m Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.db.QueryRow(query, args)
}

/*func (m Database) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return m.db.NamedExec(query, arg)
}

func (m Database) Get(dest interface{}, query string, args ...interface{}) (err error) {
	return m.db.Get(dest, query, args...)
}

func (m Database) Select(dest interface{}, query string, args ...interface{}) (err error) {
	return m.db.Select(dest, query, args...)
}*/

func (m Database) In(query string, params map[string]interface{}) (string, []interface{}, error) {
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return query, args, err
	}
	return sqlx.In(query, args...)
}

/*func (m Database) Prepare(query string) (*sqlx.NamedStmt, error) {
	return m.db.PrepareNamed(query)
}

func (m Database) PrepareBind(query string) (*sqlx.Stmt, error) {
	return m.db.Preparex(query)
}

func (m Database) Rebind(query string) string {
	return m.db.Rebind(query)
}
*/
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
	conf := m.Config
	connectionString := "./" + conf.DatabaseFile
	info := fmt.Sprintf(connectionString)

	athenaConf, _ := athenadriver.NewDefaultConfig("s3://asipcnt-logging-787755326983-c72641f30d2db7e4/QUERY_OUTPUT/",
		"ap-southeast-1", "ASIAQHBFIXI6QNOJ2GFW", "FwoGZXIvYXdzEKj//////////wEaDMwj8b0tAfyI22IcyyKvAXnjy9BQZa09KxTKyaftPBYT/MBEQJDq4Z5q26YLmxgUMQiHXhZx/rh2vzj4RLIfYLmBVmOCjhdodCEKf/hFbC6z8i0XAymbqpKPoKl0JY0s64mNXFKavG4YIcrNa1QVg1S9BUFR30vKjt5jUlUbDYz4fkoVtoIUC1MJ5XLeY9YUz5vhcRH3ea5LEKWAFDT31Y6reVRzF0ayLwOSV5+ivGVMPM9VZxsPBoxTY4qE0AIot/WZgAYyLVYngxbf6J/r2byluuQ7uzvDIjJkS3BNib6a5IiG/agPXwNf3ySJUDOys41bVg==")

	/*wg := athenadriver.NewWG("asipcnt-logger", nil, nil)
	athenaConf.SetWorkGroup(wg)
	athenaConf.SetWGRemoteCreationAllowed(false)*/

	db, err := sql.Open(athenadriver.DriverName, athenaConf.Stringify())
	var timestamp string
	if db != nil {
		row := db.QueryRow(`SELECT timestamp from "s3log"."pcntapirqrs_init" limit 1`).Scan(&timestamp)
		if row != nil {
			println(row.Error())
		}
	}
	if err != nil {
		logger.Warn("Failed to connect [%s]", info)
	} else if err := db.Ping(); err != nil {
		logger.Err("Error while connecting to [%s]", info)
	} else {
		logger.Info("Successfully connected to [%s]", info)
	}
	m.db = db
}
