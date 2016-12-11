package bhattacharya

import (
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type CoralReefLogger struct {
	Name     string
	CRLog    *log.Logger
	CRBuff   *bytes.Buffer
	Backtest bool
}

// DOC: CreateLog generates a new proprietary logger structure.
//      Note that the "prefix" set in the logger is just the log file name.
func CreateLog(name string, backtest bool) CoralReefLogger {
	buf := bytes.Buffer{}
	logger := log.New(&buf, name, log.Lshortfile)
	newLog := CoralReefLogger{Name: name, CRLog: logger, CRBuff: &buf, Backtest: backtest}
	return newLog
}

// DOC: Log calls the logging input functionality.
func (crl *CoralReefLogger) Log(input interface{}) {
	crl.CRLog.Print(input)
}

// DOC: Flush closes out and writes the results of the log to the desired file.
func (crl *CoralReefLogger) Flush() {
	n := time.Now()
	y, m, d := n.Date()
	filename := crl.Name + "-" + strconv.Itoa(y) + "-" + m.String() + "-" + strconv.Itoa(d)

    // DOC: This logic evalutes whether the log should be generated in the
    //      backtests/ or logs/ directories.
    if crl.Backtest {
        crl.pathfinder("../../data/backtests/" + filename, n)
    } else {
        crl.pathfinder("../../data/logs/" + filename, n)
    }

}

// DOC: pathfinder is a helper method designed to perform the necessary file
//      creation / writing processes. Note that the logic here simply
//      evaluates whether the file exists in the target path.
func (crl *CoralReefLogger) pathfinder(filepath string, t time.Time) {
    h, m, _ := t.Clock()
    header := "LOG OUTPUT TIME\nHOUR: " + strconv.Itoa(h) + " " + "MINUTE: " + strconv.Itoa(m)
    if _, err := os.Stat(filepath); err == nil {
        f, e := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0600)
        if e != nil {
            panic(e)
        }
        io.WriteString(f, header + "\n\n" + crl.CRBuff.String() + "\n\n")
    } else {
        f, e := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
        if e != nil {
            panic(e)
        }
        io.WriteString(f, header + "\n\n" + crl.CRBuff.String() + "\n\n")
    }
}
