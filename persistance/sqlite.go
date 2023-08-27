package persistance

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net"
	"time"
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
const loadLogsQuery string = `
SELECT module, timestamp, ip, mac, message
FROM log
ORDER BY id DESC
LIMIT ?
`

func NewSQLite(filename string) SQLite {
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

func (sqlite SQLite) Load(numLogs int) []Log {
	rows, err := sqlite.db.Query(loadLogsQuery, numLogs)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	logs := make([]Log, 0)

	for rows.Next() {
		var module string
		var timestamp time.Time
		var ip string
		var mac string
		var message string

		rows.Scan(&module, &timestamp, &ip, &mac, &message)

		parsedMAC, err := net.ParseMAC(mac)
		if err != nil {
			panic(err)
		}

		log := Log{
			Module:    module,
			Timestamp: timestamp,
			Ip:        net.ParseIP(ip),
			Mac:       parsedMAC,
			Message:   message,
		}

		logs = append(logs, log)
	}

	return logs
}
