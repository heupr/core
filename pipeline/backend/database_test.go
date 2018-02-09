package backend

import (
	"database/sql/driver"
	"io"
)

type testDB struct {
	name string
}

type testDriver struct{}

func (td testDriver) Open(name string) (driver.Conn, error) {
	db := &testDB{name: name}
	conn := &testConn{db: db}
	return conn, nil
}

type testConn struct {
	db *testDB
}

func (tc testConn) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

func (tc testConn) Close() error {
	return nil
}

func (tc testConn) Begin() (driver.Tx, error) {
	return nil, nil
}

func (tc *testConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	tr := testRows{}
	return tr, nil
}

type testStmt struct{}

func (ts testStmt) Close() error {
	return nil
}

func (ts testStmt) NumInput() int {
	return 0
}

func (ts testStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, nil
}

func (ts testStmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, nil
}

type testRows struct {
	rowsi    driver.Rows
	cancel   func()
	closed   bool
	lasterr  error
	lastcols []driver.Value
	// NOTE: See https://github.com/golang/go/blob/master/src/database/sql/fakedb_test.go#L858
}

func (tr testRows) Columns() []string {
	out := make([]string, 4)
	return out
}

func (tr testRows) Close() error {
	return nil
}

var startstop = true

func (tr testRows) Next(dest []driver.Value) error {
	if startstop {
		// NOTE: Below is a commented list of the corresponding database schema.
		dest[0] = 1     // id
		dest[1] = 1     // repo_id
		dest[2] = false // is_pull
		dest[3] = []byte(`{"issue":{"id":1,"closed_at":"2015-01-13T22:29:04Z"}}`)
		startstop = false
		return nil
	} else {
		return io.EOF
	}
}
