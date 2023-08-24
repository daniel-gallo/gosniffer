package persistance

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net"
)

type SQLite struct {
	db *sql.DB
}

const initDBQuery string = `
CREATE TABLE IF NOT EXISTS log(
    id INTEGER NOT NULL PRIMARY KEY,
    module VARCHAR,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    ip VARCHAR,
    mac VARCHAR,
    message VARCHAR
);
`

const saveLogQuery string = `
INSERT INTO log(module, ip, mac, message)
VALUES (?, ?, ?, ?)
`

func CreateSQLite(filename string) SQLite {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(initDBQuery)
	if err != nil {
		panic(err)
	}

	return SQLite{db}
}

func (sqlite SQLite) Save(module string, ip net.IP, mac net.HardwareAddr, message string) {
	res, err := sqlite.db.Exec(saveLogQuery, module, ip.String(), mac.String(), message)
	if err != nil {
		panic(err)
	}

	_, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}
}
