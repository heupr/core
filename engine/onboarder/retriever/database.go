package retriever

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseAccess interface {
	Open()
	Close()
	Query(string, ...interface{}) (*sql.Rows, error)
}

type MemSQL struct {
	db *sql.DB // NOTE: Should this be a DatabaseAccess interface value instead?
}

func (m *MemSQL) Open() {
	mysql, err := sql.Open("mysql", "root@/heupr?interpolateParams=true")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	m.db = mysql
}

func (m *MemSQL) Close() {
	m.db.Close()
}
