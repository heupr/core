package retriever

import (
    "database/sql"
    "database/sql/driver"
    "testing"
)

type testConn struct{}

type Conn interface {
	Prepare(query string) (Stmt, error)
	Close() error
	Begin() (Tx, error)
}

type testDriver struct {}

func (td testDriver) Open(name string) (sql.Conn, error) {

}







func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func newTestDB(t testing.TB, name string) *sql.DB {
	db, err := sql.Open("test", fakeDBName)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if _, err := db.Exec("WIPE"); err != nil {
		t.Fatalf("exec wipe: %v", err)
	}
	return db
}

var testMemSQL = MemSQL{}

// [ ] register new driver
// [ ] generate fake database












/*
import "database/sql"

type testRows struct {
    contents    string
}

func (tr testRows) Columns() []string{
    return []string{}
}

func (tr testRows) Close() error {
    return error
}

func (tr testRows) Next(dest []sql.Value) error {
    return error
}

type testDB struct {}

func (tdb testDB) Open() {}

func (tdb testDB) Close() {}

func (tdb testDB) Query(query string, args ...interface{}) (*testRows, error) {
    tr := &testRows{}
    tr.contents = "stuff goes here"
    return tr, nil
}

var testMemSQL = MemSQL{
    // db: testDB goes here
}
*/
