package backend

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MemSQL struct {
	db *sql.DB
}

func (m *MemSQL) Open() {
	mysql, err := sql.Open("mysql", "root@/heupr?interpolateParams=true")
	if err != nil {
		panic(err.Error()) // TODO: Proper error handling.
	}
	m.db = mysql
}

func (m *MemSQL) Close() {
	m.db.Close()
}
