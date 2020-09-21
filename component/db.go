package component

import (
	"database/sql"
	"fmt"
)

//InitializeDatabase to initialize database
func InitializeDatabase() (db *sql.DB, err error) {

	db, err = sql.Open("mysql", "")
	if err != nil {
		return nil, fmt.Errorf("failed to open DB master connection. %+v", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(50)
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping DB master. %+v", err)
	}

	return db, err
}