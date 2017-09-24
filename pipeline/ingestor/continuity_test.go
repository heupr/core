package ingestor

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt" // TEMPORARY
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"core/pipeline/frontend"
	"core/utils"
)

// This unit test employs an "easy" solution to the BoltDB filepath by simply
// overwriting the config value; however, a wrapper struct could be built
// around the BoltDB which would then implement a new interface with methods
// like open, close, etc.
// See https://github.com/boltdb/bolt/blob/master/cmd/bolt/main_test.go#L181

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

func (c *testConn) Query(query string, args []driver.Value) (driver.Rows, error) {
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

func destPopulator(dest *[]interface{}, new []interface{}) {
	*dest = new
}

func (tr testRows) Next(dest []driver.Value) error {
	if startstop {
		// NOTE: Below is a commented list of the corresponding database
		// schema.
		dest[0] = 1 // repo_id
		// NOTE: This is what the query would return if the number 2 issue was
		// missing from the table.
		dest[1] = 1     // startNum issue
		dest[2] = 3     // endNum issue
		dest[3] = false // is_pull
		startstop = false
		return nil
	} else {
		return io.EOF
	}
}

func Test_continuityCheck(t *testing.T) {
	token, err := json.Marshal(&oauth2.Token{AccessToken: "droideka"})
	if err != nil {
		t.Errorf("failed marshaling test token: %v", err)
	}
	// Set up fake BoltDB file and database.
	name := "continuity-test.db"
	utils.Config.BoltDBPath = name // NOTE: Easy config option.
	file, err := ioutil.TempFile("", name)
	if err != nil {
		t.Errorf("generate continuity test file %v", err)
	}
	file.Close()
	defer os.Remove(name)
	b, err := bolt.Open(name, 0644, nil)
	if err != nil {
		t.Errorf("opening test bolt db %v", err)
	}
	boltDB := frontend.BoltDB{DB: b}
	if err := boltDB.Initialize(); err != nil {
		t.Errorf("initialize test bolt db %v", err)
	}
	if err := boltDB.Store("token", 1, token); err != nil {
		t.Errorf("setting bolt db test values %v", err)
	}
	boltDB.DB.Close()

	// Set up fake GitHub server; note that this unit test is only meant to
	// have a single call to the GitHub server - it would be to find the
	// missing issue (issue number 2).
	mux := http.NewServeMux()
	// Returning the test repo structure.
	mux.HandleFunc("/repositories/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":1,"name":"trade-federation","owner":{"login":"nute-gunray"}}`)
	})
	// Returning issues for test repo.
	mux.HandleFunc("/repos/nute-gunray/trade-federation/issues/2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":2,"number":2}`)
	})
	server := httptest.NewServer(mux)
	testURL, _ := url.Parse(server.URL)

	NewClient = func(t oauth2.Token) *github.Client {
		c := github.NewClient(nil)
		c.BaseURL = testURL
		c.UploadURL = testURL
		return c
	}

	// Create fake MemSQL database (see above for variable settings).
	driverName := "cato-nemoidia"
	sourceName := "purse-world"
	td := testDriver{}

	sql.Register(driverName, td)
	db, err := sql.Open(driverName, sourceName)
	if err != nil {
		t.Errorf("error opening test database %v: %v", sourceName, err)
	}
	testIS := IngestorServer{
		Database: Database{
			db: db,
		},
	}

	_, _, err = testIS.continuityCheck()
	if err != nil {
		t.Errorf("continuity check test: %v", err)
	}
}
